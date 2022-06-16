package commission

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"
)

func GetCommission(ctx context.Context, appID, userID string) (float64, error) {
	return getIncoming(ctx, appID, userID)
}

func GetGoodCommissions(ctx context.Context, appID, userID string, directOnly bool) ([]*npool.GoodCommission, error) {
	return nil, nil
}
