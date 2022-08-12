package fee

import (
	"context"
	"fmt"

	constant "github.com/NpoolPlatform/cloud-hashing-apis/pkg/const"
	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"
	coininfopb "github.com/NpoolPlatform/message/npool/coininfo"
	currency "github.com/NpoolPlatform/oracle-manager/pkg/middleware/currency"
)

func amount(ctx context.Context, coinTypeID string, amount float64) (float64, error) {
	coin, err := grpc2.GetCoinInfo(ctx, &coininfopb.GetCoinInfoRequest{
		ID: coinTypeID,
	})
	if err != nil {
		return 0, fmt.Errorf("invalid coin info id: %v", err)
	}

	price, err := currency.USDPrice(ctx, coin.Name)
	if err != nil {
		return 0, fmt.Errorf("cannot get usd currency for coin: %v", err)
	}

	return amount / price, nil
}

func Amount(ctx context.Context, coinTypeID string) (float64, error) {
	return amount(ctx, coinTypeID, constant.FeeUSDTAmount)
}

func ExtraAmount(ctx context.Context, coinTypeID string) (float64, error) {
	return amount(ctx, coinTypeID, constant.ExtraPayAmount)
}

// TODO: not a fix amount
func GetCurrentFee(ctx context.Context, in *npool.GetCurrentFeeRequest) (*npool.GetCurrentFeeResponse, error) {
	amount, err := Amount(ctx, in.GetCoinTypeID())
	if err != nil {
		return nil, fmt.Errorf("fail get fee amount: %v", err)
	}

	return &npool.GetCurrentFeeResponse{
		FeeAmount: amount,
	}, nil
}
