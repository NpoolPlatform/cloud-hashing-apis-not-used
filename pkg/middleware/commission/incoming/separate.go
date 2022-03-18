package incoming

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	commissionsetting "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/commission/setting"
	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/referral"
	orderconst "github.com/NpoolPlatform/cloud-hashing-order/pkg/const"
	orderpb "github.com/NpoolPlatform/message/npool/cloud-hashing-order"

	"golang.org/x/xerrors"
)

func getRebate(ctx context.Context, appID, userID string) (float64, error) {
	settings, err := commissionsetting.GetAmountSettingsByAppUser(ctx, appID, userID)
	if err != nil {
		return 0, xerrors.Errorf("fail get amount settings: %v", err)
	}

	totalAmount := 0.0
	logger.Sugar().Infof("settings %v", settings)

	for _, setting := range settings {
		amount, err := referral.GetPeriodUSDAmount(ctx, appID, userID, setting.Start, setting.End)
		if err != nil {
			return 0, xerrors.Errorf("fail get period usd amount: %v", err)
		}
		totalAmount += amount * float64(setting.Percent) / 100.0

		logger.Sugar().Infof("commission %v %v amount %v %v %v", appID, userID, totalAmount, amount, setting.Percent)
	}

	return totalAmount, nil
}

func getOrderParentRebate(ctx context.Context, order *orderpb.OrderDetail) (float64, error) {
	inviter, err := referral.GetInviter(ctx, order.Order.AppID, order.Order.UserID)
	if err != nil {
		return 0, xerrors.Errorf("fail get inviter: %v", err)
	}

	parents, err := commissionsetting.GetAmountSettingsByAppUser(ctx, inviter.AppID, inviter.InviterID)
	if err != nil {
		return 0, xerrors.Errorf("fail get parent settings: %v", err)
	}

	childs, err := commissionsetting.GetAmountSettingsByAppUser(ctx, inviter.AppID, inviter.InviteeID)
	if err != nil {
		return 0, xerrors.Errorf("fail get child settings: %v", err)
	}

	if order.Payment == nil || order.Payment.State != orderconst.PaymentStateDone {
		return 0, nil
	}

	setting := commissionsetting.GetAmountSettingByTimestamp(parents, order.Order.CreateAt)
	if setting == nil {
		return 0, nil
	}
	parentPercent := int(setting.Percent)

	childPercent := 0
	setting = commissionsetting.GetAmountSettingByTimestamp(childs, order.Order.CreateAt)
	if setting != nil {
		childPercent = int(setting.Percent)
	}

	if parentPercent < childPercent {
		return 0, nil
	}

	orderAmount := order.Payment.Amount * order.Payment.CoinUSDCurrency
	parentAmount := orderAmount * float64(parentPercent) / 100.0

	logger.Sugar().Infof("parent commission %v %v amount %v %v %v",
		order.Order.AppID, order.Order.UserID, orderAmount, parentAmount, parentPercent)

	return parentAmount, nil
}

func getPeriodRebate(ctx context.Context, appID, userID string) (float64, error) {
	orders, err := referral.GetOrders(ctx, appID, userID)
	if err != nil {
		return 0, xerrors.Errorf("fail get orders: %v", err)
	}

	totalParentAmount := 0.0

	for _, order := range orders {
		parentAmount, err := getOrderParentRebate(ctx, order)
		if err != nil {
			return 0, xerrors.Errorf("fail get order rebate: %v", err)
		}

		totalParentAmount += parentAmount
	}

	return totalParentAmount, nil
}

func getIncomings(ctx context.Context, appID, userID string) (float64, error) {
	invitees, err := referral.GetLayeredInvitees(ctx, appID, userID)
	if err != nil {
		return 0, xerrors.Errorf("fail get layered invitees: %v", err)
	}

	rebates := map[string]float64{}

	for _, iv := range invitees {
		amount, err := getPeriodRebate(ctx, iv.AppID, iv.InviteeID)
		if err != nil {
			return 0, xerrors.Errorf("fail get rebate: %v", err)
		}

		rbAmount := rebates[iv.InviteeID]
		rbAmount += amount

		rebates[iv.InviteeID] = rbAmount
	}

	ivs, err := referral.GetInvitees(ctx, appID, userID)
	if err != nil {
		return 0, xerrors.Errorf("fail get invitees: %v", err)
	}

	totalAmount := 0.0
	for _, iv := range ivs {
		totalAmount += rebates[iv.InviteeID]
	}

	return totalAmount, nil
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
