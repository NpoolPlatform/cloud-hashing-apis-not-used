// +build !codeanalysis

package api

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/referral"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetMyInvitations(ctx context.Context, in *npool.GetMyInvitationsRequest) (*npool.GetMyInvitationsResponse, error) {
	resp, err := referral.GetMyInvitations(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get my invitations error: %w", err)
		return &npool.GetMyInvitationsResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) GetMyDirectInvitations(ctx context.Context, in *npool.GetMyDirectInvitationsRequest) (*npool.GetMyDirectInvitationsResponse, error) {
	resp, err := referral.GetMyDirectInvitations(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get my direct invitations error: %w", err)
		return &npool.GetMyDirectInvitationsResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}
