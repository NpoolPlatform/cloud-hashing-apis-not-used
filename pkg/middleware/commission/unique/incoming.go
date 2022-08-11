package unique

import (
	"context"
	"fmt"

	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"
)

func GetUniqueIncoming(ctx context.Context, appID, userID string) (float64, error) {
	return 0, fmt.Errorf("NOT IMPLEMENTED")
}

func GetUniqueGoodCommissions(ctx context.Context, appID, userID string) ([]*npool.GoodCommission, error) {
	return nil, fmt.Errorf("NOT IMPLEMENTED")
}
