package referral

import (
	"context"

	cache "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/cache"
	ordermw "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/order"
	cachekey "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/referral/cachekey"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	"golang.org/x/xerrors"
)

const cacheOrders = "referral:orders"

func GetOrders(ctx context.Context, appID, userID string) ([]*npool.Order, error) {
	myOrders := cache.GetEntry(cachekey.CacheKey(appID, userID, cacheOrders))
	if myOrders != nil {
		return myOrders.([]*npool.Order), nil
	}

	// TODO: let database to sum orders amount
	orders, err := ordermw.GetOrdersByAppUser(ctx, &npool.GetOrdersByAppUserRequest{
		AppID:  appID,
		UserID: userID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get orders: %v", err)
	}

	cache.AddEntry(cachekey.CacheKey(appID, userID, cacheOrders), orders.Infos)

	return orders.Infos, nil
}
