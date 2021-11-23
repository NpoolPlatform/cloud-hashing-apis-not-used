package order

import (
	"context"

	"github.com/NpoolPlatform/cloud-hashing-apis/message/npool"
	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"

	inspirepb "github.com/NpoolPlatform/cloud-hashing-inspire/message/npool"
	orderpb "github.com/NpoolPlatform/cloud-hashing-order/message/npool"
	coininfopb "github.com/NpoolPlatform/message/npool/coininfo"

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
	orderDetail, err := grpc.GetOrderDetail(ctx, &orderpb.GetOrderDetailRequest{
		ID: in.GetID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get order detail: %v", err)
	}

	couponAllocated, err := grpc.GetCouponAllocated(ctx, &inspirepb.GetCouponAllocatedDetailRequest{
		ID: orderDetail.Detail.CouponID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get coupon allocated detail: %v", err)
	}

	coinInfo, err := grpc.GetCoinInfo(ctx, &coininfopb.GetCoinInfoRequest{
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
	return nil, nil
}
