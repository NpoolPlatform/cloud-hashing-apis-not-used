package order

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	redis2 "github.com/NpoolPlatform/go-service-framework/pkg/redis"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	gooddetail "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/good" //nolint
	currency "github.com/NpoolPlatform/cloud-hashing-staker/pkg/middleware/currency"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	billingpb "github.com/NpoolPlatform/message/npool/cloud-hashing-billing"
	goodspb "github.com/NpoolPlatform/message/npool/cloud-hashing-goods"
	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"
	orderpb "github.com/NpoolPlatform/message/npool/cloud-hashing-order"
	coininfopb "github.com/NpoolPlatform/message/npool/coininfo"
	sphinxproxypb "github.com/NpoolPlatform/message/npool/sphinxproxy"

	orderconst "github.com/NpoolPlatform/cloud-hashing-order/pkg/const"
	paymentwatcher "github.com/NpoolPlatform/cloud-hashing-staker/pkg/middleware/payment-watcher"

	"github.com/google/uuid"

	"golang.org/x/xerrors"
)

const (
	secondsInDay = 24 * 60 * 60
)

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

func expandOrder( //nolint
	ctx context.Context,
	info *orderpb.OrderDetail,
	base bool,
	goodsDetail map[string]*npool.Good,
	coinInfos map[string]*coininfopb.CoinInfo) (*npool.Order, error) { //nolint
	var coupon *inspirepb.CouponAllocatedDetail
	invalidUUID := uuid.UUID{}.String()

	if !base && info.Order.CouponID != invalidUUID {
		resp, err := grpc2.GetOrderByAppUserCouponTypeID(ctx, &orderpb.GetOrderByAppUserCouponTypeIDRequest{
			AppID:      info.Order.AppID,
			UserID:     info.Order.UserID,
			CouponType: orderconst.FixAmountCoupon,
			CouponID:   info.Order.CouponID,
		})
		if err != nil {
			return nil, xerrors.Errorf("fail check coupon usage: %v", err)
		}
		if resp.Info != nil {
			return nil, xerrors.Errorf("fail check coupon usage")
		}

		couponAllocated, err := grpc2.GetCouponAllocated(ctx, &inspirepb.GetCouponAllocatedDetailRequest{
			ID: info.Order.CouponID,
		})
		if err != nil && info.Order.CouponID != invalidUUID {
			return nil, xerrors.Errorf("fail get coupon allocated detail: %v", err)
		}

		if couponAllocated.Info != nil {
			if couponAllocated.Info.Allocated.AppID != info.Order.AppID || couponAllocated.Info.Allocated.UserID != info.Order.UserID {
				return nil, xerrors.Errorf("invalid coupon")
			}
			coupon = couponAllocated.Info
		}
	}

	var paymentCoinInfo *coininfopb.CoinInfo
	var accountInfo *billingpb.CoinAccountInfo

	if info.Payment != nil {
		if coinInfo, ok := coinInfos[info.Payment.CoinInfoID]; ok {
			paymentCoinInfo = coinInfo
		} else {
			coinInfo, err := grpc2.GetCoinInfo(ctx, &coininfopb.GetCoinInfoRequest{
				ID: info.Payment.CoinInfoID,
			})
			if err != nil {
				return nil, xerrors.Errorf("fail get payment coin info: %v", err)
			}
			paymentCoinInfo = coinInfo.Info
			coinInfos[info.Payment.CoinInfoID] = paymentCoinInfo
		}

		if !base {
			account, err := grpc2.GetBillingAccount(ctx, &billingpb.GetCoinAccountRequest{
				ID: info.Payment.AccountID,
			})
			if err != nil {
				return nil, xerrors.Errorf("fail get payment address: %v", err)
			}
			accountInfo = account.Info
		}
	}

	var discountCoupon *inspirepb.CouponAllocatedDetail

	if !base && info.Order.DiscountCouponID != invalidUUID {
		resp, err := grpc2.GetOrderByAppUserCouponTypeID(ctx, &orderpb.GetOrderByAppUserCouponTypeIDRequest{
			AppID:      info.Order.AppID,
			UserID:     info.Order.UserID,
			CouponType: orderconst.DiscountCoupon,
			CouponID:   info.Order.DiscountCouponID,
		})
		if err != nil {
			return nil, xerrors.Errorf("fail check coupon usage: %v", err)
		}
		if resp.Info != nil {
			return nil, xerrors.Errorf("fail check coupon usage")
		}

		discount, err := grpc2.GetCouponAllocated(ctx, &inspirepb.GetCouponAllocatedDetailRequest{
			ID: info.Order.DiscountCouponID,
		})
		if err != nil {
			return nil, xerrors.Errorf("fail get discount coupon allocated detail: %v", err)
		}

		if discount.Info != nil {
			if discount.Info.Allocated.AppID != info.Order.AppID || discount.Info.Allocated.UserID != info.Order.UserID {
				return nil, xerrors.Errorf("invalid coupon")
			}
			discountCoupon = discount.Info
		}
	}

	var userSpecialReduction *inspirepb.UserSpecialReduction

	if !base && info.Order.UserSpecialReductionID != invalidUUID {
		resp, err := grpc2.GetOrderByAppUserCouponTypeID(ctx, &orderpb.GetOrderByAppUserCouponTypeIDRequest{
			AppID:      info.Order.AppID,
			UserID:     info.Order.UserID,
			CouponType: orderconst.UserSpecialReductionCoupon,
			CouponID:   info.Order.UserSpecialReductionID,
		})
		if err != nil {
			return nil, xerrors.Errorf("fail check coupon usage: %v", err)
		}
		if resp.Info != nil {
			return nil, xerrors.Errorf("fail check coupon usage")
		}

		userSpecial, err := grpc2.GetUserSpecialReduction(ctx, &inspirepb.GetUserSpecialReductionRequest{
			ID: info.Order.UserSpecialReductionID,
		})
		if err != nil {
			return nil, xerrors.Errorf("fail get user special reduction: %v", err)
		}

		if userSpecial.Info != nil {
			if userSpecial.Info.AppID != info.Order.AppID || userSpecial.Info.UserID != info.Order.UserID {
				return nil, xerrors.Errorf("invalid coupon")
			}
			userSpecialReduction = userSpecial.Info
		}
	}

	var goodInfo *npool.Good
	if detail, ok := goodsDetail[info.Order.GetGoodID()]; ok {
		goodInfo = detail
	} else {
		resp, err := gooddetail.Get(ctx, &npool.GetGoodRequest{
			ID: info.Order.GetGoodID(),
		})
		if err != nil {
			return nil, xerrors.Errorf("fail get good info: %v", err)
		}
		goodsDetail[info.Order.GetGoodID()] = resp.Info
		goodInfo = resp.Info
	}

	var appGood *goodspb.AppGoodInfo
	appGoodResp, err := grpc2.GetAppGoodByAppGood(ctx, &goodspb.GetAppGoodByAppGoodRequest{
		AppID:  info.Order.GetAppID(),
		GoodID: info.Order.GetGoodID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get app good: %v", err)
	}
	if appGoodResp.Info != nil {
		appGood = appGoodResp.Info
	}

	var promotion *goodspb.AppGoodPromotion
	promotionResp, err := grpc2.GetAppGoodPromotion(ctx, &goodspb.GetAppGoodPromotionRequest{
		ID: info.Order.PromotionID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get promotion: %v", err)
	}
	if promotionResp.Info != nil {
		if info.Order.GetAppID() != promotionResp.Info.AppID || info.Order.GetGoodID() != promotionResp.Info.GoodID {
			return nil, xerrors.Errorf("invalid promotion")
		}
		promotion = promotionResp.Info
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
	orderDetail, err := grpc2.GetOrderDetail(ctx, &orderpb.GetOrderDetailRequest{
		ID: in.GetID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get order detail: %v", err)
	}

	goods := map[string]*npool.Good{}
	coins := map[string]*coininfopb.CoinInfo{}

	detail, err := expandOrder(ctx, orderDetail.Info, false, goods, coins)
	if err != nil {
		return nil, xerrors.Errorf("fail expand order detail: %v", err)
	}

	return &npool.GetOrderResponse{
		Info: detail,
	}, nil
}

func GetOrdersByAppUser(ctx context.Context, in *npool.GetOrdersByAppUserRequest) (*npool.GetOrdersByAppUserResponse, error) {
	ordersDetail, err := grpc2.GetOrdersDetailByAppUser(ctx, &orderpb.GetOrdersDetailByAppUserRequest{
		AppID:  in.GetAppID(),
		UserID: in.GetUserID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get orders detail by app user: %v", err)
	}

	details := []*npool.Order{}
	goods := map[string]*npool.Good{}
	coins := map[string]*coininfopb.CoinInfo{}

	for _, info := range ordersDetail.Infos {
		detail, err := expandOrder(ctx, info, false, goods, coins)
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

func GetOrdersShortDetailByAppUser( //nolint
	ctx context.Context,
	in *npool.GetOrdersByAppUserRequest,
	goods map[string]*npool.Good,
	coins map[string]*coininfopb.CoinInfo) (*npool.GetOrdersByAppUserResponse,
	map[string]*npool.Good,
	map[string]*coininfopb.CoinInfo,
	error) {
	ordersDetail, err := grpc2.GetOrdersShortDetailByAppUser(ctx, &orderpb.GetOrdersShortDetailByAppUserRequest{
		AppID:  in.GetAppID(),
		UserID: in.GetUserID(),
	})
	if err != nil {
		return nil, nil, nil, xerrors.Errorf("fail get orders detail by app user: %v", err)
	}

	details := []*npool.Order{}

	for _, info := range ordersDetail.Infos {
		detail, err := expandOrder(ctx, info, true, goods, coins)
		if err != nil {
			logger.Sugar().Warnf("cannot expand order detail: %v", err)
			continue
		}

		details = append(details, detail)
	}

	return &npool.GetOrdersByAppUserResponse{
		Infos: details,
	}, goods, coins, nil
}

func GetOrdersByApp(ctx context.Context, in *npool.GetOrdersByAppRequest) (*npool.GetOrdersByAppResponse, error) {
	ordersDetail, err := grpc2.GetOrdersDetailByApp(ctx, &orderpb.GetOrdersDetailByAppRequest{
		AppID: in.GetAppID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get orders detail by app: %v", err)
	}

	details := []*npool.Order{}
	goods := map[string]*npool.Good{}
	coins := map[string]*coininfopb.CoinInfo{}

	for _, info := range ordersDetail.Infos {
		detail, err := expandOrder(ctx, info, false, goods, coins)
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
	ordersDetail, err := grpc2.GetOrdersDetailByGood(ctx, &orderpb.GetOrdersDetailByGoodRequest{
		GoodID: in.GetGoodID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get orders detail by good: %v", err)
	}

	details := []*npool.Order{}
	goods := map[string]*npool.Good{}
	coins := map[string]*coininfopb.CoinInfo{}

	for _, info := range ordersDetail.Infos {
		detail, err := expandOrder(ctx, info, false, goods, coins)
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

func SubmitOrder(ctx context.Context, in *npool.SubmitOrderRequest) (*npool.SubmitOrderResponse, error) {
	goodInfo, err := gooddetail.Get(ctx, &npool.GetGoodRequest{
		ID: in.GetGoodID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get order good info: %v", err)
	}

	if in.GetUnits() > uint32(goodInfo.Info.Good.Good.Total) {
		return nil, xerrors.Errorf("invalid units")
	}

	// Validate app id: done by gateway
	// Validate user id: done by gateway
	// Validate coupon id: done in expandOrder
	// TODO: Validate fee ids

	lockKey := fmt.Sprintf("submit-order:%v", in.GetGoodID())
	err = redis2.TryLock(lockKey, 0)
	if err != nil {
		return nil, xerrors.Errorf("fail lock good: %v", err)
	}
	defer func() {
		err := redis2.Unlock(lockKey)
		if err != nil {
			logger.Sugar().Errorf("fail unlock good: %v", err)
		}
	}()

	sold, err := grpc2.GetSoldByGood(ctx, &orderpb.GetSoldByGoodRequest{
		GoodID: in.GetGoodID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get good sold: %v", err)
	}

	if sold.Sold >= uint32(goodInfo.Info.Good.Good.Total) {
		return nil, xerrors.Errorf("good sold out")
	}

	start := (uint32(time.Now().Unix()) + secondsInDay) / secondsInDay * secondsInDay
	if start < goodInfo.Info.Good.Good.StartAt {
		start = goodInfo.Info.Good.Good.StartAt
	}

	end := start + uint32(goodInfo.Info.Good.Good.DurationDays)*secondsInDay
	promotionID := uuid.UUID{}.String()

	resp, err := grpc2.GetAppGoodPromotionByAppGoodTimestamp(ctx, &goodspb.GetAppGoodPromotionByAppGoodTimestampRequest{
		AppID:     in.GetAppID(),
		GoodID:    in.GetGoodID(),
		Timestamp: uint32(time.Now().Unix()),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get promotion: %v", err)
	}
	if resp.Info != nil {
		promotionID = resp.Info.ID
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
		ID: myOrder.Info.ID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get order detail:%v", err)
	}

	return &npool.SubmitOrderResponse{
		Info: orderDetail.Info,
	}, nil
}

func peekIdlePaymentAccount(ctx context.Context, order *npool.Order, paymentCoinInfo *coininfopb.CoinInfo) (*billingpb.CoinAccountInfo, error) {
	resp, err := grpc2.GetIdleGoodPaymentsByGoodPaymentCoin(ctx, &billingpb.GetIdleGoodPaymentsByGoodPaymentCoinRequest{
		GoodID:            order.Good.Good.Good.ID,
		PaymentCoinTypeID: paymentCoinInfo.ID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get idle good payments: %v", err)
	}

	var paymentAccount *billingpb.GoodPayment
	var lockKey string

	for _, info := range resp.Infos {
		lockKey = paymentwatcher.AccountLockKey(info.ID)
		err = redis2.TryLock(lockKey, orderconst.TimeoutSeconds*2)
		if err != nil {
			continue
		}

		paymentAccount = info
		break
	}

	if paymentAccount == nil {
		return nil, xerrors.Errorf("cannot find suitable payment account")
	}

	paymentAccount.Idle = false
	paymentAccount.OccupiedBy = "paying"

	_, err = grpc2.UpdateGoodPayment(ctx, &billingpb.UpdateGoodPaymentRequest{
		Info: paymentAccount,
	})
	if err != nil {
		xerr := redis2.Unlock(lockKey)
		if xerr != nil {
			logger.Sugar().Errorf("cannot unlock %v: %v", lockKey, xerr)
		}
		return nil, xerrors.Errorf("fail update good payment: %v", err)
	}

	resp1, err := grpc2.GetBillingAccount(ctx, &billingpb.GetCoinAccountRequest{
		ID: paymentAccount.AccountID,
	})
	if err != nil {
		xerr := redis2.Unlock(lockKey)
		if xerr != nil {
			logger.Sugar().Errorf("cannot unlock %v: %v", lockKey, xerr)
		}
		return nil, xerrors.Errorf("fail get account: %v", err)
	}

	return resp1.Info, nil
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
		if address.Info == nil || address.Info.Address == "" {
			return nil, xerrors.Errorf("fail create wallet address for %v", paymentCoinInfo.Name)
		}

		account, err := grpc2.CreateBillingAccount(ctx, &billingpb.CreateCoinAccountRequest{
			Info: &billingpb.CoinAccountInfo{
				CoinTypeID:             paymentCoinInfo.ID,
				Address:                address.Info.Address,
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
				AccountID:         account.Info.ID,
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

	if paymentCoinInfo.Info.PreSale {
		return nil, xerrors.Errorf("cannot use presale coin as payment coin")
	}
	if !paymentCoinInfo.Info.ForPay {
		return nil, xerrors.Errorf("payment coin not for pay")
	}
	if paymentCoinInfo.Info.ENV != myOrder.Info.Good.Main.ENV {
		return nil, xerrors.Errorf("payment coin env different from good coin env")
	}

	paymentCoinCurrency, err := currency.USDPrice(ctx, paymentCoinInfo.Info.Name)
	if err != nil {
		return nil, xerrors.Errorf("cannot get usd currency for payment coin: %v", err)
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
	logger.Sugar().Infof("purchase %v goods with price %v amountUSD %v amount target %v currency %v",
		myOrder.Info.Order.Order.Units,
		goodPrice,
		amountUSD,
		amountTarget,
		paymentCoinCurrency)

	// TODO: Check if idle address is available with lock
	paymentAccount, err := peekIdlePaymentAccount(ctx, myOrder.Info, paymentCoinInfo.Info)
	if err != nil {
		paymentAccount, err = createNewPaymentAccount(ctx, myOrder.Info, paymentCoinInfo.Info)
	}
	if err != nil {
		return nil, xerrors.Errorf("cannot get valid payment account: %v", err)
	}

	balance, err := grpc2.GetBalance(ctx, &sphinxproxypb.GetBalanceRequest{
		Name:    paymentCoinInfo.Info.Name,
		Address: paymentAccount.Address,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get wallet balance: %v", err)
	}
	balanceAmount := balance.Info.Balance

	// Generate payment
	myPayment, err := grpc2.CreatePayment(ctx, &orderpb.CreatePaymentRequest{
		Info: &orderpb.Payment{
			AppID:           myOrder.Info.Order.Order.AppID,
			UserID:          myOrder.Info.Order.Order.UserID,
			GoodID:          myOrder.Info.Order.Order.GoodID,
			OrderID:         myOrder.Info.Order.Order.ID,
			AccountID:       paymentAccount.ID,
			StartAmount:     balanceAmount,
			Amount:          amountTarget,
			CoinUSDCurrency: paymentCoinCurrency,
			CoinInfoID:      in.GetPaymentCoinTypeID(),
		},
	})
	if err != nil {
		return nil, xerrors.Errorf("fail create payment: %v", err)
	}

	// Generate good paying
	_, err = grpc2.CreateGoodPaying(ctx, &orderpb.CreateGoodPayingRequest{
		Info: &orderpb.GoodPaying{
			OrderID:   myOrder.Info.Order.Order.ID,
			PaymentID: myPayment.Info.ID,
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
				PaymentID:       myPayment.Info.ID,
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

	return &npool.CreateOrderPaymentResponse{
		Info: orderDetail.Info,
	}, nil
}
