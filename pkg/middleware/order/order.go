package order

import (
	"context"
	"math"
	"time"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	redis2 "github.com/NpoolPlatform/go-service-framework/pkg/redis"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	gooddetail "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/good" //nolint
	currency "github.com/NpoolPlatform/cloud-hashing-staker/pkg/middleware/currency"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	billingpb "github.com/NpoolPlatform/message/npool/cloud-hashing-billing"
	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"
	orderpb "github.com/NpoolPlatform/message/npool/cloud-hashing-order"
	coininfopb "github.com/NpoolPlatform/message/npool/coininfo"
	sphinxproxypb "github.com/NpoolPlatform/message/npool/sphinxproxy"

	orderconst "github.com/NpoolPlatform/cloud-hashing-order/pkg/const"
	paymentwatcher "github.com/NpoolPlatform/cloud-hashing-staker/pkg/middleware/payment-watcher"

	"github.com/google/uuid"

	"golang.org/x/xerrors"
)

func constructOrder(
	info *orderpb.OrderDetail,
	goodInfo *npool.Good,
	coupon *inspirepb.CouponAllocatedDetail,
	paymentCoinInfo *coininfopb.CoinInfo,
	account *billingpb.CoinAccountInfo,
	discountCoupon *inspirepb.CouponAllocatedDetail,
	userSpecial *inspirepb.UserSpecialReduction) *npool.Order {
	return &npool.Order{
		Order:                info,
		PayToAccount:         account,
		Good:                 goodInfo,
		PayWithCoin:          paymentCoinInfo,
		FixAmountCoupon:      coupon,
		DiscountCoupon:       discountCoupon,
		UserSpecialReduction: userSpecial,
		PaymentDeadline:      info.CreateAt + orderconst.TimeoutSeconds,
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

	if !base {
		couponAllocated, err := grpc2.GetCouponAllocated(ctx, &inspirepb.GetCouponAllocatedDetailRequest{
			ID: info.CouponID,
		})
		if err != nil && info.CouponID != invalidUUID {
			return nil, xerrors.Errorf("fail get coupon allocated detail: %v", err)
		}

		if couponAllocated != nil {
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

	if !base {
		discount, err := grpc2.GetCouponAllocated(ctx, &inspirepb.GetCouponAllocatedDetailRequest{
			ID: info.DiscountCouponID,
		})
		if err != nil && info.DiscountCouponID != invalidUUID {
			return nil, xerrors.Errorf("fail get discount coupon allocated detail: %v", err)
		}

		if discount != nil {
			discountCoupon = discount.Info
		}
	}

	var userSpecialReduction *inspirepb.UserSpecialReduction

	if !base {
		userSpecial, err := grpc2.GetUserSpecialReduction(ctx, &inspirepb.GetUserSpecialReductionRequest{
			ID: info.UserSpecialReductionID,
		})
		if err != nil && info.UserSpecialReductionID != invalidUUID {
			return nil, xerrors.Errorf("fail get user special reduction: %v", err)
		}

		if userSpecial != nil {
			userSpecialReduction = userSpecial.Info
		}
	}

	var goodInfo *npool.Good
	if detail, ok := goodsDetail[info.GetGoodID()]; ok {
		goodInfo = detail
	} else {
		resp, err := gooddetail.Get(ctx, &npool.GetGoodRequest{
			ID: info.GetGoodID(),
		})
		if err != nil {
			return nil, xerrors.Errorf("fail get good info: %v", err)
		}
		goodsDetail[info.GetGoodID()] = resp.Info
		goodInfo = resp.Info
	}

	return constructOrder(
		info,
		goodInfo,
		coupon,
		paymentCoinInfo,
		accountInfo,
		discountCoupon,
		userSpecialReduction), nil
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

	detail, err := expandOrder(ctx, orderDetail.Detail, false, goods, coins)
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

	for _, info := range ordersDetail.Details {
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

	for _, info := range ordersDetail.Details {
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

	for _, info := range ordersDetail.Details {
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

	for _, info := range ordersDetail.Details {
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

	// Validate app id: done by gateway
	// Validate user id: done by gateway
	// TODO: Validate coupon id
	// TODO: Validate fee ids

	start := (uint32(time.Now().Unix()) + 24*60*60) / 24 / 60 / 60 * 24 * 60 * 60
	end := start + uint32(goodInfo.Info.Good.DurationDays)*24*60*60

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
		GoodID:            order.Good.Good.ID,
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
				GoodID:            order.Good.Good.ID,
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

	paymentCoinCurrency, err := currency.USDPrice(ctx, paymentCoinInfo.Info.Name)
	if err != nil {
		return nil, xerrors.Errorf("cannot get usd currency for payment coin: %v", err)
	}

	amountUSD := float64(myOrder.Info.Order.Units) * myOrder.Info.Good.Good.Price
	if myOrder.Info.DiscountCoupon != nil {
		amountUSD *= float64(100 - myOrder.Info.DiscountCoupon.Discount.Discount)
		amountUSD /= float64(100)
	}
	if myOrder.Info.UserSpecialReduction != nil {
		amountUSD -= myOrder.Info.UserSpecialReduction.Amount
	}
	if myOrder.Info.FixAmountCoupon != nil {
		amountUSD -= myOrder.Info.FixAmountCoupon.Coupon.Denomination
	}

	amountTarget := math.Ceil(amountUSD*10000/paymentCoinCurrency) / 10000

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
			OrderID:         myOrder.Info.Order.ID,
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
			OrderID:   myOrder.Info.Order.ID,
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
				OrderID:         myOrder.Info.Order.ID,
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
		ID: myOrder.Info.Order.ID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get order detail: %v", err)
	}

	// Watch payment address and change payment state

	return &npool.CreateOrderPaymentResponse{
		Info: orderDetail.Info,
	}, nil
}
