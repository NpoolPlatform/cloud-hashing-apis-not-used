package incoming

import (
	"context"
	"fmt"
	"sort"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	commissionsetting "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/commission/setting"
	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/referral"
)

func getAmount(ctx context.Context, appID, userID string) (float64, error) {
	settings, err := commissionsetting.GetAmountSettingsByAppUser(ctx, appID, userID)
	if err != nil {
		return 0, fmt.Errorf("fail get level settings: %v", err)
	}

	if len(settings) == 0 {
		return 0, nil
	}

	sort.Slice(settings, func(i, j int) bool {
		return settings[i].Start < settings[j].Start
	})
	sort.Slice(settings, func(i, j int) bool {
		return settings[i].Amount < settings[j].Amount
	})

	levelAmount := 0.0
	lastAmount := 0.0
	lastPercent := uint32(0)

	totalAmount, err := referral.GetUSDAmount(ctx, appID, userID)
	if err != nil {
		return 0, fmt.Errorf("fail get usd amount: %v", err)
	}

	subAmount, err := referral.GetSubUSDAmount(ctx, appID, userID)
	if err != nil {
		return 0, fmt.Errorf("fail get sub usd amount: %v", err)
	}

	totalAmount += subAmount
	remainAmount := totalAmount

	for i, setting := range settings {
		if setting.End > 0 {
			continue
		}

		if i == 0 && setting.Amount == 0 {
			lastAmount = setting.Amount
			lastPercent = setting.Percent
			continue
		}

		if remainAmount == 0 {
			break
		}

		amount := 0.0
		if totalAmount < setting.Amount {
			amount = totalAmount - lastAmount
		} else {
			amount = setting.Amount - lastAmount
		}

		levelAmount += amount * float64(lastPercent) / 100.0
		logger.Sugar().Infof("amount %v level amount %v sub amount %v total amount %v last amount %v last percent %v user %v",
			amount, levelAmount, subAmount, totalAmount, lastAmount, lastPercent, userID)

		lastAmount = setting.Amount
		lastPercent = setting.Percent
		remainAmount = totalAmount - setting.Amount
	}

	levelAmount += remainAmount * float64(lastPercent) / 100.0
	logger.Sugar().Infof("remain amount %v level amount %v sub amount %v total amount %v last amount %v last percent %v user %v",
		remainAmount, levelAmount, subAmount, totalAmount, lastAmount, lastPercent, userID)

	return levelAmount, nil
}

func getRootAmount(ctx context.Context, appID, userID string) (float64, error) {
	return getAmount(ctx, appID, userID)
}

func GetKPIIncoming(ctx context.Context, appID, userID string) (float64, error) {
	rootAmount, err := getRootAmount(ctx, appID, userID)
	if err != nil {
		return 0, fmt.Errorf("fail get root amount: %v", err)
	}

	invitees, err := referral.GetInvitees(ctx, appID, userID)
	if err != nil {
		return 0, fmt.Errorf("fail get invitees: %v", err)
	}

	totalSubAmount := 0.0
	for _, iv := range invitees {
		subAmount, err := getRootAmount(ctx, iv.AppID, iv.InviteeID)
		if err != nil {
			return 0, fmt.Errorf("fail get invitee amount: %v", err)
		}
		totalSubAmount += subAmount
	}

	logger.Sugar().Infof("root amount %v sub amount %v user %v", rootAmount, totalSubAmount, userID)

	return rootAmount - totalSubAmount, nil
}
