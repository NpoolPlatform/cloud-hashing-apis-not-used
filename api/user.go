// +build !codeanalysis

package api

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

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

func (s *Server) GetMyInvitations(ctx context.Context, in *npool.GetMyInvitationsRequest) (*npool.GetMyInvitationsResponse, error) {
	resp, err := user.GetMyInvitations(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get my invitations error: %w", err)
		return &npool.GetMyInvitationsResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) GetMyDirectInvitations(ctx context.Context, in *npool.GetMyDirectInvitationsRequest) (*npool.GetMyDirectInvitationsResponse, error) {
	resp, err := user.GetMyDirectInvitations(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get my direct invitations error: %w", err)
		return &npool.GetMyDirectInvitationsResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}
