package order

import (
	"context"
	"fmt"
	"time"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	"github.com/NpoolPlatform/cloud-hashing-apis/message/npool"
	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	gooddetail "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/good-detail" //nolint

	billingpb "github.com/NpoolPlatform/cloud-hashing-billing/message/npool"
	inspirepb "github.com/NpoolPlatform/cloud-hashing-inspire/message/npool"
	orderpb "github.com/NpoolPlatform/cloud-hashing-order/message/npool"
	coininfopb "github.com/NpoolPlatform/message/npool/coininfo"
	sphinxproxypb "github.com/NpoolPlatform/message/npool/sphinxproxy"

	"github.com/google/uuid"

	"golang.org/x/xerrors"
)

func constructOrderDetail(
	info *orderpb.OrderDetail,
	goodInfo *npool.GoodDetail,
	coupon *inspirepb.CouponAllocatedDetail,
	paymentCoinInfo *coininfopb.CoinInfo,
	account *billingpb.CoinAccountInfo,
	discountCoupon *inspirepb.CouponAllocatedDetail,
	userSpecial *inspirepb.UserSpecialReduction) *npool.OrderDetail {
	gasPayings := []*npool.GasPaying{}
	for _, paying := range info.GasPayings {
		gasPayings = append(gasPayings, &npool.GasPaying{
			ID:              paying.ID,
			OrderID:         paying.OrderID,
			PaymentID:       paying.PaymentID,
			DurationMinutes: paying.DurationMinutes,
		})
	}

	compensates := []*npool.Compensate{}
	for _, cm := range info.Compensates {
		compensates = append(compensates, &npool.Compensate{
			ID:      cm.ID,
			OrderID: cm.OrderID,
			Start:   cm.Start,
			End:     cm.End,
			Message: cm.Message,
		})
	}

	outOfGases := []*npool.OutOfGas{}
	for _, ofg := range info.OutOfGases {
		outOfGases = append(outOfGases, &npool.OutOfGas{
			ID:      ofg.ID,
			OrderID: ofg.OrderID,
			Start:   ofg.Start,
			End:     ofg.End,
		})
	}

	var myCoupon *npool.Coupon
	if coupon != nil {
		myCoupon = &npool.Coupon{
			ID:     coupon.ID,
			UserID: coupon.UserID,
			AppID:  coupon.AppID,
			Pool: &npool.CouponPool{
				ID:           coupon.Coupon.ID,
				AppID:        coupon.Coupon.AppID,
				Denomination: coupon.Coupon.Denomination,
				Start:        coupon.Coupon.Start,
				DurationDays: coupon.Coupon.DurationDays,
				Message:      coupon.Coupon.Message,
				Name:         coupon.Coupon.Name,
			},
		}
	}

	var myDiscount *npool.Discount
	discountAmount := uint32(0)

	if discountCoupon != nil {
		myDiscount = &npool.Discount{
			ID:     discountCoupon.ID,
			AppID:  discountCoupon.AppID,
			UserID: discountCoupon.UserID,
			Pool: &npool.DiscountPool{
				ID:           discountCoupon.Discount.ID,
				AppID:        discountCoupon.Discount.AppID,
				Discount:     discountCoupon.Discount.Discount,
				Start:        discountCoupon.Discount.Start,
				DurationDays: discountCoupon.Discount.DurationDays,
				Message:      discountCoupon.Discount.Message,
				Name:         discountCoupon.Discount.Name,
			},
		}
		discountAmount = discountCoupon.Discount.Discount
	}

	var specialReduction *npool.UserSpecialReduction
	reductionAmount := float64(0)

	if userSpecial != nil {
		specialReduction = &npool.UserSpecialReduction{
			ID:           userSpecial.ID,
			AppID:        userSpecial.AppID,
			UserID:       userSpecial.UserID,
			Amount:       userSpecial.Amount,
			Start:        userSpecial.Start,
			DurationDays: userSpecial.DurationDays,
		}
		reductionAmount = userSpecial.Amount
	}

	var coinInfo *npool.CoinInfo

	if paymentCoinInfo != nil {
		coinInfo = &npool.CoinInfo{
			ID:      paymentCoinInfo.ID,
			Name:    paymentCoinInfo.Name,
			PreSale: paymentCoinInfo.PreSale,
			Logo:    paymentCoinInfo.Logo,
			Unit:    paymentCoinInfo.Unit,
		}
	}

	var myPayment *npool.Payment

	if info.Payment != nil {
		myPayment = &npool.Payment{
			ID:      info.Payment.ID,
			OrderID: info.Payment.OrderID,
			Account: &npool.Account{
				ID:         account.ID,
				CoinTypeID: account.CoinTypeID,
				Address:    account.Address,
				AppID:      account.AppID,
				UserID:     account.UserID,
			},
			Amount:                info.Payment.Amount,
			CoinInfo:              coinInfo,
			State:                 info.Payment.State,
			ChainTransactionID:    info.Payment.ChainTransactionID,
			PlatformTransactionID: info.Payment.PlatformTransactionID,
		}
	}

	var goodPaying *npool.GoodPaying

	if info.GoodPaying != nil {
		goodPaying = &npool.GoodPaying{
			ID:        info.GoodPaying.ID,
			OrderID:   info.GoodPaying.OrderID,
			PaymentID: info.GoodPaying.PaymentID,
		}
	}

	return &npool.OrderDetail{
		ID:                     info.ID,
		Good:                   goodInfo,
		AppID:                  info.AppID,
		UserID:                 info.UserID,
		Units:                  info.Units,
		Discount:               discountAmount,
		SpecialReductionAmount: reductionAmount,
		DiscountCoupon:         myDiscount,
		UserSpecialReduction:   specialReduction,
		GoodPaying:             goodPaying,
		GasPayings:             gasPayings,
		Compensates:            compensates,
		OutOfGases:             outOfGases,
		Payment:                myPayment,
		Start:                  info.Start,
		End:                    info.End,
		Coupon:                 myCoupon,
	}
}

