//go:build !codeanalysis
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

func (s *Server) CreateInvitationCode(ctx context.Context, in *npool.CreateInvitationCodeRequest) (*npool.CreateInvitationCodeResponse, error) {
	code, err := referral.CreateInvitationCode(
		ctx,
		in.GetAppID(), in.GetUserID(), in.GetTargetUserID(), in.GetLangID(),
		in.GetInviterName(), in.GetInviteeName(),
		in.GetInfo())
	if err != nil {
		logger.Sugar().Errorf("create invitation code error: %w", err)
		return &npool.CreateInvitationCodeResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &npool.CreateInvitationCodeResponse{
		Info: code,
	}, nil
}
