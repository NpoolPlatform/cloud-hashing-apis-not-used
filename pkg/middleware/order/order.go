package order

import (
	"context"
	"fmt"
	"sort"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	cache "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/cache"
	gooddetail "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/good"

	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	billingpb "github.com/NpoolPlatform/message/npool/cloud-hashing-billing"
	goodspb "github.com/NpoolPlatform/message/npool/cloud-hashing-goods"
	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"
	orderpb "github.com/NpoolPlatform/message/npool/cloud-hashing-order"
	coininfopb "github.com/NpoolPlatform/message/npool/coininfo"

	orderconst "github.com/NpoolPlatform/cloud-hashing-order/pkg/const"

	"github.com/google/uuid"

	appgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/appgood"
	appgoodmgrpb "github.com/NpoolPlatform/message/npool/good/mgr/v1/appgood"
)

const (
	secondsInDay          = 24 * 60 * 60
	cacheGood             = "order:good"
	cacheCoin             = "order:coin"
	cacheFixAmount        = "order:fix-amount"
	cacheDiscount         = "order:discount"
	cacheSpecialOffer     = "order:user-special-offer"
	cacheAppGood          = "order:app-good"
	cacheAppGoodPromotion = "order:app-good-promotion"
)

func cacheKey(key, id string) string {
	return fmt.Sprintf("%v:%v", key, id)
}

func constructOrder(
	info *orderpb.OrderDetail,
	goodInfo *npool.Good,
	coupon *inspirepb.CouponAllocatedDetail,
	paymentCoinInfo *coininfopb.CoinInfo,
	account *billingpb.CoinAccountInfo,
	discountCoupon *inspirepb.CouponAllocatedDetail,
	userSpecial *inspirepb.UserSpecialReduction,
	appGood *goodspb.AppGoodInfo,
	promotion *goodspb.AppGoodPromotion,
) *npool.Order {
	return &npool.Order{
		Order:                info,
		PayToAccount:         account,
		Good:                 goodInfo,
		PayWithCoin:          paymentCoinInfo,
		FixAmountCoupon:      coupon,
		DiscountCoupon:       discountCoupon,
		UserSpecialReduction: userSpecial,
		AppGood:              appGood,
		Promotion:            promotion,
		PaymentDeadline:      info.Order.CreateAt + orderconst.TimeoutSeconds,
	}
}

