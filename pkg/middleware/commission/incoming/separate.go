package incoming

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	commissionsetting "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/commission/setting"
	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/referral"
	orderconst "github.com/NpoolPlatform/cloud-hashing-order/pkg/const"
	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"
	orderpb "github.com/NpoolPlatform/message/npool/cloud-hashing-order"

	"golang.org/x/xerrors"
)

func getRebate(ctx context.Context, appID, userID string) (float64, error) {
	settings, err := commissionsetting.GetAmountSettingsByAppUser(ctx, appID, userID)
	if err != nil {
		return 0, xerrors.Errorf("fail get amount settings: %v", err)
	}

	totalAmount := 0.0

	for _, setting := range settings {
		amount, err := referral.GetPeriodUSDAmount(ctx, appID, userID, setting.Start, setting.End)
		if err != nil {
			return 0, xerrors.Errorf("fail get period usd amount: %v", err)
		}
		totalAmount += amount * float64(setting.Percent) / 100.0
	}

	return totalAmount, nil
}

func getOrderParentRebate(ctx context.Context, order *orderpb.OrderDetail, roots, nexts []*inspirepb.AppPurchaseAmountSetting) (float64, error) {
	if order.Payment == nil || order.Payment.State != orderconst.PaymentStateDone {
		return 0, nil
	}

	setting := commissionsetting.GetAmountSettingByTimestamp(roots, order.Order.CreateAt)
	if setting == nil {
		return 0, nil
	}
	rootPercent := int(setting.Percent)

	nextPercent := 0
	setting = commissionsetting.GetAmountSettingByTimestamp(nexts, order.Order.CreateAt)
	if setting != nil {
		nextPercent = int(setting.Percent)
	}

	if rootPercent < nextPercent {
		return 0, nil
	}

	orderAmount := order.Payment.Amount * order.Payment.CoinUSDCurrency
	rootAmount := orderAmount * float64(rootPercent-nextPercent) / 100.0

	logger.Sugar().Infof("order %v | %v root %v | %v next %v user %v",
		order.Order.ID, orderAmount, rootAmount, rootPercent, nextPercent, order.Order.UserID)

	return rootAmount, nil
}

func getPeriodRebate(ctx context.Context, appID, userID string, roots, nexts []*inspirepb.AppPurchaseAmountSetting) (float64, error) {
	orders, err := referral.GetOrders(ctx, appID, userID)
	if err != nil {
		return 0, xerrors.Errorf("fail get orders: %v", err)
	}

	totalRootAmount := 0.0

	for _, order := range orders {
		rootAmount, err := getOrderParentRebate(ctx, order, roots, nexts)
		if err != nil {
			return 0, xerrors.Errorf("fail get order rebate: %v", err)
		}

		totalRootAmount += rootAmount
	}

	return totalRootAmount, nil
}

func getIncomings(ctx context.Context, appID, userID string) (float64, error) {
	roots, err := commissionsetting.GetAmountSettingsByAppUser(ctx, appID, userID)
	if err != nil {
		return 0, xerrors.Errorf("fail get amount settings: %v", err)
	}

	invitees, err := referral.GetLayeredInvitees(ctx, appID, userID)
	if err != nil {
		return 0, xerrors.Errorf("fail get layered invitees: %v", err)
	}

	totalRootAmount := 0.0

	for _, iv := range invitees {
		nexts, err := commissionsetting.GetAmountSettingsByAppUser(ctx, iv.AppID, iv.InviteeID)
		if err != nil {
			return 0, xerrors.Errorf("fail get amount settings: %v", err)
		}

		rootAmount, err := getPeriodRebate(ctx, iv.AppID, iv.InviteeID, roots, nexts)
		if err != nil {
			return 0, xerrors.Errorf("fail get rebate: %v", err)
		}

		totalRootAmount += rootAmount
	}

	return totalRootAmount, nil
}

func GetSeparateIncoming(ctx context.Context, appID, userID string) (float64, error) {
	incoming, err := getRebate(ctx, appID, userID)
	if err != nil {
		return 0, xerrors.Errorf("fail get total incoming: %v", err)
	}

	subAmount, err := getIncomings(ctx, appID, userID)
	if err != nil {
		return 0, xerrors.Errorf("fail get sub incomings: %v", err)
	}

	incoming += subAmount

	return incoming, nil
}
