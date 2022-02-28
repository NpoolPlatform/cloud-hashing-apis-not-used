// +build !codeanalysis

package api

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	billingstate "github.com/NpoolPlatform/cloud-hashing-billing/pkg/const"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	user "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/user"
	withdraw "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/withdraw"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetCommissionByAppUser(ctx context.Context, in *npool.GetCommissionByAppUserRequest) (*npool.GetCommissionByAppUserResponse, error) {
	amount, err := user.GetCommission(in.GetAppID(), in.GetUserID())
	if err != nil {
		logger.Sugar().Errorf("get commission error: %v", err)
		return &npool.GetCommissionByAppUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	coinTypeID, err := withdraw.CommissionCoinTypeID(ctx)
	if err != nil {
		logger.Sugar().Errorf("get commission coin type id error: %v", err)
		return &npool.GetCommissionByAppUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	outcoming, err := withdraw.Outcoming(ctx, in.GetAppID(), in.GetUserID(), coinTypeID, billingstate.WithdrawTypeCommission)
	if err != nil {
		logger.Sugar().Errorf("get commission withdraw error: %v", err)
		return &npool.GetCommissionByAppUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	if amount < outcoming {
		logger.Sugar().Errorf("commission amount error")
		return &npool.GetCommissionByAppUserResponse{}, status.Error(codes.Internal, "invalid error")
	}

	return &npool.GetCommissionByAppUserResponse{
		Info: &npool.Commission{
			Total:   amount,
			Balance: amount - outcoming,
		},
	}, nil
}