func expandOrder(ctx context.Context, info *orderpb.OrderDetail, base bool) (*npool.Order, error) { //nolint
	var coupon *inspirepb.CouponAllocatedDetail
	var err error
	invalidUUID := uuid.UUID{}.String()

	if !base && info.Order.CouponID != invalidUUID {
		couponAllocated := cache.GetEntry(cacheKey(cacheFixAmount, info.Order.CouponID), func(data []byte) (interface{}, error) {
			return cache.UnmarshalCouponAllocated(data)
		})
		if couponAllocated != nil {
			coupon = couponAllocated.(*inspirepb.CouponAllocatedDetail) //nolint
		}

		if coupon == nil {
			order, err := grpc2.GetOrderByAppUserCouponTypeID(ctx, &orderpb.GetOrderByAppUserCouponTypeIDRequest{
				AppID:      info.Order.AppID,
				UserID:     info.Order.UserID,
				CouponType: orderconst.FixAmountCoupon,
				CouponID:   info.Order.CouponID,
			})
			if err != nil {
				return nil, fmt.Errorf("fail check coupon usage: %v", err)
			}

			if order != nil && order.ID != info.Order.ID {
				return nil, fmt.Errorf("fail check coupon usage")
			}

			couponAllocated, err := grpc2.GetCouponAllocated(ctx, &inspirepb.GetCouponAllocatedDetailRequest{
				ID: info.Order.CouponID,
			})
			if err != nil && info.Order.CouponID != invalidUUID {
				return nil, fmt.Errorf("fail get coupon allocated detail: %v", err)
			}

			if couponAllocated != nil {
				if couponAllocated.Allocated.AppID != info.Order.AppID || couponAllocated.Allocated.UserID != info.Order.UserID {
					return nil, fmt.Errorf("invalid coupon")
				}
				coupon = couponAllocated
				cache.AddEntry(cacheKey(cacheFixAmount, info.Order.CouponID), coupon)
			}
		}
	}

	var paymentCoinInfo *coininfopb.CoinInfo
	var accountInfo *billingpb.CoinAccountInfo

	if info.Payment != nil {
		coinInfo := cache.GetEntry(cacheKey(cacheCoin, info.Payment.CoinInfoID), func(data []byte) (interface{}, error) {
			return cache.UnmarshalCoinInfo(data)
		})
		if coinInfo != nil {
			paymentCoinInfo = coinInfo.(*coininfopb.CoinInfo) //nolint
		}

		if paymentCoinInfo == nil {
			paymentCoinInfo, err := grpc2.GetCoinInfo(ctx, &coininfopb.GetCoinInfoRequest{
				ID: info.Payment.CoinInfoID,
			})
			if err != nil || paymentCoinInfo == nil {
				return nil, fmt.Errorf("fail get payment coin info: %v", err)
			}

			cache.AddEntry(cacheKey(cacheCoin, info.Payment.CoinInfoID), paymentCoinInfo)
		}

		if !base {
			account, err := grpc2.GetBillingAccount(ctx, &billingpb.GetCoinAccountRequest{
				ID: info.Payment.AccountID,
			})
			if err != nil {
				return nil, fmt.Errorf("fail get payment address: %v", err)
			}
			accountInfo = account
		}
	}

	var discountCoupon *inspirepb.CouponAllocatedDetail

	if !base && info.Order.DiscountCouponID != invalidUUID {
		discount := cache.GetEntry(cacheKey(cacheDiscount, info.Order.DiscountCouponID), func(data []byte) (interface{}, error) {
			return cache.UnmarshalCouponAllocated(data)
		})
		if discount != nil {
			discountCoupon = discount.(*inspirepb.CouponAllocatedDetail) //nolint
		}

		if discountCoupon == nil {
			order, err := grpc2.GetOrderByAppUserCouponTypeID(ctx, &orderpb.GetOrderByAppUserCouponTypeIDRequest{
				AppID:      info.Order.AppID,
				UserID:     info.Order.UserID,
				CouponType: orderconst.DiscountCoupon,
				CouponID:   info.Order.DiscountCouponID,
			})
			if err != nil {
				return nil, fmt.Errorf("fail check coupon usage: %v", err)
			}

			if order != nil && order.ID != info.Order.ID {
				return nil, fmt.Errorf("fail check coupon usage")
			}

			discount, err := grpc2.GetCouponAllocated(ctx, &inspirepb.GetCouponAllocatedDetailRequest{
				ID: info.Order.DiscountCouponID,
			})
			if err != nil {
				return nil, fmt.Errorf("fail get discount coupon allocated detail: %v", err)
			}

			if discount != nil {
				if discount.Allocated.AppID != info.Order.AppID || discount.Allocated.UserID != info.Order.UserID {
					return nil, fmt.Errorf("invalid coupon")
				}
				discountCoupon = discount
				cache.AddEntry(cacheKey(cacheDiscount, info.Order.DiscountCouponID), discountCoupon)
			}
		}
	}

	var userSpecialReduction *inspirepb.UserSpecialReduction

	if !base && info.Order.UserSpecialReductionID != invalidUUID {
		special := cache.GetEntry(cacheKey(cacheSpecialOffer, info.Order.UserSpecialReductionID), func(data []byte) (interface{}, error) {
			return cache.UnmarshalSpecialOffer(data)
		})
		if special != nil {
			userSpecialReduction = special.(*inspirepb.UserSpecialReduction) //nolint
		}

		if userSpecialReduction == nil {
			order, err := grpc2.GetOrderByAppUserCouponTypeID(ctx, &orderpb.GetOrderByAppUserCouponTypeIDRequest{
				AppID:      info.Order.AppID,
				UserID:     info.Order.UserID,
				CouponType: orderconst.UserSpecialReductionCoupon,
				CouponID:   info.Order.UserSpecialReductionID,
			})
			if err != nil {
				return nil, fmt.Errorf("fail check coupon usage: %v", err)
			}

			if order != nil && order.ID != info.Order.ID {
				return nil, fmt.Errorf("fail check coupon usage")
			}

			userSpecial, err := grpc2.GetUserSpecialReduction(ctx, &inspirepb.GetUserSpecialReductionRequest{
				ID: info.Order.UserSpecialReductionID,
			})
			if err != nil {
				return nil, fmt.Errorf("fail get user special reduction: %v", err)
			}

			if userSpecial != nil {
				if userSpecial.AppID != info.Order.AppID || userSpecial.UserID != info.Order.UserID {
					return nil, fmt.Errorf("invalid coupon")
				}
				userSpecialReduction = userSpecial
				cache.AddEntry(cacheKey(cacheDiscount, info.Order.UserSpecialReductionID), userSpecialReduction)
			}
		}
	}

	var goodInfo *npool.Good

	cacheGoodInfo := cache.GetEntry(cacheKey(cacheGood, info.Order.GetGoodID()), func(data []byte) (interface{}, error) {
		return cache.UnmarshalGood(data)
	})
	if cacheGoodInfo != nil {
		goodInfo = cacheGoodInfo.(*npool.Good) //nolint
	}

	if goodInfo == nil {
		resp, err := gooddetail.Get(ctx, &npool.GetGoodRequest{
			ID: info.Order.GetGoodID(),
		})
		if err != nil || resp.Info == nil {
			return nil, fmt.Errorf("fail get good info: %v", err)
		}
		goodInfo = resp.Info
		cache.AddEntry(cacheKey(cacheGood, info.Order.GetGoodID()), goodInfo)
	}

	var appGood *goodspb.AppGoodInfo

	cacheAppGoodInfo := cache.GetEntry(cacheKey(cacheAppGood, fmt.Sprintf(info.Order.GetAppID(), info.Order.GetGoodID())), func(data []byte) (interface{}, error) {
		return cache.UnmarshalAppGoodInfo(data)
	})
	if cacheAppGoodInfo != nil {
		appGood = cacheAppGoodInfo.(*goodspb.AppGoodInfo) //nolint
	}

	if appGood == nil {
		appgoodmwcli.GetGoods(ctx, &appgoodmgrpb.Conds{
			ID:      nil,
			AppID:   nil,
			GoodID:  nil,
			GoodIDs: nil,
		}, 0, 0)

		appGood, err = grpc2.GetAppGoodByAppGood(ctx, &goodspb.GetAppGoodByAppGoodRequest{
			AppID:  info.Order.GetAppID(),
			GoodID: info.Order.GetGoodID(),
		})
		if err != nil {
			return nil, fmt.Errorf("fail get app good: %v", err)
		}
		cache.AddEntry(cacheKey(cacheAppGood, fmt.Sprintf(info.Order.GetAppID(), info.Order.GetGoodID())), appGood)
	}

	var promotion *goodspb.AppGoodPromotion

	cachedAppGoodPromotion := cache.GetEntry(cacheKey(cacheAppGoodPromotion, info.Order.PromotionID), func(data []byte) (interface{}, error) {
		return cache.UnmarshalAppGoodPromotion(data)
	})
	if cachedAppGoodPromotion != nil {
		promotion = cachedAppGoodPromotion.(*goodspb.AppGoodPromotion) //nolint
	}

	if promotion == nil {
		promotion, err = grpc2.GetAppGoodPromotion(ctx, &goodspb.GetAppGoodPromotionRequest{
			ID: info.Order.PromotionID,
		})
		if err != nil {
			return nil, fmt.Errorf("fail get promotion: %v", err)
		}
		if promotion != nil {
			if info.Order.GetAppID() != promotion.AppID || info.Order.GetGoodID() != promotion.GoodID {
				return nil, fmt.Errorf("invalid promotion")
			}
		}
		cache.AddEntry(cacheKey(cacheAppGoodPromotion, info.Order.PromotionID), promotion)
	}

	return constructOrder(
		info,
		goodInfo,
		coupon,
		paymentCoinInfo,
		accountInfo,
		discountCoupon,
		userSpecialReduction,
		appGood,
		promotion,
	), nil
}

func GetOrdersByAppUser(ctx context.Context, in *npool.GetOrdersByAppUserRequest) (*npool.GetOrdersByAppUserResponse, error) {
	orders, err := grpc2.GetOrdersDetailByAppUser(ctx, &orderpb.GetOrdersDetailByAppUserRequest{
		AppID:  in.GetAppID(),
		UserID: in.GetUserID(),
	})
	if err != nil {
		return nil, fmt.Errorf("fail get orders detail by app user: %v", err)
	}

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].Order.CreateAt > orders[j].Order.CreateAt
	})

	details := []*npool.Order{}
	for _, info := range orders {
		detail, err := expandOrder(ctx, info, false)
		if err != nil {
			logger.Sugar().Warnf("cannot expand order detail: %v", err)
			continue
		}

		details = append(details, detail)
	}

	return &npool.GetOrdersByAppUserResponse{
		Infos: details,
	}, nil
}
