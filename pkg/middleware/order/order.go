package order

import (
	"context"
	"fmt"
	"time"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	gooddetail "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/good" //nolint
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	billingpb "github.com/NpoolPlatform/message/npool/cloud-hashing-billing"
	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"
	orderpb "github.com/NpoolPlatform/message/npool/cloud-hashing-order"
	coininfopb "github.com/NpoolPlatform/message/npool/coininfo"
	sphinxproxypb "github.com/NpoolPlatform/message/npool/sphinxproxy"

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

	// Validate app id
	// Validate user id
	// Validate coupon id
	// Validate fee ids

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

func CreateOrderPayment(ctx context.Context, in *npool.CreateOrderPaymentRequest) (*npool.CreateOrderPaymentResponse, error) {
	myOrder, err := GetOrder(ctx, &npool.GetOrderRequest{
		ID: in.GetOrderID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get order: %v", err)
	}

	// Caculate amount
	amount := float64(myOrder.Info.Order.Units) * myOrder.Info.Good.Good.Price

	// TODO: All should validate duration days
	// User discount info
	if myOrder.Info.DiscountCoupon != nil {
		amount *= float64(100 - myOrder.Info.DiscountCoupon.Discount.Discount)
	}
	amount /= float64(100)
	// Extra reduction
	if myOrder.Info.UserSpecialReduction != nil {
		amount -= myOrder.Info.UserSpecialReduction.Amount
	}
	// Coupon amount
	if myOrder.Info.FixAmountCoupon != nil {
		amount -= myOrder.Info.FixAmountCoupon.Coupon.Denomination
	}

	// Validate payment coin info id
	coinInfo, err := grpc2.GetCoinInfo(ctx, &coininfopb.GetCoinInfoRequest{
		ID: in.GetPaymentCoinTypeID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("invalid coin info id: %v", err)
	}

	// Check if idle address is available
	idle := false
	var coinAddress string

	if coinInfo.Info.PreSale {
		coinAddress = fmt.Sprintf("placeholder-%v", uuid.New())
	} else {
		// Generate transaction address
		address, err := grpc2.CreateCoinAddress(ctx, &sphinxproxypb.CreateWalletRequest{
			Name: coinInfo.Info.Name,
		})
		if err != nil {
			return nil, xerrors.Errorf("fail create wallet address: %v", err)
		}
		coinAddress = address.Info.Address
	}

	// Create billing account
	account, err := grpc2.CreateBillingAccount(ctx, &billingpb.CreateCoinAccountRequest{
		Info: &billingpb.CoinAccountInfo{
			CoinTypeID:  in.GetPaymentCoinTypeID(),
			Address:     coinAddress,
			GeneratedBy: "platform",
			UsedFor:     "paying",
			AppID:       myOrder.Info.Order.AppID,
			UserID:      myOrder.Info.Order.UserID,
		},
	})
	if err != nil {
		return nil, xerrors.Errorf("fail create billing account: %v", err)
	}

	balanceAmount := float64(0)

	if !idle && !coinInfo.Info.PreSale {
		balance, err := grpc2.GetBalance(ctx, &sphinxproxypb.GetBalanceRequest{
			Name:    coinInfo.Info.Name,
			Address: coinAddress,
		})
		if err != nil {
			return nil, xerrors.Errorf("fail get wallet balance: %v", err)
		}
		balanceAmount = balance.Info.Balance
	}

	// Generate payment
	myPayment, err := grpc2.CreatePayment(ctx, &orderpb.CreatePaymentRequest{
		Info: &orderpb.Payment{
			OrderID:     myOrder.Info.Order.ID,
			AccountID:   account.Info.ID,
			StartAmount: balanceAmount,
			Amount:      amount,
			CoinInfoID:  in.GetPaymentCoinTypeID(),
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
