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

func (s *Server) UpdatePassword(ctx context.Context, in *npool.UpdatePasswordRequest) (*npool.UpdatePasswordResponse, error) {
	resp, err := user.UpdatePassword(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("update password error: %w", err)
		return &npool.UpdatePasswordResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) UpdatePasswordByAppUser(ctx context.Context, in *npool.UpdatePasswordByAppUserRequest) (*npool.UpdatePasswordByAppUserResponse, error) {
	resp, err := user.UpdatePasswordByAppUser(ctx, in, true)
	if err != nil {
		logger.Sugar().Errorf("update password by app user error: %w", err)
		return &npool.UpdatePasswordByAppUserResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) UpdateEmailAddress(ctx context.Context, in *npool.UpdateEmailAddressRequest) (*npool.UpdateEmailAddressResponse, error) {
	resp, err := user.UpdateEmailAddress(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("update email address error: %w", err)
		return &npool.UpdateEmailAddressResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) UpdatePhoneNO(ctx context.Context, in *npool.UpdatePhoneNORequest) (*npool.UpdatePhoneNOResponse, error) {
	resp, err := user.UpdatePhoneNO(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("update phone NO error: %w", err)
		return &npool.UpdatePhoneNOResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) GetCommissionByAppUser(ctx context.Context, in *npool.GetCommissionByAppUserRequest) (*npool.GetCommissionByAppUserResponse, error) {
	amount, err := user.GetCommission(in.GetAppID(), in.GetUserID())
	if err != nil {
		logger.Sugar().Errorf("get commission error: %v", err)
		return &npool.GetCommissionByAppUserResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &npool.GetCommissionByAppUserResponse{
		Info: &npool.Commission{
			Amount: amount,
		},
	}, nil
}
