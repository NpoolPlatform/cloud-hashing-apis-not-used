package referral

import (
	"context"
	"fmt"

	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"
)

func GetReferrals(ctx context.Context, in *npool.GetReferralsRequest) (*npool.GetReferralsResponse, error) { //nolint
	referrals, err := getReferrals(ctx, in.GetAppID(), in.GetUserID())
	if err != nil {
		return nil, fmt.Errorf("fail get referrals: %v", err)
	}

	_referral, err := getReferral(ctx, in.GetAppID(), in.GetUserID())
	if err != nil {
		return nil, fmt.Errorf("fail get referral: %v", err)
	}

	_referrals := []*npool.Referral{_referral}
	_referrals = append(_referrals, referrals...)

	return &npool.GetReferralsResponse{
		Infos: _referrals,
	}, nil
}

func GetLayeredReferrals(ctx context.Context, in *npool.GetLayeredReferralsRequest) (*npool.GetLayeredReferralsResponse, error) { //nolint
	referrals, err := getLayeredReferrals(ctx, in.GetAppID(), in.GetUserID())
	if err != nil {
		return nil, fmt.Errorf("fail get referrals: %v", err)
	}

	_referral, err := getReferral(ctx, in.GetAppID(), in.GetUserID())
	if err != nil {
		return nil, fmt.Errorf("fail get referral: %v", err)
	}

	_referrals := []*npool.Referral{_referral}
	_referrals = append(_referrals, referrals...)

	return &npool.GetLayeredReferralsResponse{
		Infos: _referrals,
	}, nil
}
