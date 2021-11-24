package order

import (
	"context"

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

func GetOrderDetail(ctx context.Context, in *npool.GetOrderDetailRequest) (*npool.GetOrderDetailResponse, error) {
	orderDetail, err := grpc2.GetOrderDetail(ctx, &orderpb.GetOrderDetailRequest{
		ID: in.GetID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get order detail: %v", err)
	}

	couponAllocated, err := grpc2.GetCouponAllocated(ctx, &inspirepb.GetCouponAllocatedDetailRequest{
		ID: orderDetail.Detail.CouponID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get coupon allocated detail: %v", err)
	}

	coinInfo, err := grpc2.GetCoinInfo(ctx, &coininfopb.GetCoinInfoRequest{
		ID: orderDetail.Detail.Payment.CoinInfoID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get payment coin info: %v", err)
	}

	return &npool.GetOrderDetailResponse{
		Detail: constructOrderDetail(orderDetail.Detail, couponAllocated.Info, coinInfo.Info),
	}, nil
}

func GetOrdersDetailByAppUser(ctx context.Context, in *npool.GetOrdersDetailByAppUserRequest) (*npool.GetOrdersDetailByAppUserResponse, error) {
	return nil, nil
}

func GetOrdersDetailByApp(ctx context.Context, in *npool.GetOrdersDetailByAppRequest) (*npool.GetOrdersDetailByAppResponse, error) {
	return nil, nil
}

func GetOrdersDetailByGood(ctx context.Context, in *npool.GetOrdersDetailByGoodRequest) (*npool.GetOrdersDetailByGoodResponse, error) {
	return nil, nil
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
			UsedFor:     "payment",
			AppID:       in.GetAppID(),
			UserID:      in.GetUserID(),
		},
	})
	if err != nil {
		return nil, xerrors.Errorf("fail create billing account: %v", err)
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
			OrderID:    myOrder.Info.ID,
			AccountID:  account.Info.ID,
			Amount:     amount,
			CoinInfoID: in.GetPaymentCoinTypeID(),
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
