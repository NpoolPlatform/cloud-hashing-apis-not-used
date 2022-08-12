package coupon

import (
	"context"
	"fmt"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"
	orderpb "github.com/NpoolPlatform/message/npool/cloud-hashing-order"

	orderconst "github.com/NpoolPlatform/cloud-hashing-order/pkg/const"
)

func GetCouponsByAppUser(ctx context.Context, in *npool.GetCouponsByAppUserRequest) (*npool.GetCouponsByAppUserResponse, error) {
	infos, err := grpc2.GetCouponsAllocatedByAppUser(ctx, &inspirepb.GetCouponsAllocatedDetailByAppUserRequest{
		AppID:  in.GetAppID(),
		UserID: in.GetUserID(),
	})
	if err != nil {
		return nil, fmt.Errorf("fail get coupons: %v", err)
	}

	coupons := []*npool.Coupon{}
	for _, info := range infos {
		order, err := grpc2.GetOrderByAppUserCouponTypeID(ctx, &orderpb.GetOrderByAppUserCouponTypeIDRequest{
			AppID:      in.AppID,
			UserID:     in.UserID,
			CouponType: info.Allocated.Type,
			CouponID:   info.Allocated.ID,
		})
		if err != nil {
			return nil, fmt.Errorf("fail get order for coupon: %v", err)
		}

		coupons = append(coupons, &npool.Coupon{
			Coupon: info,
			Order:  order,
		})
	}

	specials, err := grpc2.GetUserSpecialReductionsByAppUser(ctx, &inspirepb.GetUserSpecialReductionsByAppUserRequest{
		AppID:  in.GetAppID(),
		UserID: in.GetUserID(),
	})
	if err != nil {
		return nil, fmt.Errorf("fail get user special offer: %v", err)
	}

	offers := []*npool.UserSpecial{}
	for _, info := range specials {
		order, err := grpc2.GetOrderByAppUserCouponTypeID(ctx, &orderpb.GetOrderByAppUserCouponTypeIDRequest{
			AppID:      in.AppID,
			UserID:     in.UserID,
			CouponType: orderconst.UserSpecialReductionCoupon,
			CouponID:   info.ID,
		})
		if err != nil {
			return nil, fmt.Errorf("fail get order for coupon: %v", err)
		}
		offers = append(offers, &npool.UserSpecial{
			Coupon: info,
			Order:  order,
		})
	}

	return &npool.GetCouponsByAppUserResponse{
		Info: &npool.Coupons{
			Coupons: coupons,
			Offers:  offers,
		},
	}, nil
}
