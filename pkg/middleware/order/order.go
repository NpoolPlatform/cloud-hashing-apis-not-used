package order

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	constant "github.com/NpoolPlatform/cloud-hashing-apis/pkg/const"
	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	cache "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/cache"
	fee "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/fee"
	gooddetail "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/good"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	oraclecli "github.com/NpoolPlatform/oracle-manager/pkg/client"
	oracleconst "github.com/NpoolPlatform/oracle-manager/pkg/const"
	currency "github.com/NpoolPlatform/oracle-manager/pkg/middleware/currency"

	billingpb "github.com/NpoolPlatform/message/npool/cloud-hashing-billing"
	goodspb "github.com/NpoolPlatform/message/npool/cloud-hashing-goods"
	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"
	orderpb "github.com/NpoolPlatform/message/npool/cloud-hashing-order"
	coininfopb "github.com/NpoolPlatform/message/npool/coininfo"
	sphinxproxypb "github.com/NpoolPlatform/message/npool/sphinxproxy"

	orderconst "github.com/NpoolPlatform/cloud-hashing-order/pkg/const"
	accountlock "github.com/NpoolPlatform/staker-manager/pkg/middleware/account"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"google.golang.org/protobuf/types/known/structpb"

	stockcli "github.com/NpoolPlatform/stock-manager/pkg/client"
	stockconst "github.com/NpoolPlatform/stock-manager/pkg/const"

	"github.com/google/uuid"

	"golang.org/x/xerrors"
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
		couponAllocated := cache.GetEntry(cacheKey(cacheFixAmount, info.Order.CouponID))
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
				return nil, xerrors.Errorf("fail check coupon usage: %v", err)
			}

			if order != nil && order.ID != info.Order.ID {
				return nil, xerrors.Errorf("fail check coupon usage")
			}

			couponAllocated, err := grpc2.GetCouponAllocated(ctx, &inspirepb.GetCouponAllocatedDetailRequest{
				ID: info.Order.CouponID,
			})
			if err != nil && info.Order.CouponID != invalidUUID {
				return nil, xerrors.Errorf("fail get coupon allocated detail: %v", err)
			}

			if couponAllocated != nil {
				if couponAllocated.Allocated.AppID != info.Order.AppID || couponAllocated.Allocated.UserID != info.Order.UserID {
					return nil, xerrors.Errorf("invalid coupon")
				}
				coupon = couponAllocated
				cache.AddEntry(cacheKey(cacheFixAmount, info.Order.CouponID), coupon)
			}
		}
	}

	var paymentCoinInfo *coininfopb.CoinInfo
	var accountInfo *billingpb.CoinAccountInfo

	if info.Payment != nil {
		coinInfo := cache.GetEntry(cacheKey(cacheCoin, info.Payment.CoinInfoID))
		if coinInfo != nil {
			paymentCoinInfo = coinInfo.(*coininfopb.CoinInfo) //nolint
		}

		if paymentCoinInfo == nil {
			paymentCoinInfo, err := grpc2.GetCoinInfo(ctx, &coininfopb.GetCoinInfoRequest{
				ID: info.Payment.CoinInfoID,
			})
			if err != nil || paymentCoinInfo == nil {
				return nil, xerrors.Errorf("fail get payment coin info: %v", err)
			}

			cache.AddEntry(cacheKey(cacheCoin, info.Payment.CoinInfoID), paymentCoinInfo)
		}

		if !base {
			account, err := grpc2.GetBillingAccount(ctx, &billingpb.GetCoinAccountRequest{
				ID: info.Payment.AccountID,
			})
			if err != nil {
				return nil, xerrors.Errorf("fail get payment address: %v", err)
			}
			accountInfo = account
		}
	}

	var discountCoupon *inspirepb.CouponAllocatedDetail

	if !base && info.Order.DiscountCouponID != invalidUUID {
		discount := cache.GetEntry(cacheKey(cacheDiscount, info.Order.DiscountCouponID))
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
				return nil, xerrors.Errorf("fail check coupon usage: %v", err)
			}

			if order != nil && order.ID != info.Order.ID {
				return nil, xerrors.Errorf("fail check coupon usage")
			}

			discount, err := grpc2.GetCouponAllocated(ctx, &inspirepb.GetCouponAllocatedDetailRequest{
				ID: info.Order.DiscountCouponID,
			})
			if err != nil {
				return nil, xerrors.Errorf("fail get discount coupon allocated detail: %v", err)
			}

			if discount != nil {
				if discount.Allocated.AppID != info.Order.AppID || discount.Allocated.UserID != info.Order.UserID {
					return nil, xerrors.Errorf("invalid coupon")
				}
				discountCoupon = discount
				cache.AddEntry(cacheKey(cacheDiscount, info.Order.DiscountCouponID), discountCoupon)
			}
		}
	}

	var userSpecialReduction *inspirepb.UserSpecialReduction

	if !base && info.Order.UserSpecialReductionID != invalidUUID {
		special := cache.GetEntry(cacheKey(cacheSpecialOffer, info.Order.UserSpecialReductionID))
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
				return nil, xerrors.Errorf("fail check coupon usage: %v", err)
			}

			if order != nil && order.ID != info.Order.ID {
				return nil, xerrors.Errorf("fail check coupon usage")
			}

			userSpecial, err := grpc2.GetUserSpecialReduction(ctx, &inspirepb.GetUserSpecialReductionRequest{
				ID: info.Order.UserSpecialReductionID,
			})
			if err != nil {
				return nil, xerrors.Errorf("fail get user special reduction: %v", err)
			}

			if userSpecial != nil {
				if userSpecial.AppID != info.Order.AppID || userSpecial.UserID != info.Order.UserID {
					return nil, xerrors.Errorf("invalid coupon")
				}
				userSpecialReduction = userSpecial
				cache.AddEntry(cacheKey(cacheDiscount, info.Order.UserSpecialReductionID), userSpecialReduction)
			}
		}
	}

	var goodInfo *npool.Good

	cacheGoodInfo := cache.GetEntry(cacheKey(cacheGood, info.Order.GetGoodID()))
	if cacheGoodInfo != nil {
		goodInfo = cacheGoodInfo.(*npool.Good) //nolint
	}

	if goodInfo == nil {
		resp, err := gooddetail.Get(ctx, &npool.GetGoodRequest{
			ID: info.Order.GetGoodID(),
		})
		if err != nil || resp.Info == nil {
			return nil, xerrors.Errorf("fail get good info: %v", err)
		}
		goodInfo = resp.Info
		cache.AddEntry(cacheKey(cacheGood, info.Order.GetGoodID()), goodInfo)
	}

	var appGood *goodspb.AppGoodInfo

	cacheAppGoodInfo := cache.GetEntry(cacheKey(cacheAppGood, fmt.Sprintf(info.Order.GetAppID(), info.Order.GetGoodID())))
	if cacheAppGoodInfo != nil {
		appGood = cacheAppGoodInfo.(*goodspb.AppGoodInfo) //nolint
	}

	if appGood == nil {
		appGood, err = grpc2.GetAppGoodByAppGood(ctx, &goodspb.GetAppGoodByAppGoodRequest{
			AppID:  info.Order.GetAppID(),
			GoodID: info.Order.GetGoodID(),
		})
		if err != nil {
			return nil, xerrors.Errorf("fail get app good: %v", err)
		}
		cache.AddEntry(cacheKey(cacheAppGood, fmt.Sprintf(info.Order.GetAppID(), info.Order.GetGoodID())), appGood)
	}

	var promotion *goodspb.AppGoodPromotion

	cachedAppGoodPromotion := cache.GetEntry(cacheKey(cacheAppGoodPromotion, info.Order.PromotionID))
	if cachedAppGoodPromotion != nil {
		promotion = cachedAppGoodPromotion.(*goodspb.AppGoodPromotion) //nolint
	}

	if promotion == nil {
		promotion, err = grpc2.GetAppGoodPromotion(ctx, &goodspb.GetAppGoodPromotionRequest{
			ID: info.Order.PromotionID,
		})
		if err != nil {
			return nil, xerrors.Errorf("fail get promotion: %v", err)
		}
		if promotion != nil {
			if info.Order.GetAppID() != promotion.AppID || info.Order.GetGoodID() != promotion.GoodID {
				return nil, xerrors.Errorf("invalid promotion")
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

func GetOrder(ctx context.Context, in *npool.GetOrderRequest) (*npool.GetOrderResponse, error) {
	order, err := grpc2.GetOrderDetail(ctx, &orderpb.GetOrderDetailRequest{
		ID: in.GetID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get order detail: %v", err)
	}

	detail, err := expandOrder(ctx, order, false)
	if err != nil {
		return nil, xerrors.Errorf("fail expand order detail: %v", err)
	}

	return &npool.GetOrderResponse{
		Info: detail,
	}, nil
}

func GetOrdersByAppUser(ctx context.Context, in *npool.GetOrdersByAppUserRequest) (*npool.GetOrdersByAppUserResponse, error) {
	orders, err := grpc2.GetOrdersDetailByAppUser(ctx, &orderpb.GetOrdersDetailByAppUserRequest{
		AppID:  in.GetAppID(),
		UserID: in.GetUserID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get orders detail by app user: %v", err)
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

func GetOrdersShortDetailByAppUser(ctx context.Context, in *npool.GetOrdersByAppUserRequest) (*npool.GetOrdersByAppUserResponse, error) {
	orders, err := grpc2.GetOrdersShortDetailByAppUser(ctx, &orderpb.GetOrdersShortDetailByAppUserRequest{
		AppID:  in.GetAppID(),
		UserID: in.GetUserID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get orders detail by app user: %v", err)
	}

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].Order.CreateAt > orders[j].Order.CreateAt
	})

	details := []*npool.Order{}

	for _, info := range orders {
		detail, err := expandOrder(ctx, info, true)
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

func GetOrdersByApp(ctx context.Context, in *npool.GetOrdersByAppRequest) (*npool.GetOrdersByAppResponse, error) {
	orders, err := grpc2.GetOrdersDetailByApp(ctx, &orderpb.GetOrdersDetailByAppRequest{
		AppID: in.GetAppID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get orders detail by app: %v", err)
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

	return &npool.GetOrdersByAppResponse{
		Infos: details,
	}, nil
}

func GetOrdersByGood(ctx context.Context, in *npool.GetOrdersByGoodRequest) (*npool.GetOrdersByGoodResponse, error) {
	orders, err := grpc2.GetOrdersDetailByGood(ctx, &orderpb.GetOrdersDetailByGoodRequest{
		GoodID: in.GetGoodID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get orders detail by good: %v", err)
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

	return &npool.GetOrdersByGoodResponse{
		Infos: details,
	}, nil
}

func SubmitOrder(ctx context.Context, in *npool.SubmitOrderRequest) (*npool.SubmitOrderResponse, error) { //nolint
	payments, err := grpc2.GetPaymentsByAppUserState(ctx, &orderpb.GetPaymentsByAppUserStateRequest{
		AppID:  in.GetAppID(),
		UserID: in.GetUserID(),
		State:  orderconst.PaymentStateWait,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get wait payments: %v", err)
	}
	if len(payments) > constant.MaxUnpaidOrder {
		return nil, xerrors.Errorf("too many unpaid orders")
	}

	appGood, err := grpc2.GetAppGoodByAppGood(ctx, &goodspb.GetAppGoodByAppGoodRequest{
		AppID:  in.GetAppID(),
		GoodID: in.GetGoodID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get app good: %v", err)
	}
	if appGood == nil {
		return nil, xerrors.Errorf("fail get app good")
	}
	if !appGood.Online {
		return nil, xerrors.Errorf("good offline by app")
	}
	if appGood.Price <= 0 {
		return nil, xerrors.Errorf("good price invalid")
	}

	goodInfo, err := gooddetail.Get(ctx, &npool.GetGoodRequest{
		ID: in.GetGoodID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get order good info: %v", err)
	}

	if appGood.PurchaseLimit > 0 && int32(in.GetUnits()) > appGood.PurchaseLimit {
		return nil, xerrors.Errorf("too many units in a single order")
	}

	// Validate app id: done by gateway
	// Validate user id: done by gateway
	// Validate coupon id: done in expandOrder
	// TODO: Validate fee ids

	start := (uint32(time.Now().Unix()) + secondsInDay) / secondsInDay * secondsInDay
	if start < goodInfo.Info.Good.Good.StartAt {
		start = goodInfo.Info.Good.Good.StartAt
	}

	end := start + uint32(goodInfo.Info.Good.Good.DurationDays)*secondsInDay
	promotionID := uuid.UUID{}.String()

	promotion, err := grpc2.GetAppGoodPromotionByAppGoodTimestamp(ctx, &goodspb.GetAppGoodPromotionByAppGoodTimestampRequest{
		AppID:     in.GetAppID(),
		GoodID:    in.GetGoodID(),
		Timestamp: uint32(time.Now().Unix()),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get promotion: %v", err)
	}
	if promotion != nil {
		promotionID = promotion.ID
	}

	// Generate order
	myOrder, err := grpc2.CreateOrder(ctx, &orderpb.CreateOrderRequest{
		Info: &orderpb.Order{
			GoodID:                 in.GetGoodID(),
			AppID:                  in.GetAppID(),
			UserID:                 in.GetUserID(),
			Units:                  in.GetUnits(),
			Start:                  start,
			End:                    end,
			CouponID:               in.GetCouponID(),
			DiscountCouponID:       in.GetDiscountCouponID(),
			UserSpecialReductionID: in.GetUserSpecialReductionID(),
			PromotionID:            promotionID,
		},
	})
	if err != nil {
		return nil, xerrors.Errorf("fail create order: %v", err)
	}

	orderDetail, err := GetOrder(ctx, &npool.GetOrderRequest{
		ID: myOrder.ID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get order detail:%v", err)
	}

	return &npool.SubmitOrderResponse{
		Info: orderDetail.Info,
	}, nil
}

func peekIdlePaymentAccount(ctx context.Context, order *npool.Order, paymentCoinInfo *coininfopb.CoinInfo) (*billingpb.CoinAccountInfo, error) {
	payments, err := grpc2.GetIdleGoodPaymentsByGoodPaymentCoin(ctx, &billingpb.GetIdleGoodPaymentsByGoodPaymentCoinRequest{
		GoodID:            order.Good.Good.Good.ID,
		PaymentCoinTypeID: paymentCoinInfo.ID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get idle good payments: %v", err)
	}

	var paymentAccount *billingpb.GoodPayment

	for _, info := range payments {
		if uint32(time.Now().Unix()) <= info.AvailableAt {
			continue
		}

		if !info.Idle {
			continue
		}

		err = accountlock.Lock(info.AccountID)
		if err != nil {
			continue
		}

		paymentAccount = info
		break
	}

	if paymentAccount == nil {
		return nil, xerrors.Errorf("cannot find suitable payment account")
	}

	account, err := grpc2.GetBillingAccount(ctx, &billingpb.GetCoinAccountRequest{
		ID: paymentAccount.AccountID,
	})
	if err != nil {
		xerr := accountlock.Unlock(paymentAccount.AccountID)
		if xerr != nil {
			logger.Sugar().Errorf("cannot unlock %v: %v", paymentAccount.AccountID, xerr)
		}
		return nil, xerrors.Errorf("fail get account: %v", err)
	}

	paymentAccount.Idle = false
	paymentAccount.OccupiedBy = "paying"

	_, err = grpc2.UpdateGoodPayment(ctx, &billingpb.UpdateGoodPaymentRequest{
		Info: paymentAccount,
	})
	if err != nil {
		xerr := accountlock.Unlock(paymentAccount.AccountID)
		if xerr != nil {
			logger.Sugar().Errorf("cannot unlock %v: %v", paymentAccount.AccountID, xerr)
		}
		return nil, xerrors.Errorf("fail update good payment: %v", err)
	}

	return account, nil
}

func createNewPaymentAccount(ctx context.Context, order *npool.Order, paymentCoinInfo *coininfopb.CoinInfo) (*billingpb.CoinAccountInfo, error) {
	successCreated := 0

	for i := 0; i < 5; i++ {
		address, err := grpc2.CreateCoinAddress(ctx, &sphinxproxypb.CreateWalletRequest{
			Name: paymentCoinInfo.Name,
		})
		if err != nil {
			return nil, xerrors.Errorf("fail create wallet address: %v", err)
		}
		if address == nil || address.Address == "" {
			return nil, xerrors.Errorf("fail create wallet address for %v", paymentCoinInfo.Name)
		}

		account, err := grpc2.CreateBillingAccount(ctx, &billingpb.CreateCoinAccountRequest{
			Info: &billingpb.CoinAccountInfo{
				CoinTypeID:             paymentCoinInfo.ID,
				Address:                address.Address,
				PlatformHoldPrivateKey: true,
			},
		})
		if err != nil {
			return nil, xerrors.Errorf("fail create billing account: %v", err)
		}

		_, err = grpc2.CreateGoodPayment(ctx, &billingpb.CreateGoodPaymentRequest{
			Info: &billingpb.GoodPayment{
				GoodID:            order.Good.Good.Good.ID,
				PaymentCoinTypeID: paymentCoinInfo.ID,
				AccountID:         account.ID,
			},
		})
		if err != nil {
			return nil, xerrors.Errorf("fail create good payment: %v", err)
		}

		successCreated++
	}

	if successCreated > 0 {
		return peekIdlePaymentAccount(ctx, order, paymentCoinInfo)
	}

	return nil, xerrors.Errorf("SHOULD NOT BE HERE")
}

func CreateOrderPayment(ctx context.Context, in *npool.CreateOrderPaymentRequest) (*npool.CreateOrderPaymentResponse, error) { //nolint
	myOrder, err := GetOrder(ctx, &npool.GetOrderRequest{
		ID: in.GetOrderID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get order: %v", err)
	}
	if myOrder.Info.Order.Payment != nil {
		return &npool.CreateOrderPaymentResponse{
			Info: myOrder.Info,
		}, nil
	}

	paymentDeadline := time.Unix(int64(myOrder.Info.PaymentDeadline), 0)
	if time.Now().After(paymentDeadline) {
		return nil, xerrors.Errorf("order expired")
	}

	paymentCoinInfo, err := grpc2.GetCoinInfo(ctx, &coininfopb.GetCoinInfoRequest{
		ID: in.GetPaymentCoinTypeID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("invalid coin info id: %v", err)
	}

	if paymentCoinInfo.PreSale {
		return nil, xerrors.Errorf("cannot use presale coin as payment coin")
	}
	if !paymentCoinInfo.ForPay {
		return nil, xerrors.Errorf("payment coin not for pay")
	}
	if paymentCoinInfo.ENV != myOrder.Info.Good.Main.ENV {
		return nil, xerrors.Errorf("payment coin env different from good coin env")
	}

	paymentLiveCoinCurrency, err := currency.USDPrice(ctx, paymentCoinInfo.Name)
	if err != nil {
		return nil, xerrors.Errorf("cannot get usd currency for payment coin: %v", err)
	}

	paymentCoinCurrency := paymentLiveCoinCurrency
	paymentLocalCoinCurrency := 0.0

	payCurrency, err := oraclecli.GetCurrencyOnly(ctx,
		cruder.NewFilterConds().
			WithCond(oracleconst.FieldAppID, cruder.EQ, structpb.NewStringValue(myOrder.Info.Order.Order.AppID)).
			WithCond(oracleconst.FieldCoinTypeID, cruder.EQ, structpb.NewStringValue(in.GetPaymentCoinTypeID())))
	if err != nil {
		logger.Sugar().Errorf("fail get pay currency info: %v", err)
	}
	if payCurrency != nil {
		paymentCoinCurrency = payCurrency.AppPriceVSUSDT
		paymentLocalCoinCurrency = payCurrency.PriceVSUSDT
	}

	goodPrice := myOrder.Info.Good.Good.Good.Price
	if myOrder.Info.AppGood != nil {
		goodPrice = myOrder.Info.AppGood.Price
	}
	if myOrder.Info.Promotion != nil {
		goodPrice = myOrder.Info.Promotion.Price
	}

	amountUSD := float64(myOrder.Info.Order.Order.Units) * goodPrice
	if myOrder.Info.DiscountCoupon != nil {
		discount := myOrder.Info.DiscountCoupon.Discount
		if discount.Start < uint32(time.Now().Unix()) && uint32(time.Now().Unix()) < discount.Start+uint32(discount.DurationDays)*secondsInDay {
			amountUSD *= float64(100 - myOrder.Info.DiscountCoupon.Discount.Discount)
			amountUSD /= float64(100)
		}
	}
	if myOrder.Info.UserSpecialReduction != nil {
		userSpecial := myOrder.Info.UserSpecialReduction
		if userSpecial.Start < uint32(time.Now().Unix()) && uint32(time.Now().Unix()) < userSpecial.Start+uint32(userSpecial.DurationDays)*secondsInDay {
			amountUSD -= myOrder.Info.UserSpecialReduction.Amount
		}
	}
	if myOrder.Info.FixAmountCoupon != nil {
		coupon := myOrder.Info.FixAmountCoupon.Coupon
		if coupon.Start < uint32(time.Now().Unix()) && uint32(time.Now().Unix()) < coupon.Start+uint32(coupon.DurationDays)*secondsInDay {
			amountUSD -= myOrder.Info.FixAmountCoupon.Coupon.Denomination
		}
	}

	if amountUSD < 0 {
		amountUSD = 0
	}

	amountTarget := math.Ceil(amountUSD*10000/paymentCoinCurrency) / 10000
	extraAmount, err := fee.ExtraAmount(ctx, in.GetPaymentCoinTypeID())
	if err != nil {
		return nil, xerrors.Errorf("fail get extra payment amount: %v", err)
	}

	stock, err := stockcli.GetStockOnly(ctx, cruder.NewFilterConds().
		WithCond(stockconst.StockFieldGoodID, cruder.EQ, structpb.NewStringValue(myOrder.Info.AppGood.GoodID)))
	if err != nil || stock == nil {
		return nil, xerrors.Errorf("fail get good stock: %v", err)
	}

	stock, err = stockcli.AddStockFields(ctx, stock.ID, cruder.NewFilterFields().
		WithField(stockconst.StockFieldLocked, structpb.NewNumberValue(float64(myOrder.Info.Order.Order.Units))))
	if err != nil {
		return nil, xerrors.Errorf("fail add locked stock: %v", err)
	}

	defer func() {
		if err != nil {
			logger.Sugar().Errorf("try revert locked stock: %v", err)
			_, err = stockcli.AddStockFields(ctx, stock.ID, cruder.NewFilterFields().
				WithField(stockconst.StockFieldLocked, structpb.NewNumberValue(float64(int32(myOrder.Info.Order.Order.Units)*-1))))
			if err != nil {
				logger.Sugar().Errorf("fail sub locked stock: %v", err)
			}
		}
	}()

	logger.Sugar().Infof("purchase %v goods with price %v amountUSD %v amount target %v currency %v",
		myOrder.Info.Order.Order.Units,
		goodPrice,
		amountUSD,
		amountTarget,
		paymentCoinCurrency)

	// TODO: Check if idle address is available with lock
	paymentAccount, err := peekIdlePaymentAccount(ctx, myOrder.Info, paymentCoinInfo)
	if err != nil {
		paymentAccount, err = createNewPaymentAccount(ctx, myOrder.Info, paymentCoinInfo)
	}
	if err != nil || paymentAccount == nil {
		return nil, xerrors.Errorf("cannot get valid payment account: %v", err)
	}

	balance, err := grpc2.GetBalance(ctx, &sphinxproxypb.GetBalanceRequest{
		Name:    paymentCoinInfo.Name,
		Address: paymentAccount.Address,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get wallet balance: %v", err)
	}
	balanceAmount := balance.Balance

	// Generate payment
	myPayment, err := grpc2.CreatePayment(ctx, &orderpb.CreatePaymentRequest{
		Info: &orderpb.Payment{
			AppID:                myOrder.Info.Order.Order.AppID,
			UserID:               myOrder.Info.Order.Order.UserID,
			GoodID:               myOrder.Info.Order.Order.GoodID,
			OrderID:              myOrder.Info.Order.Order.ID,
			AccountID:            paymentAccount.ID,
			StartAmount:          balanceAmount,
			Amount:               amountTarget,
			CoinUSDCurrency:      paymentCoinCurrency,
			LocalCoinUSDCurrency: paymentLocalCoinCurrency,
			LiveCoinUSDCurrency:  paymentLiveCoinCurrency,
			CoinInfoID:           in.GetPaymentCoinTypeID(),
		},
	})
	if err != nil {
		return nil, xerrors.Errorf("fail create payment: %v", err)
	}

	// Generate good paying
	_, err = grpc2.CreateGoodPaying(ctx, &orderpb.CreateGoodPayingRequest{
		Info: &orderpb.GoodPaying{
			OrderID:   myOrder.Info.Order.Order.ID,
			PaymentID: myPayment.ID,
		},
	})
	if err != nil {
		return nil, xerrors.Errorf("fail create good paying: %v", err)
	}

	// Generate gas payings
	for _, fee := range in.GetFees() {
		_, err := grpc2.CreateGasPaying(ctx, &orderpb.CreateGasPayingRequest{
			Info: &orderpb.GasPaying{
				OrderID:         myOrder.Info.Order.Order.ID,
				PaymentID:       myPayment.ID,
				DurationMinutes: fee.DurationDays * 24 * 60,
				FeeTypeID:       fee.ID,
			},
		})
		if err != nil {
			return nil, xerrors.Errorf("fail create fee: %v", err)
		}
	}

	orderDetail, err := GetOrder(ctx, &npool.GetOrderRequest{
		ID: myOrder.Info.Order.Order.ID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get order detail: %v", err)
	}

	// Watch payment address and change payment state
	orderDetail.Info.Order.Payment.Amount += extraAmount

	return &npool.CreateOrderPaymentResponse{
		Info: orderDetail.Info,
	}, nil
}
