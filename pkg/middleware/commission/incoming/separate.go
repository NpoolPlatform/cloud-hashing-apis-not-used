package incoming

import (
	"context"
	"time"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	commissionsetting "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/commission/setting"
	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/referral"
	orderconst "github.com/NpoolPlatform/cloud-hashing-order/pkg/const"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"
	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"

	"golang.org/x/xerrors"
)

func dayBeginning() uint32 {
	return uint32(time.Now().Unix() / 24 / 60 / 60 * 24 * 60 * 60)
}

func getRebate(ctx context.Context, appID, userID string) (float64, error) {
	settings, err := commissionsetting.GetAmountSettingsByAppUser(ctx, appID, userID)
	if err != nil {
		return 0, xerrors.Errorf("fail get amount settings: %v", err)
	}

	totalAmount := 0.0

	for _, setting := range settings {
		if setting.Start >= dayBeginning() {
			continue
		}

		end := setting.End
		if end == 0 || end >= dayBeginning() {
			end = dayBeginning()
		}

		amount, err := referral.GetPeriodUSDAmount(ctx, appID, userID, setting.Start, end)
		if err != nil {
			return 0, xerrors.Errorf("fail get period usd amount: %v", err)
		}
		totalAmount += amount * float64(setting.Percent) / 100.0

		logger.Sugar().Infof("amount %v percent %v user %v", amount, setting.Percent, userID)
	}

	return totalAmount, nil
}

func getOrderParentRebate(_ context.Context, order *npool.Order, roots, nexts []*inspirepb.AppPurchaseAmountSetting) float64 {
	if order.Order.Payment == nil || order.Order.Payment.State != orderconst.PaymentStateDone {
		return 0
	}

	if order.Order.Order.CreateAt > dayBeginning() {
		return 0
	}

	setting := commissionsetting.GetAmountSettingByTimestamp(roots, order.Order.Order.CreateAt)
	if setting == nil {
		return 0
	}
	rootPercent := int(setting.Percent)

	nextPercent := 0
	setting = commissionsetting.GetAmountSettingByTimestamp(nexts, order.Order.Order.CreateAt)
	if setting != nil {
		nextPercent = int(setting.Percent)
	}

	if rootPercent < nextPercent {
		return 0
	}

	orderAmount := order.Order.Payment.Amount * order.Order.Payment.CoinUSDCurrency
	rootAmount := orderAmount * float64(rootPercent-nextPercent) / 100.0

	logger.Sugar().Infof("sub order %v | %v root %v | %v next %v user %v",
		order.Order.Order.ID, orderAmount, rootAmount, rootPercent, nextPercent,
		order.Order.Order.UserID)

	return rootAmount
}

func getPeriodRebate(ctx context.Context, appID, userID string, roots, nexts []*inspirepb.AppPurchaseAmountSetting) (float64, error) {
	orders, err := referral.GetOrders(ctx, appID, userID)
	if err != nil {
		return 0, xerrors.Errorf("fail get orders: %v", err)
	}

	totalRootAmount := 0.0

	for _, order := range orders {
		totalRootAmount += getOrderParentRebate(ctx, order, roots, nexts)
	}

	return totalRootAmount, nil
}

func findRootInviter(_ context.Context, rootUserID, inviteeID string, invitees []*inspirepb.RegistrationInvitation) (string, error) {
getInviter:
	for _, iv := range invitees {
		if iv.InviteeID == inviteeID {
			if iv.InviterID == rootUserID {
				return iv.InviteeID, nil
			}
			inviteeID = iv.InviterID
			goto getInviter
		}
	}
	return "", xerrors.Errorf("cannot find root inviter")
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
		inviteeID, err := findRootInviter(ctx, userID, iv.InviteeID, invitees)
		if err != nil {
			return 0, xerrors.Errorf("fail find root inviter: %v", err)
		}

		nexts, err := commissionsetting.GetAmountSettingsByAppUser(ctx, iv.AppID, inviteeID)
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
