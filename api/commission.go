// +build !codeanalysis

package api

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	billingstate "github.com/NpoolPlatform/cloud-hashing-billing/pkg/const"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	commission "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/commission"
	withdraw "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/withdraw"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetCommissionByAppUser(ctx context.Context, in *npool.GetCommissionByAppUserRequest) (*npool.GetCommissionByAppUserResponse, error) {
	amount, err := commission.GetCommission(ctx, in.GetAppID(), in.GetUserID())
	if err != nil {
		logger.Sugar().Errorf("get commission error: %v", err)
		return &npool.GetCommissionByAppUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	coinTypeID, err := withdraw.CommissionCoinTypeID(ctx)
	if err != nil {
		logger.Sugar().Errorf("get commission coin type id error: %v", err)
		return &npool.GetCommissionByAppUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	outcoming, err := withdraw.Outcoming(ctx, in.GetAppID(), in.GetUserID(), coinTypeID, billingstate.WithdrawTypeCommission, true)
	if err != nil {
		logger.Sugar().Errorf("get commission withdraw error: %v", err)
		return &npool.GetCommissionByAppUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	if amount < outcoming {
		logger.Sugar().Errorf("commission amount error %v < %v", amount, outcoming)
		return &npool.GetCommissionByAppUserResponse{}, status.Error(codes.Internal, "invalid error")
	}

	return &npool.GetCommissionByAppUserResponse{
		Info: &npool.Commission{
			Total:   amount,
			Balance: amount - outcoming,
		},
	}, nil
}

func (s *Server) GetGoodCommissions(ctx context.Context, in *npool.GetGoodCommissionsRequest) (*npool.GetGoodCommissionsResponse, error) {
	commissions, err := commission.GetGoodCommissions(ctx, in.GetAppID(), in.GetUserID())
	if err != nil {
		logger.Sugar().Errorf("get good commission error: %v", err)
		return &npool.GetGoodCommissionsResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &npool.GetGoodCommissionsResponse{
		Infos: commissions,
	}, nil
}

func (s *Server) GetUserGoodCommissions(ctx context.Context, in *npool.GetUserGoodCommissionsRequest) (*npool.GetUserGoodCommissionsResponse, error) {
	commissions, err := commission.GetGoodCommissions(ctx, in.GetAppID(), in.GetTargetUserID())
	if err != nil {
		logger.Sugar().Errorf("get good commission error: %v", err)
		return &npool.GetUserGoodCommissionsResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &npool.GetUserGoodCommissionsResponse{
		Infos: commissions,
	}, nil
}

func (s *Server) GetAmountSettings(ctx context.Context, in *npool.GetAmountSettingsRequest) (*npool.GetAmountSettingsResponse, error) {
	settings, err := commission.GetAmountSettings(ctx, in.GetAppID(), in.GetUserID())
	if err != nil {
		logger.Sugar().Errorf("get amount settings error: %v", err)
		return &npool.GetAmountSettingsResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &npool.GetAmountSettingsResponse{
		Infos: settings,
	}, nil
}

func (s *Server) CreateAmountSetting(ctx context.Context, in *npool.CreateAmountSettingRequest) (*npool.CreateAmountSettingResponse, error) {
	settings, err := commission.CreateAmountSetting(ctx, in.GetAppID(), in.GetUserID(), in.GetTargetUserID(), in.GetInfo())
	if err != nil {
		logger.Sugar().Errorf("create amount settings error: %v", err)
		return &npool.CreateAmountSettingResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &npool.CreateAmountSettingResponse{
		Infos: settings,
	}, nil
}
