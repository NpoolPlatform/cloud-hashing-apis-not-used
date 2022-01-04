package user

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	"github.com/NpoolPlatform/cloud-hashing-apis/message/npool"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	order "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/order"

	appmgrpb "github.com/NpoolPlatform/application-management/message/npool"
	inspirepb "github.com/NpoolPlatform/cloud-hashing-inspire/message/npool"
	orderconst "github.com/NpoolPlatform/cloud-hashing-order/pkg/const"
	usermgrpb "github.com/NpoolPlatform/user-management/message/npool"

	"golang.org/x/xerrors"
)

func convertUserinfo(info *usermgrpb.UserBasicInfo) *npool.UserInfo {
	return &npool.UserInfo{
		UserID:         info.UserID,
		Username:       info.Username,
		Avatar:         info.Avatar,
		Age:            info.Age,
		Gender:         info.Gender,
		Region:         info.Region,
		Birthday:       info.Birthday,
		Country:        info.Country,
		Province:       info.Province,
		City:           info.City,
		PhoneNumber:    info.PhoneNumber,
		EmailAddress:   info.EmailAddress,
		CreateAt:       info.CreateAt,
		LoginTimes:     info.LoginTimes,
		KycVerify:      info.KycVerify,
		GaVerify:       info.GaVerify,
		GaLogin:        info.GaLogin,
		SignupMethod:   info.SignupMethod,
		Career:         info.Career,
		DisplayName:    info.DisplayName,
		FirstName:      info.FirstName,
		LastName:       info.LastName,
		StreetAddress1: info.StreetAddress1,
		StreetAddress2: info.StreetAddress2,
	}
}

