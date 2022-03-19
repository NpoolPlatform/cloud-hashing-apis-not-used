package commission

import (
	"context"

	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/commission/incoming"
	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/commission/setting"

	"golang.org/x/xerrors"
)

func getIncoming(ctx context.Context, appID, userID string) (float64, error) {
	kpi, err := setting.KPISetting(ctx, appID)
	if err != nil {
		return 0, xerrors.Errorf("fail get kpi setting: %v", err)
	}

	if kpi {
		return incoming.GetKPIIncoming(ctx, appID, userID)
	}

	unique, err := setting.UniqueSetting(ctx, appID)
	if err != nil {
		return 0, xerrors.Errorf("fail get unique setting: %v", err)
	}

	if unique {
		return incoming.GetUniqueIncoming(ctx, appID, userID)
	}

	return incoming.GetSeparateIncoming(ctx, appID, userID)
}
