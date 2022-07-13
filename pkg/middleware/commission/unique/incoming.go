package unique

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	"golang.org/x/xerrors"
)

func GetUniqueIncoming(ctx context.Context, appID, userID string) (float64, error) {
	return 0, xerrors.Errorf("NOT IMPLEMENTED")
}

func GetUniqueGoodCommissions(ctx context.Context, appID, userID string) ([]*npool.GoodCommission, error) {
	return nil, xerrors.Errorf("NOT IMPLEMENTED")
}
