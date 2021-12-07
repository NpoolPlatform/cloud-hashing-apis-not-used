// +build !codeanalysis

package api

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	"github.com/NpoolPlatform/cloud-hashing-apis/message/npool"

	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/user" //nolint

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Signup(ctx context.Context, in *npool.SignupRequest) (*npool.SignupResponse, error) {
	resp, err := user.Signup(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("sign up error: %w", err)
		return &npool.SignupResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}
