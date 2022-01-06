package user

import (
	"context"
	"sync"
	"time"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	"github.com/NpoolPlatform/cloud-hashing-apis/message/npool"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	order "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/order"

	inspirepb "github.com/NpoolPlatform/cloud-hashing-inspire/message/npool"
	orderconst "github.com/NpoolPlatform/cloud-hashing-order/pkg/const"
	coininfopb "github.com/NpoolPlatform/message/npool/coininfo"
	usermgrpb "github.com/NpoolPlatform/user-management/message/npool"

	"golang.org/x/xerrors"
)

var (
	appInvitations        = map[string]map[string]map[string]*npool.Invitation{}
	appInviterUserInfos   = map[string]map[string]*npool.InvitationUserInfo{}
	appInviterLastUpdates = map[string]map[string]time.Time{}
	notifier              = make(chan struct{})
	mutex                 = sync.Mutex{}
)

func addWatcher(appID, inviterID string) {
	mutex.Lock()
	defer mutex.Unlock()

	if _, ok := appInvitations[appID]; !ok {
		appInvitations[appID] = map[string]map[string]*npool.Invitation{}
	}
	appInvitation := appInvitations[appID]
	if _, ok := appInvitation[inviterID]; !ok {
		appInvitation[inviterID] = map[string]*npool.Invitation{}
	}
	appInvitations[appID] = appInvitation

	if _, ok := appInviterUserInfos[appID]; !ok {
		appInviterUserInfos[appID] = map[string]*npool.InvitationUserInfo{}
	}
	appInviterUserInfo := appInviterUserInfos[appID]
	if _, ok := appInviterUserInfos[inviterID]; !ok {
		appInviterUserInfo[inviterID] = &npool.InvitationUserInfo{}
	}
	appInviterUserInfos[appID] = appInviterUserInfo

	if _, ok := appInviterLastUpdates[appID]; !ok {
		appInviterLastUpdates[appID] = map[string]time.Time{}
	}
	appInviterLastUpdate := appInviterLastUpdates[appID]
	if _, ok := appInviterLastUpdates[inviterID]; !ok {
		appInviterLastUpdate[inviterID] = time.Time{}
	}
	appInviterLastUpdates[appID] = appInviterLastUpdate

	go func() {
		notifier <- struct{}{}
	}()
}

func update() time.Duration {
	appInviters := map[string][]string{}

	mutex.Lock()
	for appID, inviterMap := range appInvitations {
		if _, ok := appInviters[appID]; !ok {
			appInviters[appID] = []string{}
		}
		myInviters := appInviters[appID]
		for inviterID := range inviterMap {
			myInviters = append(myInviters, inviterID)
		}
		appInviters[appID] = myInviters
	}
	mutex.Unlock()

	logger.Sugar().Infof("run async updater at %v", time.Now())
	toNext := 24 * time.Hour

	for appID, inviters := range appInviters {
		for _, inviterID := range inviters {
			nextSync := appInviterLastUpdates[appID][inviterID].Add(24 * time.Hour)
			if time.Now().Before(nextSync) {
				curToNext := time.Until(nextSync)
				if curToNext < toNext {
					toNext = curToNext
				}
				continue
			}

			logger.Sugar().Infof("run async updater for %v at %v", inviterID, time.Now())

			invitations, userInfo, err := getInvitations(appID, inviterID, false)
			if err != nil {
				logger.Sugar().Errorf("fail get invitations: %v", err)
				continue
			}

			mutex.Lock()
			appInvitations[appID][inviterID] = invitations
			appInviterUserInfos[appID][inviterID] = userInfo
			appInviterLastUpdates[appID][inviterID] = time.Now()
			mutex.Unlock()
		}
	}

	return toNext
}

func Run() {
	var timeout time.Duration

	timer := time.NewTimer(24 * time.Hour)

	for {
		select {
		case <-timer.C:
			timeout = update()
		case <-notifier:
			timeout = update()
		}
		timer.Reset(timeout)
	}
}

func getFullInvitations(appID, inviterID string) (map[string]*npool.Invitation, *npool.InvitationUserInfo, error) {
	mutex.Lock()
	invitations := appInvitations[appID][inviterID]
	userInfo := appInviterUserInfos[appID][inviterID]
	mutex.Unlock()

	if len(invitations) > 0 {
		return invitations, userInfo, nil
	}

	invitations, userInfo, err := getInvitations(appID, inviterID, false)
	if err != nil {
		return nil, nil, xerrors.Errorf("fail get invitations: %v", err)
	}

	return invitations, userInfo, nil
}

