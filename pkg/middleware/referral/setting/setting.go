package setting

import (
	"context"
	"fmt"
	"sort"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	cache "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/cache"
	cachekey "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/referral/cachekey"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"
	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"

	"github.com/google/uuid"
)

func getCoins(ctx context.Context) ([]*inspirepb.CommissionCoinSetting, error) {
	return grpc2.GetCommissionCoinSettings(ctx, &inspirepb.GetCommissionCoinSettingsRequest{})
}

func GetUsingCoin(ctx context.Context) (*inspirepb.CommissionCoinSetting, error) {
	coins, err := getCoins(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail get coins: %v", err)
	}

	for _, coin := range coins {
		if coin.Using {
			return coin, nil
		}
	}

	return nil, fmt.Errorf("no using coin")
}

func UniqueSetting(ctx context.Context, appID string) (bool, error) {
	setting, err := grpc2.GetAppCommissionSettingByApp(ctx, &inspirepb.GetAppCommissionSettingByAppRequest{
		AppID: appID,
	})
	if err != nil || setting == nil {
		return false, fmt.Errorf("fail get app commission setting: %v", err)
	}

	return setting.UniqueSetting, nil
}

func KPISetting(ctx context.Context, appID string) (bool, error) {
	setting, err := grpc2.GetAppCommissionSettingByApp(ctx, &inspirepb.GetAppCommissionSettingByAppRequest{
		AppID: appID,
	})
	if err != nil || setting == nil {
		return false, fmt.Errorf("fail get app commission setting: %v", err)
	}

	return setting.KPISetting, nil
}

func getAmountSettingsByApp(ctx context.Context, appID string) (*inspirepb.AppPurchaseAmountSetting, error) { //nolint
	// TODO: for unique app commission setting
	return nil, nil
}

const cacheFor = "purchase:amount:settings"

func UpdateAmountSettingsCache(ctx context.Context, appID, userID string, settings []*inspirepb.AppPurchaseAmountSetting) {
	cache.AddEntry(cachekey.CacheKey(appID, userID, cacheFor), settings)
}

func GetAmountSettingsByAppUser(ctx context.Context, appID, userID string) ([]*inspirepb.AppPurchaseAmountSetting, error) {
	mySettings := cache.GetEntry(cachekey.CacheKey(appID, userID, cacheFor), func(data []byte) (interface{}, error) {
		return cache.UnmarshalAmountSettings(data)
	})
	if mySettings != nil {
		return mySettings.([]*inspirepb.AppPurchaseAmountSetting), nil
	}

	settings, err := grpc2.GetAppPurchaseAmountSettingsByAppUser(ctx, &inspirepb.GetAppPurchaseAmountSettingsByAppUserRequest{
		AppID:  appID,
		UserID: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("fail get app purchase amount setting: %v", err)
	}

	sort.Slice(settings, func(i, j int) bool {
		return settings[i].Start < settings[j].Start
	})
	sort.Slice(settings, func(i, j int) bool {
		return settings[i].Amount < settings[j].Amount
	})

	var lastSetting *inspirepb.AppPurchaseAmountSetting
	for _, setting := range settings {
		if lastSetting == nil {
			continue
		}
		if setting.Start != lastSetting.End {
			return nil, fmt.Errorf("invalid purchase amount setting: %v", err)
		}
	}

	UpdateAmountSettingsCache(ctx, appID, userID, settings)

	return settings, nil
}

func GetOrderAmountSetting(settings []*inspirepb.AppPurchaseAmountSetting, order *npool.Order) *inspirepb.AppPurchaseAmountSetting {
	invalidID := uuid.UUID{}.String()
	for _, s := range settings {
		if s.Amount > 0 {
			continue
		}
		if s.Start <= order.Order.Order.CreateAt && (order.Order.Order.CreateAt < s.End || s.End == 0) {
			if s.GoodID == invalidID || s.GoodID == "" || s.GoodID == order.Order.Order.GoodID {
				return s
			}
		}
	}
	return nil
}

func GetGoodAmountSetting(settings []*inspirepb.AppPurchaseAmountSetting, goodID string) *inspirepb.AppPurchaseAmountSetting {
	invalidID := uuid.UUID{}.String()
	for _, s := range settings {
		if s.Amount > 0 {
			continue
		}
		if s.End == 0 {
			if s.GoodID == invalidID || s.GoodID == "" || s.GoodID == goodID {
				return s
			}
		}
	}
	return nil
}
