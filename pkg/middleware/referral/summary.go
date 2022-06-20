package referral

import (
	"context"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	cache "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/cache"
	cachekey "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/referral/cachekey"
	appusermgrpb "github.com/NpoolPlatform/message/npool/appusermgr"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"
	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"

	"golang.org/x/xerrors"
)

const (
	cacheReferralUser         = "referral:user"
	cacheReferralExtra        = "referral:extra"
	cacheLayeredCoinSummaries = "referral:layered:coin:summaries"
	cacheLayeredGoodSummaries = "referral:layered:good:summaries"
)

func getReferralUser(ctx context.Context, appID, userID string) (*appusermgrpb.AppUser, error) {
	user := cache.GetEntry(cachekey.CacheKey(appID, userID, cacheReferralUser))
	if user != nil {
		return user.(*appusermgrpb.AppUser), nil
	}

	user, err := grpc2.GetAppUserByAppUser(ctx, &appusermgrpb.GetAppUserByAppUserRequest{
		AppID:  appID,
		UserID: userID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get app user: %v", err)
	}
	if user.(*appusermgrpb.AppUser) == nil {
		return nil, xerrors.Errorf("invalid app user")
	}

	cache.AddEntry(cachekey.CacheKey(appID, userID, cacheReferralUser), user)

	return user.(*appusermgrpb.AppUser), nil
}

func getReferralExtra(ctx context.Context, appID, userID string) (*appusermgrpb.AppUserExtra, error) {
	extra := cache.GetEntry(cachekey.CacheKey(appID, userID, cacheReferralExtra))
	if extra != nil {
		return extra.(*appusermgrpb.AppUserExtra), nil
	}

	extra, err := grpc2.GetAppUserExtraByAppUser(ctx, &appusermgrpb.GetAppUserExtraByAppUserRequest{
		AppID:  appID,
		UserID: userID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get app user extra: %v", err)
	}

	cache.AddEntry(cachekey.CacheKey(appID, userID, cacheReferralExtra), extra)

	return extra.(*appusermgrpb.AppUserExtra), nil
}

func getLayeredGoodSummaries(ctx context.Context, appID, userID string) ([]*npool.GoodSummary, error) {
	mySummaries := cache.GetEntry(cachekey.CacheKey(appID, userID, cacheLayeredGoodSummaries))
	if mySummaries != nil {
		return mySummaries.([]*npool.GoodSummary), nil
	}

	goodSummaries, err := getGoodSummaries(ctx, appID, userID)
	if err != nil {
		return nil, xerrors.Errorf("fail get good summaries: %v", err)
	}

	invitees, err := GetLayeredInvitees(ctx, appID, userID)
	if err != nil {
		return nil, xerrors.Errorf("fail get invitees: %v", err)
	}

	sums := make([]*npool.GoodSummary, len(goodSummaries))
	for i, sum := range goodSummaries {
		sums[i] = &npool.GoodSummary{
			GoodID:     sum.GoodID,
			CoinTypeID: sum.CoinTypeID,
			CoinName:   sum.CoinName,
			Units:      sum.Units,
			Unit:       sum.Unit,
			Amount:     sum.Amount,
			Percent:    sum.Percent,
		}
	}

	for _, iv := range invitees {
		summaries, err := getGoodSummaries(ctx, iv.AppID, iv.InviteeID)
		if err != nil {
			return nil, xerrors.Errorf("fail get good summaries: %v", err)
		}

		for _, sum1 := range sums {
			for _, sum2 := range summaries {
				if sum1.GoodID == sum2.GoodID {
					sum1.Units += sum2.Units
					sum1.Amount += sum2.Amount
				}
			}
		}
	}

	if len(sums) > 0 {
		cache.AddEntry(cachekey.CacheKey(appID, userID, cacheLayeredGoodSummaries), sums)
	}

	return sums, nil
}

func getLayeredCoinSummaries(ctx context.Context, appID, userID string) ([]*npool.CoinSummary, error) {
	mySummaries := cache.GetEntry(cachekey.CacheKey(appID, userID, cacheLayeredCoinSummaries))
	if mySummaries != nil {
		return mySummaries.([]*npool.CoinSummary), nil
	}

	coinSummaries, err := getCoinSummaries(ctx, appID, userID)
	if err != nil {
		return nil, xerrors.Errorf("fail get coin summaries: %v", err)
	}

	invitees, err := GetLayeredInvitees(ctx, appID, userID)
	if err != nil {
		return nil, xerrors.Errorf("fail get invitees: %v", err)
	}

	sums := make([]*npool.CoinSummary, len(coinSummaries))
	for i, sum := range coinSummaries {
		sums[i] = &npool.CoinSummary{
			CoinTypeID: sum.CoinTypeID,
			CoinName:   sum.CoinName,
			Units:      sum.Units,
			Unit:       sum.Unit,
			Amount:     sum.Amount,
		}
	}

	for _, iv := range invitees {
		summaries, err := getCoinSummaries(ctx, iv.AppID, iv.InviteeID)
		if err != nil {
			return nil, xerrors.Errorf("fail get coin summaries: %v", err)
		}

		for _, sum1 := range sums {
			for _, sum2 := range summaries {
				if sum1.CoinTypeID == sum2.CoinTypeID {
					sum1.Units += sum2.Units
					sum1.Amount += sum2.Amount
				}
			}
		}
	}

	if len(sums) > 0 {
		cache.AddEntry(cachekey.CacheKey(appID, userID, cacheLayeredCoinSummaries), sums)
	}

	return sums, nil
}

func getReferral(ctx context.Context, appID, userID string) (*npool.Referral, error) {
	user, err := getReferralUser(ctx, appID, userID)
	if err != nil {
		return nil, xerrors.Errorf("fail get referral user: %v", err)
	}
	if user == nil {
		return nil, xerrors.Errorf("invalid referral user")
	}

	extra, err := getReferralExtra(ctx, appID, userID)
	if err != nil {
		return nil, xerrors.Errorf("fail get referral extra: %v", err)
	}

	inviter, err := GetInviter(ctx, appID, userID)
	if err != nil {
		return nil, xerrors.Errorf("fail get inviter: %v", err)
	}

	amount, err := GetUSDAmount(ctx, appID, userID)
	if err != nil {
		return nil, xerrors.Errorf("fail get usd amount: %v", err)
	}

	subAmount, err := GetSubUSDAmount(ctx, appID, userID)
	if err != nil {
		return nil, xerrors.Errorf("fail get sub usd amount: %v", err)
	}

	invitees, err := GetInvitees(ctx, appID, userID)
	if err != nil {
		return nil, xerrors.Errorf("fail get invitees: %v", err)
	}

	coinSummaries, err := getLayeredCoinSummaries(ctx, appID, userID)
	if err != nil {
		return nil, xerrors.Errorf("fail get coin summaries: %v", err)
	}

	goodSummaries, err := getLayeredGoodSummaries(ctx, appID, userID)
	if err != nil {
		return nil, xerrors.Errorf("fail get good summaries: %v", err)
	}

	code, err := grpc2.GetUserInvitationCodeByAppUser(ctx, &inspirepb.GetUserInvitationCodeByAppUserRequest{
		AppID:  appID,
		UserID: userID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get user invitation code: %v", err)
	}

	return &npool.Referral{
		User:          user,
		Extra:         extra,
		Invitation:    inviter,
		USDAmount:     amount,
		SubUSDAmount:  subAmount,
		Kol:           code != nil,
		InvitedCount:  uint32(len(invitees)),
		Summaries:     coinSummaries,
		GoodSummaries: goodSummaries,
	}, nil
}

func getReferrals(ctx context.Context, appID, userID string) ([]*npool.Referral, error) {
	invitees, err := GetInvitees(ctx, appID, userID)
	if err != nil {
		return nil, xerrors.Errorf("fail get invitees: %v", err)
	}

	referrals := []*npool.Referral{}

	for _, iv := range invitees {
		referral, err := getReferral(ctx, iv.AppID, iv.InviteeID)
		if err != nil {
			return nil, xerrors.Errorf("fail get referral: %v", err)
		}
		referrals = append(referrals, referral)
	}

	return referrals, nil
}

func getLayeredReferrals(ctx context.Context, appID, userID string) ([]*npool.Referral, error) {
	invitees, err := GetLayeredInvitees(ctx, appID, userID)
	if err != nil {
		return nil, xerrors.Errorf("fail get invitees: %v", err)
	}

	referrals := []*npool.Referral{}

	for _, iv := range invitees {
		referral, err := getReferral(ctx, iv.AppID, iv.InviteeID)
		if err != nil {
			return nil, xerrors.Errorf("fail get referral: %v", err)
		}
		referrals = append(referrals, referral)
	}

	return referrals, nil
}