func getDirectInvitations(appID, inviterID string) (map[string]*npool.Invitation, *npool.InvitationUserInfo, error) {
	mutex.Lock()
	invitations := appInvitations[appID][inviterID]
	userInfo := appInviterUserInfos[appID][inviterID]
	mutex.Unlock()

	if len(invitations) > 0 {
		return invitations, userInfo, nil
	}

	invitations, userInfo, err := getInvitations(appID, inviterID, true)
	if err != nil {
		return nil, nil, xerrors.Errorf("fail get invitations: %v", err)
	}

	return invitations, userInfo, nil
}

func getInvitationUserInfo( //nolint
	appID, inviteeID string,
	myGoods map[string]*npool.GoodDetail,
	myCoins map[string]*coininfopb.CoinInfo) (*npool.InvitationUserInfo,
	map[string]*npool.GoodDetail,
	map[string]*coininfopb.CoinInfo,
	error) {
	ctx := context.Background()

	inviteeResp, err := grpc2.GetUser(ctx, &usermgrpb.GetUserRequest{
		AppID:  appID,
		UserID: inviteeID,
	})
	if err != nil {
		return nil, myGoods, myCoins, xerrors.Errorf("fail get invitee %v user info: %v", inviteeID, err)
	}

	resp1, err := grpc2.GetUserInvitationCodeByAppUser(ctx, &inspirepb.GetUserInvitationCodeByAppUserRequest{
		AppID:  appID,
		UserID: inviteeResp.Info.UserID,
	})
	if err != nil {
		return nil, myGoods, myCoins, xerrors.Errorf("fail get user invitation code: %v", err)
	}

	summarys := map[string]*npool.InvitationSummary{}

	resp2, goods, coins, err := order.GetOrdersShortDetailByAppUser(ctx, &npool.GetOrdersDetailByAppUserRequest{
		AppID:  appID,
		UserID: inviteeResp.Info.UserID,
	}, myGoods, myCoins)
	if err != nil {
		return nil, myGoods, myCoins, xerrors.Errorf("fail get orders detail by app user: %v", err)
	}

	myGoods = goods
	myCoins = coins

	for _, orderInfo := range resp2.Details {
		if orderInfo.Payment == nil {
			continue
		}

		if orderInfo.Payment.State != orderconst.PaymentStateDone {
			continue
		}

		if _, ok := summarys[orderInfo.Good.CoinInfo.ID]; !ok {
			summarys[orderInfo.Good.CoinInfo.ID] = &npool.InvitationSummary{}
		}

		summary := summarys[orderInfo.Good.CoinInfo.ID]
		summary.Units += orderInfo.Units
		summary.Amount += orderInfo.Payment.Amount
		summarys[orderInfo.Good.CoinInfo.ID] = summary
	}

	kol := false
	if resp1.Info != nil {
		kol = true
	}

	return &npool.InvitationUserInfo{
		UserID:       inviteeResp.Info.UserID,
		Username:     inviteeResp.Info.Username,
		Avatar:       inviteeResp.Info.Avatar,
		EmailAddress: inviteeResp.Info.EmailAddress,
		Kol:          kol,
		MySummarys:   summarys,
	}, myGoods, myCoins, nil
}

