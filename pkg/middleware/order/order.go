package order

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	"github.com/NpoolPlatform/cloud-hashing-apis/message/npool"
	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	gooddetail "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/good-detail" //nolint

	billingpb "github.com/NpoolPlatform/cloud-hashing-billing/message/npool"
	inspirepb "github.com/NpoolPlatform/cloud-hashing-inspire/message/npool"
	orderpb "github.com/NpoolPlatform/cloud-hashing-order/message/npool"
	coininfopb "github.com/NpoolPlatform/message/npool/coininfo"
	tradingpb "github.com/NpoolPlatform/message/npool/trading"

	"golang.org/x/xerrors"
)

func constructOrderDetail(
	info *orderpb.OrderDetail,
	coupon *inspirepb.CouponAllocatedDetail,
	paymentCoinInfo *coininfopb.CoinInfo) *npool.OrderDetail {
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

	return &npool.OrderDetail{
		ID:                     info.ID,
		GoodID:                 info.GoodID,
		AppID:                  info.AppID,
		UserID:                 info.UserID,
		Units:                  info.Units,
		Discount:               info.Discount,
		SpecialReductionAmount: info.SpecialReductionAmount,
		GoodPaying: &npool.GoodPaying{
			ID:        info.GoodPaying.ID,
			OrderID:   info.GoodPaying.OrderID,
			PaymentID: info.GoodPaying.PaymentID,
		},
		GasPayings:  gasPayings,
		Compensates: compensates,
		OutOfGases:  outOfGases,
		Payment: &npool.Payment{
			ID:        info.Payment.ID,
			OrderID:   info.Payment.OrderID,
			AccountID: info.Payment.AccountID,
			Amount:    info.Payment.Amount,
			CoinInfo: &npool.CoinInfo{
				ID:      paymentCoinInfo.ID,
				Name:    paymentCoinInfo.Name,
				PreSale: paymentCoinInfo.PreSale,
				Logo:    "",
				Unit:    paymentCoinInfo.Unit,
			},
			State:                 info.Payment.State,
			ChainTransactionID:    info.Payment.ChainTransactionID,
			PlatformTransactionID: info.Payment.PlatformTransactionID,
		},
		Start: info.Start,
		End:   info.End,
		Coupon: &npool.Coupon{
			ID:     coupon.ID,
			UserID: coupon.UserID,
			AppID:  coupon.AppID,
			Used:   coupon.Used,
			Pool: &npool.CouponPool{
				ID:           coupon.Coupon.ID,
				AppID:        coupon.Coupon.AppID,
				Denomination: coupon.Coupon.Denomination,
				Start:        coupon.Coupon.Start,
				DurationDays: coupon.Coupon.DurationDays,
				Message:      coupon.Coupon.Message,
				Name:         coupon.Coupon.Message,
			},
		},
	}
}

func expandDetail(ctx context.Context, info *orderpb.OrderDetail) (*npool.OrderDetail, error) {
	couponAllocated, err := grpc2.GetCouponAllocated(ctx, &inspirepb.GetCouponAllocatedDetailRequest{
		ID: info.CouponID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get coupon allocated detail: %v", err)
	}

	coinInfo, err := grpc2.GetCoinInfo(ctx, &coininfopb.GetCoinInfoRequest{
		ID: info.Payment.CoinInfoID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get payment coin info: %v", err)
	}

	return constructOrderDetail(info, couponAllocated.Info, coinInfo.Info), nil
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
	goodInfo, err := gooddetail.Get(ctx, in.GetGoodID())
	if err != nil {
		return nil, xerrors.Errorf("fail get order good info: %v", err)
	}

	// Validate app id
	// Validate user id
	// Validate coupon id
	// Validate fee ids
	// Caculate amount
	amount := 1273.0

	// Validate payment coin info id
	coinInfo, err := grpc2.GetCoinInfo(ctx, &coininfopb.GetCoinInfoRequest{
		ID: in.GetPaymentCoinTypeID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("invalid coin info id: %v", err)
	}

	// Check if idle address is available
	idle := false

	// Generate transaction address
	address, err := grpc2.CreateCoinAddress(ctx, &tradingpb.CreateWalletRequest{
		CoinName: coinInfo.Info.Name,
		UUID:     in.GetUserID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail create wallet address: %v", err)
	}

	// Create billing account
	account, err := grpc2.CreateBillingAccount(ctx, &billingpb.CreateCoinAccountRequest{
		Info: &billingpb.CoinAccountInfo{
			CoinTypeID:  in.GetPaymentCoinTypeID(),
			Address:     address.Info.Address,
			GeneratedBy: "platform",
			UsedFor:     "paying",
			AppID:       in.GetAppID(),
			UserID:      in.GetUserID(),
		},
	})
	if err != nil {
		return nil, xerrors.Errorf("fail create billing account: %v", err)
	}

	balanceAmount := float64(0)

	if !idle {
		balance, err := grpc2.GetWalletBalance(ctx, &tradingpb.GetWalletBalanceRequest{
			Info: &tradingpb.EntAccount{
				CoinName: coinInfo.Info.Name,
				Address:  account.Info.Address,
			},
		})
		if err != nil {
			return nil, xerrors.Errorf("fail get wallet balance: %v", err)
		}
		balanceAmount = balance.AmountFloat64
	}

	start := (goodInfo.Start + 24*60*60) / 24 / 60 / 60 * 24 * 60 * 60
	end := start + uint32(goodInfo.DurationDays)*24*60*60

	// Generate order
	myOrder, err := grpc2.CreateOrder(ctx, &orderpb.CreateOrderRequest{
		Info: &orderpb.Order{
			GoodID:   in.GetGoodID(),
			AppID:    in.GetAppID(),
			UserID:   in.GetUserID(),
			Units:    in.GetUnits(),
			Start:    start,
			End:      end,
			CouponID: in.GetCouponID(),
		},
	})
	if err != nil {
		return nil, xerrors.Errorf("fail create order: %v", err)
	}

	// Generate payment
	myPayment, err := grpc2.CreatePayment(ctx, &orderpb.CreatePaymentRequest{
		Info: &orderpb.Payment{
			OrderID:     myOrder.Info.ID,
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
			OrderID:   myOrder.Info.ID,
			PaymentID: myPayment.Info.ID,
		},
	})
	if err != nil {
		return nil, xerrors.Errorf("fail create good paying: %v", err)
	}

	// Generate gas payings
	for _, fee := range in.Fees {
		_, err := grpc2.CreateGasPaying(ctx, &orderpb.CreateGasPayingRequest{
			Info: &orderpb.GasPaying{
				OrderID:         myOrder.Info.ID,
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
		ID: myOrder.Info.ID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get order detail: %v", err)
	}

	// Watch payment address and change payment state

	return &npool.SubmitOrderResponse{
		Info: orderDetail.Detail,
	}, nil
}
