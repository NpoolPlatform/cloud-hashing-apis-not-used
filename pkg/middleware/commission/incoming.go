package commission

import (
	"context"

	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/referral"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"
	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"

	"golang.org/x/xerrors"
)

type Incoming struct {
	SuperUSDAmount float64
	USDAmount      float64 // Finally how much i and my followers got = (mineUSD + subUSD) * myPercent
	NetUSDAmount   float64 // Finally how much i got = USDAmount - SubUSDAmount
	SubUSDAmount   float64 // Finally how much i send = sum(subUSD * (myPercent - subPercent))
}

func getRebate(ctx context.Context, appID, userID string) (*Incoming, error) {
	settings, err := getAmountSettingsByAppUser(ctx, appID, userID)
	if err != nil {
		return 0, xerrors.Errorf("fail get amount settings: %v", err)
	}

	totalAmount := 0

	for _, setting := range settings {
		amount, err := referral.GetPeriodUSDAmount(ctx, appID, userID, setting.Start, setting.End)
		if err != nil {
			return 0, xerrors.Errorf("fail get period usd amount: %v", err)
		}
		totalAmount += amount * setting.Percent
	}

	return &Incoming{
		SuperUSDAmount: 0,
		USDAmount:      totalAmount,
		NetUSDAmount:   totalAmount,
	}, nil
}

func getPeriodRebate(ctx context.Context, appID, userID string, parentSettings []*inspirepb.AppPurchaseAmountSetting) (float64, error) {
	settings, err := getAmountSettingsByAppUser(ctx, appID, userID)
	if err != nil {
		return 0, xerrors.Errorf("fail get amount settings: %v", err)
	}

	totalAmount := 0
	for _, setting := range settings {
		amount, err := referral.GetPeriodUSDAmount(ctx, appID, userID, setting.Start, setting.End)
		if err != nil {
			return 0, xerrors.Errorf("fail get period usd amount: %v", err)
		}
		totalAmount += amount * setting.Percent
	}

	return &Incoming{
		SuperUSDAmount: 0,
		USDAmount:      totalAmount,
		NetUSDAmount:   totalAmount,
	}, nil
}

func getIncomings(ctx context.Context, appID, userID string) ([]*Incoming, error) {
	settings, err := getAmountSettingsByAppUser(ctx, appID, userID)
	if err != nil {
		return 0, xerrors.Errorf("fail get amount settings: %v", err)
	}

	invitees, err := referral.GetInvitees(ctx, appID, userID)
	if err != nil {
		return 0, xerrors.Errorf("fail get invitees: %v", err)
	}

	for _, iv := range invitees {
		rebase, err := getPeriodRebate(ctx, iv.AppID, iv.InviteeID)
		if err != nil {
			return 0, xerrors.Errorf("fail get rebate: %v", err)
		}

		subAmount, err := getLayeredPeriodUSDAmount(ctx, appID, userID, setting.Start, setting.End)
		if err != nil {
			return 0, xerrors.Errorf("fail get period sub usd amount: %v", err)
		}
		totalSubAmount += subAmount
	}
}

func getIncoming(ctx context.Context, appID, userID string) (*Incoming, error) {
	incoming, err := getRebate(ctx, appID, userID)
	if err != nil {
		return 0, xerrors.Errorf("fail get total incoming: %v", err)
	}

	incomings, err := getIncomings(ctx, appID, userID)
	if err != nil {
		return 0, xerrors.Errorf("fail get sub incomings: %v", err)
	}

	// TODO: add incomings

	return incoming, nil
}
