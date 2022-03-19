package incoming

import (
	"context"
	"time"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	commissionsetting "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/commission/setting"
	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/referral"

	"golang.org/x/xerrors"
)

func getPeriodLevelAmount(ctx context.Context, appID, userID string, threshold float64, percent, start, end uint32) (float64, error) {
	if end == 0 {
		end = uint32(time.Now().Unix())
	}

	periodAmount, err := referral.GetPeriodUSDAmount(ctx, appID, userID, 0, end)
	if err != nil {
		return 0, xerrors.Errorf("fail get period usd amount: %v", err)
	}

	periodSubAmount, err := referral.GetPeriodSubUSDAmount(ctx, appID, userID, 0, end)
	if err != nil {
		return 0, xerrors.Errorf("fail get period sub usd amount: %v", err)
	}

	periodAmount += periodSubAmount
	if periodAmount < threshold {
		logger.Sugar().Warnf("threshold %v percent %v %v~%v period amount %v user %v",
			threshold, percent, start, end, periodAmount, userID)
		return 0, nil
	}

	periodAmount, err = referral.GetPeriodUSDAmount(ctx, appID, userID, start, end)
	if err != nil {
		return 0, xerrors.Errorf("fail get period usd amount: %v", err)
	}

	periodSubAmount, err = referral.GetPeriodSubUSDAmount(ctx, appID, userID, start, end)
	if err != nil {
		return 0, xerrors.Errorf("fail get period sub usd amount: %v", err)
	}

	periodAmount += periodSubAmount
	levelAmount := periodAmount * float64(percent) / 100.0

	logger.Sugar().Infof("threshold %v percent %v %v~%v level amount %v user %v period amount %v",
		threshold, percent, start, end, levelAmount, userID, periodAmount)

	return levelAmount, nil
}

func getLevelAmount(ctx context.Context, appID, userID string) (float64, error) {
	settings, err := commissionsetting.GetAmountSettingsByAppUser(ctx, appID, userID)
	if err != nil {
		return 0, xerrors.Errorf("fail get amount setting: %v", err)
	}

	levelAmount := 0.0

	for _, setting := range settings {
		amount, err := getPeriodLevelAmount(ctx, appID, userID, setting.Amount,
			setting.Percent, setting.Start, setting.End)
		if err != nil {
			return 0, xerrors.Errorf("fail get period level amount: %v", err)
		}
		levelAmount += amount
	}

	return levelAmount, nil
}

func getRootAmount(ctx context.Context, appID, userID string) (float64, error) {
	return getLevelAmount(ctx, appID, userID)
}

func GetKPIIncoming(ctx context.Context, appID, userID string) (float64, error) {
	rootAmount, err := getRootAmount(ctx, appID, userID)
	if err != nil {
		return 0, xerrors.Errorf("fail get root amount: %v", err)
	}

	invitees, err := referral.GetInvitees(ctx, appID, userID)
	if err != nil {
		return 0, xerrors.Errorf("fail get invitees: %v", err)
	}

	totalSubAmount := 0.0
	for _, iv := range invitees {
		subAmount, err := getRootAmount(ctx, iv.AppID, iv.InviteeID)
		if err != nil {
			return 0, xerrors.Errorf("fail get invitee amount: %v", err)
		}
		totalSubAmount += subAmount
	}

	return rootAmount - totalSubAmount, nil
}
