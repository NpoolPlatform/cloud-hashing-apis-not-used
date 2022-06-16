package commission

import (
	"context"

	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"
)

func GetAmountSettings(ctx context.Context, appID, userID string) ([]*inspirepb.AppPurchaseAmountSetting, error) {
	return nil, nil
}

func CreateAmountSetting(ctx context.Context, appID, userID, targetUserID string, setting *inspirepb.AppPurchaseAmountSetting) ([]*inspirepb.AppPurchaseAmountSetting, error) {
	return nil, nil
}
