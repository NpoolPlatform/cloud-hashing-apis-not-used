package referral

import (
	"context"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	cache "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/cache"
	orderpb "github.com/NpoolPlatform/message/npool/cloud-hashing-order"

	"golang.org/x/xerrors"
)

const cacheOrders = "referral:orders"

func getOrders(ctx context.Context, appID, userID string) ([]*orderpb.OrderDetail, error) {
	myOrders := cache.GetEntry(cacheKey(appID, userID, cacheOrders))
	if myOrders != nil {
		return myOrders.([]*orderpb.OrderDetail), nil
	}

	// TODO: let database to sum orders amount
	orders, err := grpc2.GetOrdersDetailByAppUser(ctx, &orderpb.GetOrdersDetailByAppUserRequest{
		AppID:  appID,
		UserID: userID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get orders: %v", err)
	}

	cache.AddEntry(cacheKey(appID, userID, cacheOrders), orders)

	return orders, nil
}
