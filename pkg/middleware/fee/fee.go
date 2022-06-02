package fee

import (
	"context"

	constant "github.com/NpoolPlatform/cloud-hashing-apis/pkg/const"
	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"
	coininfopb "github.com/NpoolPlatform/message/npool/coininfo"
	currency "github.com/NpoolPlatform/staker-manager/pkg/middleware/currency"

	"golang.org/x/xerrors"
)

func Amount(ctx context.Context, coinTypeID string) (float64, error) {
	coin, err := grpc2.GetCoinInfo(ctx, &coininfopb.GetCoinInfoRequest{
		ID: coinTypeID,
	})
	if err != nil {
		return 0, xerrors.Errorf("invalid coin info id: %v", err)
	}

	price, err := currency.USDPrice(ctx, coin.Name)
	if err != nil {
		return 0, xerrors.Errorf("cannot get usd currency for coin: %v", err)
	}

	return constant.FeeUSDTAmount / price, nil
}

// TODO: not a fix amount
func GetCurrentFee(ctx context.Context, in *npool.GetCurrentFeeRequest) (*npool.GetCurrentFeeResponse, error) {
	amount, err := Amount(ctx, in.GetCoinTypeID())
	if err != nil {
		return nil, xerrors.Errorf("fail get fee amount: %v", err)
	}

	return &npool.GetCurrentFeeResponse{
		FeeAmount: amount,
	}, nil
}
