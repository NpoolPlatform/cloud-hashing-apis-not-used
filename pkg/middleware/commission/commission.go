package commission

import (
	"context"
)

func GetCommission(ctx context.Context, appID, userID string) (float64, error) {
	return getIncoming(ctx, appID, userID)
}
