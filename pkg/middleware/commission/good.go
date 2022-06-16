package commission

import (
	"context"

	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/referral/setting"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	"golang.org/x/xerrors"
)

func getGoodCommissions(ctx context.Context, appID, userID string) ([]*npool.GoodCommission, error) {
	kpi, err := setting.KPISetting(ctx, appID)
	if err != nil {
		return nil, xerrors.Errorf("fail get kpi setting: %v", err)
	}

	if kpi {
		return getKPIGoodCommissions(ctx, appID, userID)
	}

	unique, err := setting.UniqueSetting(ctx, appID)
	if err != nil {
		return nil, xerrors.Errorf("fail get unique setting: %v", err)
	}

	if unique {
		return getUniqueGoodCommissions(ctx, appID, userID)
	}

	return getSeparateGoodCommissions(ctx, appID, userID)
}
