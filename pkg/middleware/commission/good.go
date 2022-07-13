package commission

import (
	"context"

	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/referral/setting"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/commission/kpi"
	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/commission/separate"
	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/commission/unique"

	"golang.org/x/xerrors"
)

func getGoodCommissions(ctx context.Context, appID, userID string) ([]*npool.GoodCommission, error) {
	_kpi, err := setting.KPISetting(ctx, appID)
	if err != nil {
		return nil, xerrors.Errorf("fail get kpi setting: %v", err)
	}

	if _kpi {
		return kpi.GetKPIGoodCommissions(ctx, appID, userID)
	}

	_unique, err := setting.UniqueSetting(ctx, appID)
	if err != nil {
		return nil, xerrors.Errorf("fail get unique setting: %v", err)
	}

	if _unique {
		return unique.GetUniqueGoodCommissions(ctx, appID, userID)
	}

	return separate.GetSeparateGoodCommissions(ctx, appID, userID)
}
