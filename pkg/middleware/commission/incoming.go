package commission

import (
	"context"

	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/referral/setting"

	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/commission/kpi"
	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/commission/separate"
	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/commission/unique"

	"golang.org/x/xerrors"
)

func getIncoming(ctx context.Context, appID, userID string) (float64, error) {
	_kpi, err := setting.KPISetting(ctx, appID)
	if err != nil {
		return 0, xerrors.Errorf("fail get kpi setting: %v", err)
	}

	if _kpi {
		return kpi.GetKPIIncoming(ctx, appID, userID)
	}

	_unique, err := setting.UniqueSetting(ctx, appID)
	if err != nil {
		return 0, xerrors.Errorf("fail get unique setting: %v", err)
	}

	if _unique {
		return unique.GetUniqueIncoming(ctx, appID, userID)
	}

	return separate.GetSeparateIncoming(ctx, appID, userID)
}
