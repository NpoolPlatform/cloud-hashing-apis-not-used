package coupon

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"
	orderpb "github.com/NpoolPlatform/message/npool/order/mw/v1/order"

	npoolpb "github.com/NpoolPlatform/message/npool"
	ordercli "github.com/NpoolPlatform/order-middleware/pkg/client/order"

	inspireconst "github.com/NpoolPlatform/cloud-hashing-inspire/pkg/const"
)

const (
	FixAmountCoupon            = inspireconst.CouponTypeCoupon
	DiscountCoupon             = inspireconst.CouponTypeDiscount
	UserSpecialReductionCoupon = "user-special-reduction"
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
		conds := &orderpb.Conds{
			AppID: &npoolpb.StringVal{
				Op:    cruder.EQ,
				Value: in.AppID,
			},
			UserID: &npoolpb.StringVal{
				Op:    cruder.EQ,
				Value: in.UserID,
			},
		}

		switch info.Allocated.Type {
		case FixAmountCoupon:
			conds.FixAmountCouponID = &npoolpb.StringVal{
				Op:    cruder.EQ,
				Value: info.Allocated.ID,
			}
		case DiscountCoupon:
			conds.DiscountCouponID = &npoolpb.StringVal{
				Op:    cruder.EQ,
				Value: info.Allocated.ID,
			}
		case UserSpecialReductionCoupon:
			conds.UserSpecialReductionID = &npoolpb.StringVal{
				Op:    cruder.EQ,
				Value: info.Allocated.ID,
			}
		}

		order, err := ordercli.GetOrderOnly(ctx, conds)
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
		order, err := ordercli.GetOrderOnly(ctx, &orderpb.Conds{
			AppID: &npoolpb.StringVal{
				Op:    cruder.EQ,
				Value: in.AppID,
			},
			UserID: &npoolpb.StringVal{
				Op:    cruder.EQ,
				Value: in.UserID,
			},
			UserSpecialReductionID: &npoolpb.StringVal{
				Op:    cruder.EQ,
				Value: info.ID,
			},
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