func Signup(ctx context.Context, in *npool.SignupRequest) (*npool.SignupResponse, error) { //nolint
	invitationCode := in.GetInvitationCode()
	inviterID := ""

	appResp, err := grpc2.GetApp(ctx, &appmgrpb.GetApplicationRequest{
		AppID: in.GetAppID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get app: %v", err)
	}

	if appResp.Info.InvitationCodeMust {
		if invitationCode == "" {
			return nil, xerrors.Errorf("invitation code is must")
		}
	}

	if invitationCode != "" {
		getByCodeResp, err := grpc2.GetUserInvitationCodeByCode(ctx, &inspirepb.GetUserInvitationCodeByCodeRequest{
			Code: invitationCode,
		})
		if err != nil {
			return nil, xerrors.Errorf("fail get user invitation code: %v", err)
		}

		if getByCodeResp.Info == nil {
			if appResp.Info.InvitationCodeMust {
				return nil, xerrors.Errorf("fail get invitation code")
			}
		} else {
			if getByCodeResp.Info.AppID != in.GetAppID() {
				return nil, xerrors.Errorf("invalid invitation code for app")
			}
			inviterID = getByCodeResp.Info.UserID
		}
	}

	signupResp, err := grpc2.Signup(ctx, &usermgrpb.SignupRequest{
		Username:     in.GetUsername(),
		Password:     in.GetPassword(),
		EmailAddress: in.GetEmailAddress(),
		PhoneNumber:  in.GetPhoneNumber(),
		Code:         in.GetVerificationCode(),
		AppID:        in.GetAppID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail signup: %v", err)
	}

	if invitationCode != "" && inviterID != "" {
		_, err = grpc2.CreateRegistrationInvitation(ctx, &inspirepb.CreateRegistrationInvitationRequest{
			Info: &inspirepb.RegistrationInvitation{
				AppID:     in.GetAppID(),
				InviterID: inviterID,
				InviteeID: signupResp.Info.UserID,
			},
		})
		if err != nil {
			return nil, xerrors.Errorf("fail create registration invitation: %v", err)
		}
	}

	return &npool.SignupResponse{
		Info: convertUserinfo(signupResp.Info),
	}, nil
}

func getInvitations(ctx context.Context, appID, reqInviterID string, directOnly bool) (map[string]*npool.Invitation, error) { //nolint
	_, err := grpc2.GetUser(ctx, &usermgrpb.GetUserRequest{
		AppID:  appID,
		UserID: reqInviterID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get inviter %v user information: %v", reqInviterID, err)
	}

	goon := true
	invitations := map[string]*npool.Invitation{}
	invitations[reqInviterID] = &npool.Invitation{
		Invitees: []*npool.InvitationUserInfo{},
	}
	inviters := map[string]struct{}{}

	// TODO: process deadloop
	for goon {
		goon = false

		for inviterID, _ := range invitations { //nolint
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

			for _, info := range resp.Infos {
				if info.AppID != appID || info.InviterID != inviterID {
					logger.Sugar().Errorf("invalid inviter id or app id")
					continue
				}

				inviteeResp, err := grpc2.GetUser(ctx, &usermgrpb.GetUserRequest{
					AppID:  appID,
					UserID: info.InviteeID,
				})
				if err != nil {
					logger.Sugar().Errorf("fail get invitee %v user info: %v", info.InviteeID, err)
					continue
				}

				resp1, err := grpc2.GetUserInvitationCodeByAppUser(ctx, &inspirepb.GetUserInvitationCodeByAppUserRequest{
					AppID:  appID,
					UserID: inviteeResp.Info.UserID,
				})
				if err != nil {
					logger.Sugar().Errorf("fail get user invitation code: %v", err)
					continue
				}

				resp2, err := order.GetOrdersDetailByAppUser(ctx, &npool.GetOrdersDetailByAppUserRequest{
					AppID:  appID,
					UserID: inviteeResp.Info.UserID,
				})
				if err != nil {
					logger.Sugar().Errorf("fail get orders detail by app user: %v", err)
					continue
				}

				summarys := map[string]*npool.InvitationSummary{}

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

				if _, ok := invitations[inviterID]; !ok {
					invitations[inviterID] = &npool.Invitation{
						Invitees: []*npool.InvitationUserInfo{},
					}
				}

				invitations[inviterID].Invitees = append(
					invitations[inviterID].Invitees, &npool.InvitationUserInfo{
						UserID:       inviteeResp.Info.UserID,
						Username:     inviteeResp.Info.Username,
						Avatar:       inviteeResp.Info.Avatar,
						EmailAddress: inviteeResp.Info.EmailAddress,
						Kol:          kol,
						Summarys:     summarys,
					})

				if _, ok := invitations[inviteeResp.Info.UserID]; !ok {
					invitations[inviteeResp.Info.UserID] = &npool.Invitation{
						Invitees: []*npool.InvitationUserInfo{},
					}
				}

				goon = true
			}
		}
	}

	invitation := invitations[reqInviterID]

	for _, invitee := range invitation.Invitees {
		curInviteeIDs := []string{invitee.UserID}
		foundInvitees := map[string]struct{}{}
		goon := true

		for goon {
			goon = false

			for _, curInviteeID := range curInviteeIDs {
				if _, ok := foundInvitees[curInviteeID]; ok {
					continue
				}

				foundInvitees[curInviteeID] = struct{}{}

				invitation := invitations[curInviteeID]

				for _, iv := range invitation.Invitees {
					curInviteeIDs = append(curInviteeIDs, iv.UserID)

					for coinID, summary := range iv.Summarys {
						if _, ok := invitee.Summarys[coinID]; !ok {
							invitee.Summarys[coinID] = &npool.InvitationSummary{}
						}
						// TODO: process different payment coin type
						mySummary := invitee.Summarys[coinID]
						mySummary.Units += summary.Units
						mySummary.Amount += summary.Amount
						invitee.Summarys[coinID] = mySummary
					}

					goon = true
				}
			}
		}
	}

	if directOnly {
		return map[string]*npool.Invitation{
			reqInviterID: invitation,
		}, nil
	}

	invitations[reqInviterID] = invitation

	return invitations, nil
}

func GetMyInvitations(ctx context.Context, in *npool.GetMyInvitationsRequest) (*npool.GetMyInvitationsResponse, error) { //nolint
	invitations, err := getInvitations(ctx, in.GetAppID(), in.GetInviterID(), false)
	if err != nil {
		return nil, xerrors.Errorf("fail get invitations: %v", err)
	}
	return &npool.GetMyInvitationsResponse{
		Infos: invitations,
	}, nil
}

func GetMyDirectInvitations(ctx context.Context, in *npool.GetMyDirectInvitationsRequest) (*npool.GetMyDirectInvitationsResponse, error) { //nolint
	invitations, err := getInvitations(ctx, in.GetAppID(), in.GetInviterID(), true)
	if err != nil {
		return nil, xerrors.Errorf("fail get invitations: %v", err)
	}
	return &npool.GetMyDirectInvitationsResponse{
		Infos: invitations,
	}, nil
}
