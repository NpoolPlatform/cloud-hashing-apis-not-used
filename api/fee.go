//go:build !codeanalysis
// +build !codeanalysis

package api

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/fee"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetCurrentFee(ctx context.Context, in *npool.GetCurrentFeeRequest) (*npool.GetCurrentFeeResponse, error) {
	resp, err := fee.GetCurrentFee(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get gas fee error: %w", err)
		return &npool.GetCurrentFeeResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}
