package account

import (
	"context"
	"fmt"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	billingpb "github.com/NpoolPlatform/message/npool/cloud-hashing-billing"
	coininfopb "github.com/NpoolPlatform/message/npool/coininfo"
	sphinxproxypb "github.com/NpoolPlatform/message/npool/sphinxproxy"
)

func CreatePlatformCoinAccount(ctx context.Context, in *npool.CreatePlatformCoinAccountRequest) (*npool.CreatePlatformCoinAccountResponse, error) {
	coinInfo, err := grpc2.GetCoinInfo(ctx, &coininfopb.GetCoinInfoRequest{
		ID: in.GetCoinTypeID(),
	})
	if err != nil {
		return nil, fmt.Errorf("invalid coin info id: %v", err)
	}
	if coinInfo == nil {
		return nil, fmt.Errorf("invalid coin info id")
	}

	address, err := grpc2.CreateCoinAddress(ctx, &sphinxproxypb.CreateWalletRequest{
		Name: coinInfo.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("fail create wallet address: %v", err)
	}

	account, err := grpc2.CreateBillingAccount(ctx, &billingpb.CreateCoinAccountRequest{
		Info: &billingpb.CoinAccountInfo{
			CoinTypeID:             coinInfo.ID,
			Address:                address.Address,
			PlatformHoldPrivateKey: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("fail create billing account: %v", err)
	}

	return &npool.CreatePlatformCoinAccountResponse{
		Info: account,
	}, nil
}

func CreateUserCoinAccount(ctx context.Context, in *npool.CreateUserCoinAccountRequest) (*npool.CreateUserCoinAccountResponse, error) {
	coinInfo, err := grpc2.GetCoinInfo(ctx, &coininfopb.GetCoinInfoRequest{
		ID: in.GetInfo().GetCoinTypeID(),
	})
	if err != nil {
		return nil, fmt.Errorf("invalid coin info id: %v", err)
	}
	if coinInfo == nil {
		return nil, fmt.Errorf("invalid coin info id")
	}

	_, err = grpc2.GetBalance(ctx, &sphinxproxypb.GetBalanceRequest{
		Name:    coinInfo.Name,
		Address: in.GetInfo().GetAddress(),
	})
	if err != nil {
		return nil, fmt.Errorf("fail get wallet balance: %v", err)
	}

	info := in.GetInfo()
	info.PlatformHoldPrivateKey = false

	account, err := grpc2.CreateBillingAccount(ctx, &billingpb.CreateCoinAccountRequest{
		Info: info,
	})
	if err != nil {
		return nil, fmt.Errorf("fail create billing account: %v", err)
	}

	return &npool.CreateUserCoinAccountResponse{
		Info: account,
	}, nil
}
