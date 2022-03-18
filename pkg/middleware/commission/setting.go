package commission

import (
	"context"
	"sort"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	cache "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/cache"
	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"

	"golang.org/x/xerrors"
)

func getCoins(ctx context.Context) ([]*inspirepb.CommissionCoinSetting, error) {
	return grpc2.GetCommissionCoinSettings(ctx, &inspirepb.GetCommissionCoinSettings{})
}

func uniqueSetting(ctx context.Context, appID string) (bool, error) {
	setting, err := grpc2.GetAppCommissionSettingByApp(ctx, &inspirepb.GetAppCommissionSettingByAppRequest{
		AppID: appID,
	})
	if err != nil {
		return false, xerrors.Errorf("fail get app commission setting")
	}

	return setting.UniqueSetting, nil
}

func getAmountSettingsByApp(ctx context.Context, appID string) (*inspirepb.AppPurchaseAmountSetting, error) {
	// TODO: for unique app commission setting
	return nil, nil
}

func getAmountSettingsByAppUser(ctx context.Context, appID, userID string) ([]*inspirepb.AppPurchaseAmountSetting, error) {
	cacheFor := "purchase:amount:settings"

	mySettings := cache.GetEntry(cacheKey(appID, userID, cacheFor))
	if settings != nil {
		return mySettings.([]*inspirepb.AppPurchaseAmountSetting), nil
	}

	settings, err := grpc2.GetAppPurchaseAmountSettingsByAppUser(ctx, &inspirepb.GetAppPurchaseAmountSettingsByAppUserRequest{
		AppID:  appID,
		UserID: userID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get app purchase amount setting: %v", err)
	}

	sort.Slice(settings, func(i, j) {
		return settings[i].Start < settings[j].Start
	})

	lastSetting := *inspirepb.AppPurchaseAmountSetting
	for _, setting := range settings {
		if lastSetting == nil {
			continue
		}
		if setting.Start != lastSetting.End {
			return nil, xerrors.Errorf("invalid purchase amount setting: %v", err)
		}
	}

	cache.AddEntry(cacheKey(appID, userID, cacheFor), settings)

	return settings, nil
}

func getAmountSettingByTimestamp(settings []*inspirepb.AppPurchaseAmountSetting, timestamp uint32) *inspirepb.AppPurchaseAmountSetting {
	var setting *inspirepb.AppPurchaseAmountSetting
	for _, s := range settings {
		if s.Start <= timestamp && (timestamp < s.End || s.End == 0) {
			return s
		}
	}
	return nil
}
