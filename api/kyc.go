package api

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	mw "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/kyc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetKycByAppUser(ctx context.Context, in *npool.GetKycByAppUserRequest) (*npool.GetKycByAppUserResponse, error) {
	resp, err := mw.GetByAppUser(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get kyc error: %w", err)
		return &npool.GetKycByAppUserResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}
