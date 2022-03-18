package commission

import (
	"context"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
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
