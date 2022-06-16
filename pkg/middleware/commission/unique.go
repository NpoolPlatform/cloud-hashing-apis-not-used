package commission

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	"golang.org/x/xerrors"
)

func getUniqueIncoming(ctx context.Context, appID, userID string) (float64, error) {
	return 0, xerrors.Errorf("NOT IMPLEMENTED")
}

func getUniqueGoodCommissions(ctx context.Context, appID, userID string) ([]*npool.GoodCommission, error) {
	return nil, xerrors.Errorf("NOT IMPLEMENTED")
}