func getInvitations(appID, reqInviterID string, directOnly bool) (map[string]*npool.Invitation, *npool.InvitationUserInfo, error) { //nolint
	ctx := context.Background()

	_, err := grpc2.GetUser(ctx, &usermgrpb.GetUserRequest{
		AppID:  appID,
		UserID: reqInviterID,
	})
	if err != nil {
		return nil, nil, xerrors.Errorf("fail get inviter %v user information: %v", reqInviterID, err)
	}

	goon := true
	invitations := map[string]*npool.Invitation{}
	invitations[reqInviterID] = &npool.Invitation{
		Invitees: []*npool.InvitationUserInfo{},
	}
	inviters := map[string]struct{}{}
	myGoods := map[string]*npool.GoodDetail{}
	myCoins := map[string]*coininfopb.CoinInfo{}
	myCounts := map[string]uint32{}

	inviterUserInfo, goods, coins, err := getInvitationUserInfo(appID, reqInviterID, myGoods, myCoins)
	if err != nil {
		return nil, nil, xerrors.Errorf("fail get inviter %v user info: %v", reqInviterID, err)
	}

	myGoods = goods
	myCoins = coins
	layer := 0

	// TODO: process deadloop
	for goon {
		goon = false

		for inviterID := range invitations { //nolint
			if _, ok := inviters[inviterID]; ok {
				continue
			}

			inviters[inviterID] = struct{}{}

			resp, err := grpc2.GetRegistrationInvitationsByAppInviter(ctx, &inspirepb.GetRegistrationInvitationsByAppInviterRequest{
				AppID:     appID,
				InviterID: inviterID,
			})
			if err != nil {
				logger.Sugar().Errorf("fail get invitations by inviter %v: %v", inviterID, err)
				continue
			}

			myCounts[inviterID] = uint32(len(resp.Infos))

			for i, info := range resp.Infos {
				logger.Sugar().Infof("%v of %v layer %v user %v invited count %v", i, len(resp.Infos), layer, inviterID, myCounts[inviterID])

				if info.AppID != appID || info.InviterID != inviterID {
					logger.Sugar().Errorf("invalid inviter id or app id")
					continue
				}

				userInfo, goods, coins, err := getInvitationUserInfo(appID, info.InviteeID, myGoods, myCoins)
				if err != nil {
					logger.Sugar().Errorf("fail get invitation user info: %v", err)
					continue
				}

				userInfo.JoinDate = info.CreateAt
				myGoods = goods
				myCoins = coins

				if _, ok := invitations[inviterID]; !ok {
					invitations[inviterID] = &npool.Invitation{
						Invitees: []*npool.InvitationUserInfo{},
					}
				}

				invitations[inviterID].Invitees = append(invitations[inviterID].Invitees, userInfo)

				if !directOnly || layer < 2 {
					if _, ok := invitations[userInfo.UserID]; !ok {
						invitations[userInfo.UserID] = &npool.Invitation{
							Invitees: []*npool.InvitationUserInfo{},
						}
					}
				}

				goon = true
			}
		}

		layer++
	}

	invitation := invitations[reqInviterID]
	inviterUserInfo.InvitedCount = uint32(len(invitation.Invitees))
	if inviterUserInfo.Summarys == nil {
		inviterUserInfo.Summarys = map[string]*npool.InvitationSummary{}
	}

	for _, invitee := range invitation.Invitees {
		curInviteeIDs := []string{invitee.UserID}
		foundInvitees := map[string]struct{}{}
		goon := true

		invitee.InvitedCount = myCounts[invitee.UserID]

		for coinID, summary := range invitee.MySummarys {
			if _, ok := inviterUserInfo.Summarys[coinID]; !ok {
				inviterUserInfo.Summarys[coinID] = &npool.InvitationSummary{}
			}
			mySummary := inviterUserInfo.Summarys[coinID]
			mySummary.Units += summary.Units
			mySummary.Amount += summary.Amount
			invitee.Summarys[coinID] = mySummary
			inviterUserInfo.Summarys[coinID] = mySummary
		}

		for goon {
			goon = false

			for _, curInviteeID := range curInviteeIDs {
				if _, ok := foundInvitees[curInviteeID]; ok {
					continue
				}

				foundInvitees[curInviteeID] = struct{}{}

				invitation, ok := invitations[curInviteeID]
				if !ok {
					continue
				}

				for _, iv := range invitation.Invitees {
					curInviteeIDs = append(curInviteeIDs, iv.UserID)

					logger.Sugar().Infof("caculate %v invited count %v", iv.UserID, myCounts[iv.UserID])
					iv.InvitedCount = myCounts[iv.UserID]

					if invitee.Summarys == nil {
						invitee.Summarys = map[string]*npool.InvitationSummary{}
					}

					for coinID, summary := range iv.MySummarys {
						if _, ok := invitee.Summarys[coinID]; !ok {
							invitee.Summarys[coinID] = &npool.InvitationSummary{}
						}
						// TODO: process different payment coin type
						mySummary := invitee.Summarys[coinID]
						mySummary.Units += summary.Units
						mySummary.Amount += summary.Amount
						invitee.Summarys[coinID] = mySummary

						if _, ok := inviterUserInfo.Summarys[coinID]; !ok {
							inviterUserInfo.Summarys[coinID] = &npool.InvitationSummary{}
						}
						mySummary = inviterUserInfo.Summarys[coinID]
						mySummary.Units += summary.Units
						mySummary.Amount += summary.Amount
						invitee.Summarys[coinID] = mySummary
						inviterUserInfo.Summarys[coinID] = mySummary
					}

					goon = true
				}
			}
		}
	}

	if directOnly {
		return map[string]*npool.Invitation{
			reqInviterID: invitation,
		}, inviterUserInfo, nil
	}

	invitations[reqInviterID] = invitation

	return invitations, inviterUserInfo, nil
}
