//go:build !codeanalysis
// +build !codeanalysis

package api

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"go.opentelemetry.io/otel"

	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	constant "github.com/NpoolPlatform/cloud-hashing-apis/pkg/message/const"
	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Signup(ctx context.Context, in *npool.SignupRequest) (*npool.SignupResponse, error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "Signup")
	defer span.End()

	resp, err := user.Signup(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("sign up error: %w", err)
		return &npool.SignupResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) UpdatePassword(ctx context.Context, in *npool.UpdatePasswordRequest) (*npool.UpdatePasswordResponse, error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdatePassword")
	defer span.End()

	resp, err := user.UpdatePassword(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("update password error: %w", err)
		return &npool.UpdatePasswordResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) UpdatePasswordByAppUser(ctx context.Context, in *npool.UpdatePasswordByAppUserRequest) (*npool.UpdatePasswordByAppUserResponse, error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdatePasswordByAppUser")
	defer span.End()

	resp, err := user.UpdatePasswordByAppUser(ctx, in, true)
	if err != nil {
		logger.Sugar().Errorf("update password by app user error: %w", err)
		return &npool.UpdatePasswordByAppUserResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) UpdateEmailAddress(ctx context.Context, in *npool.UpdateEmailAddressRequest) (*npool.UpdateEmailAddressResponse, error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateEmailAddress")
	defer span.End()

	resp, err := user.UpdateEmailAddress(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("update email address error: %w", err)
		return &npool.UpdateEmailAddressResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) UpdatePhoneNO(ctx context.Context, in *npool.UpdatePhoneNORequest) (*npool.UpdatePhoneNOResponse, error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdatePhoneNO")
	defer span.End()

	resp, err := user.UpdatePhoneNO(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("update phone NO error: %w", err)
		return &npool.UpdatePhoneNOResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) UpdateAccount(ctx context.Context, in *npool.UpdateAccountRequest) (*npool.UpdateAccountResponse, error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateAccount")
	defer span.End()

	resp, err := user.UpdateAccount(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("update account error: %v", err)
		return &npool.UpdateAccountResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) UpdateAppUserExtra(ctx context.Context, in *npool.UpdateAppUserExtraRequest) (*npool.UpdateAppUserExtraResponse, error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateAppUserExtra")
	defer span.End()

	resp, err := user.UpdateAppUserExtra(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("update app user extra error: %v", err)
		return &npool.UpdateAppUserExtraResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) CreateAppUserExtra(ctx context.Context, in *npool.CreateAppUserExtraRequest) (*npool.CreateAppUserExtraResponse, error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateAppUserExtra")
	defer span.End()

	resp, err := user.CreateAppUserExtra(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("create app user extra error: %v", err)
		return &npool.CreateAppUserExtraResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}
