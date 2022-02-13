// +build !codeanalysis

package api

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/account"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreatePlatformCoinAccount(ctx context.Context, in *npool.CreatePlatformCoinAccountRequest) (*npool.CreatePlatformCoinAccountResponse, error) {
	resp, err := account.CreatePlatformCoinAccount(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("create platform coin account error: %w", err)
		return &npool.CreatePlatformCoinAccountResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) CreateUserCoinAccount(ctx context.Context, in *npool.CreateUserCoinAccountRequest) (*npool.CreateUserCoinAccountResponse, error) {
	resp, err := account.CreateUserCoinAccount(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("create user coin account error: %w", err)
		return &npool.CreateUserCoinAccountResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}
