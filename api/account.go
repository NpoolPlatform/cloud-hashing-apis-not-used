//go:build !codeanalysis
// +build !codeanalysis

package api

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"go.opentelemetry.io/otel"

	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	constant "github.com/NpoolPlatform/cloud-hashing-apis/pkg/message/const"
	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/account"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreatePlatformCoinAccount(ctx context.Context, in *npool.CreatePlatformCoinAccountRequest) (*npool.CreatePlatformCoinAccountResponse, error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreatePlatformCoinAccount")
	defer span.End()

	resp, err := account.CreatePlatformCoinAccount(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("create platform coin account error: %w", err)
		return &npool.CreatePlatformCoinAccountResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) CreateUserCoinAccount(ctx context.Context, in *npool.CreateUserCoinAccountRequest) (*npool.CreateUserCoinAccountResponse, error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateUserCoinAccount")
	defer span.End()

	resp, err := account.CreateUserCoinAccount(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("create user coin account error: %w", err)
		return &npool.CreateUserCoinAccountResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}