func expandDetail(ctx context.Context, info *orderpb.OrderDetail) (*npool.OrderDetail, error) {
	var coupon *inspirepb.CouponAllocatedDetail

	couponAllocated, err := grpc2.GetCouponAllocated(ctx, &inspirepb.GetCouponAllocatedDetailRequest{
		ID: info.CouponID,
	})
	invalidUUID := uuid.UUID{}.String()
	if err != nil && info.CouponID != invalidUUID {
		return nil, xerrors.Errorf("fail get coupon allocated detail: %v", err)
	}

	if couponAllocated != nil {
		coupon = couponAllocated.Info
	}

	var paymentCoinInfo *coininfopb.CoinInfo
	var accountInfo *billingpb.CoinAccountInfo

	if info.Payment != nil {
		coinInfo, err := grpc2.GetCoinInfo(ctx, &coininfopb.GetCoinInfoRequest{
			ID: info.Payment.CoinInfoID,
		})
		if err != nil {
			return nil, xerrors.Errorf("fail get payment coin info: %v", err)
		}
		paymentCoinInfo = coinInfo.Info

		account, err := grpc2.GetBillingAccount(ctx, &billingpb.GetCoinAccountRequest{
			ID: info.Payment.AccountID,
		})
		if err != nil {
			return nil, xerrors.Errorf("fail get payment address: %v", err)
		}
		accountInfo = account.Info
	}

	var discountCoupon *inspirepb.CouponAllocatedDetail

	discount, err := grpc2.GetCouponAllocated(ctx, &inspirepb.GetCouponAllocatedDetailRequest{
		ID: info.DiscountCouponID,
	})
	if err != nil && info.DiscountCouponID != invalidUUID {
		return nil, xerrors.Errorf("fail get discount coupon allocated detail: %v", err)
	}

	if discount != nil {
		discountCoupon = discount.Info
	}

	var userSpecialReduction *inspirepb.UserSpecialReduction

	userSpecial, err := grpc2.GetUserSpecialReduction(ctx, &inspirepb.GetUserSpecialReductionRequest{
		ID: info.UserSpecialReductionID,
	})
	if err != nil && info.UserSpecialReductionID != invalidUUID {
		return nil, xerrors.Errorf("fail get user special reduction: %v", err)
	}

	if userSpecial != nil {
		userSpecialReduction = userSpecial.Info
	}

	goodInfo, err := gooddetail.Get(ctx, &npool.GetGoodDetailRequest{
		ID: info.GetGoodID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get good info: %v", err)
	}

	return constructOrderDetail(
		info,
		goodInfo.Detail,
		coupon,
		paymentCoinInfo,
		accountInfo,
		discountCoupon,
		userSpecialReduction), nil
}

func GetOrderDetail(ctx context.Context, in *npool.GetOrderDetailRequest) (*npool.GetOrderDetailResponse, error) {
	orderDetail, err := grpc2.GetOrderDetail(ctx, &orderpb.GetOrderDetailRequest{
		ID: in.GetID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get order detail: %v", err)
	}

	detail, err := expandDetail(ctx, orderDetail.Detail)
	if err != nil {
		return nil, xerrors.Errorf("fail expand order detail: %v", err)
	}

	return &npool.GetOrderDetailResponse{
		Detail: detail,
	}, nil
}

func GetOrdersDetailByAppUser(ctx context.Context, in *npool.GetOrdersDetailByAppUserRequest) (*npool.GetOrdersDetailByAppUserResponse, error) {
	ordersDetail, err := grpc2.GetOrdersDetailByAppUser(ctx, &orderpb.GetOrdersDetailByAppUserRequest{
		AppID:  in.GetAppID(),
		UserID: in.GetUserID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get orders detail by app user: %v", err)
	}

	details := []*npool.OrderDetail{}
	for _, info := range ordersDetail.Details {
		detail, err := expandDetail(ctx, info)
		if err != nil {
			logger.Sugar().Warnf("cannot expand order detail: %v", err)
			continue
		}

		details = append(details, detail)
	}

	return &npool.GetOrdersDetailByAppUserResponse{
		Details: details,
	}, nil
}

func GetOrdersDetailByApp(ctx context.Context, in *npool.GetOrdersDetailByAppRequest) (*npool.GetOrdersDetailByAppResponse, error) {
	ordersDetail, err := grpc2.GetOrdersDetailByApp(ctx, &orderpb.GetOrdersDetailByAppRequest{
		AppID: in.GetAppID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get orders detail by app: %v", err)
	}

	details := []*npool.OrderDetail{}
	for _, info := range ordersDetail.Details {
		detail, err := expandDetail(ctx, info)
		if err != nil {
			logger.Sugar().Warnf("cannot expand order detail: %v", err)
			continue
		}

		details = append(details, detail)
	}

	return &npool.GetOrdersDetailByAppResponse{
		Details: details,
	}, nil
}

func GetOrdersDetailByGood(ctx context.Context, in *npool.GetOrdersDetailByGoodRequest) (*npool.GetOrdersDetailByGoodResponse, error) {
	ordersDetail, err := grpc2.GetOrdersDetailByGood(ctx, &orderpb.GetOrdersDetailByGoodRequest{
		GoodID: in.GetGoodID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get orders detail by good: %v", err)
	}

	details := []*npool.OrderDetail{}
	for _, info := range ordersDetail.Details {
		detail, err := expandDetail(ctx, info)
		if err != nil {
			logger.Sugar().Warnf("cannot expand order detail: %v", err)
			continue
		}

		details = append(details, detail)
	}

	return &npool.GetOrdersDetailByGoodResponse{
		Details: details,
	}, nil
}

func SubmitOrder(ctx context.Context, in *npool.SubmitOrderRequest) (*npool.SubmitOrderResponse, error) {
	goodInfo, err := gooddetail.Get(ctx, &npool.GetGoodDetailRequest{
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
	end := start + uint32(goodInfo.Detail.DurationDays)*24*60*60

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

	return &npool.SubmitOrderResponse{
		Info: &npool.Order{
			ID:                     myOrder.Info.ID,
			GoodID:                 myOrder.Info.GoodID,
			AppID:                  myOrder.Info.AppID,
			UserID:                 myOrder.Info.UserID,
			Units:                  myOrder.Info.Units,
			DiscountCouponID:       myOrder.Info.DiscountCouponID,
			UserSpecialReductionID: myOrder.Info.UserSpecialReductionID,
			Start:                  myOrder.Info.Start,
			End:                    myOrder.Info.End,
			CouponID:               myOrder.Info.CouponID,
		},
	}, nil
}

func CreateOrderPayment(ctx context.Context, in *npool.CreateOrderPaymentRequest) (*npool.CreateOrderPaymentResponse, error) {
	myOrder, err := GetOrderDetail(ctx, &npool.GetOrderDetailRequest{
		ID: in.GetOrderID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get order: %v", err)
	}

	goodInfo, err := gooddetail.Get(ctx, &npool.GetGoodDetailRequest{
		ID: myOrder.Detail.Good.ID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get order good info: %v", err)
	}

	// Caculate amount
	amount := float64(myOrder.Detail.Units) * goodInfo.Detail.Price

	// TODO: All should validate duration days
	// User discount info
	amount *= float64(100 - myOrder.Detail.Discount)
	amount /= float64(100)
	// Extra reduction
	amount -= myOrder.Detail.SpecialReductionAmount
	// Coupon amount
	if myOrder.Detail.Coupon != nil {
		amount -= myOrder.Detail.Coupon.Pool.Denomination
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
			AppID:       myOrder.Detail.AppID,
			UserID:      myOrder.Detail.UserID,
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
			OrderID:     myOrder.Detail.ID,
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
			OrderID:   myOrder.Detail.ID,
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
				OrderID:         myOrder.Detail.ID,
				PaymentID:       myPayment.Info.ID,
				DurationMinutes: fee.DurationDays * 24 * 60,
				FeeTypeID:       fee.ID,
			},
		})
		if err != nil {
			return nil, xerrors.Errorf("fail create fee: %v", err)
		}
	}

	orderDetail, err := GetOrderDetail(ctx, &npool.GetOrderDetailRequest{
		ID: myOrder.Detail.ID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get order detail: %v", err)
	}

	// Watch payment address and change payment state

	return &npool.CreateOrderPaymentResponse{
		Info: orderDetail.Detail,
	}, nil
}
