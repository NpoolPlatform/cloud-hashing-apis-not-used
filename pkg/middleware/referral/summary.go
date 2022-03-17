package referral

import (
	"context"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	appusermgrpb "github.com/NpoolPlatform/message/npool/appusermgr"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	"golang.org/x/xerrors"
)

func getReferral(ctx context.Context, appID, userID string) (*npool.Referral, error) {
	user, err := grpc2.GetAppUserByAppUser(ctx, &appusermgrpb.GetAppUserByAppUserRequest{
		AppID:  appID,
		UserID: userID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get app user: %v", err)
	}
	if user == nil {
		return nil, xerrors.Errorf("invalid app user")
	}

	extra, err := grpc2.GetAppUserExtraByAppUser(ctx, &appusermgrpb.GetAppUserExtraByAppUserRequest{
		AppID:  appID,
		UserID: userID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get app user extra: %v", err)
	}

	inviter, err := getInviter(ctx, appID, userID)
	if err != nil {
		return nil, xerrors.Errorf("fail get inviter: %v", err)
	}

	amount, err := getUSDAmount(ctx, appID, userID)
	if err != nil {
		return nil, xerrors.Errorf("fail get usd amount: %v", err)
	}

	subAmount, err := getSubUSDAmount(ctx, appID, userID)
	if err != nil {
		return nil, xerrors.Errorf("fail get sub usd amount: %v", err)
	}

	invitees, err := getInvitees(ctx, appID, userID)
	if err != nil {
		return nil, xerrors.Errorf("fail get invitees: %v", err)
	}

	return &npool.Referral{
		User:         user,
		Extra:        extra,
		Invitation:   inviter,
		USDAmount:    amount,
		SubUSDAmount: subAmount,
		Kol:          len(invitees) > 0,
	}, nil
}

func getReferrals(ctx context.Context, appID, userID string) ([]*npool.Referral, error) {
	invitees, err := getInvitees(ctx, appID, userID)
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
	invitees, err := getLayeredInvitees(ctx, appID, userID)
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
