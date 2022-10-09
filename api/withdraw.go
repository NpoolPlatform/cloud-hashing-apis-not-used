//go:build !codeanalysis
// +build !codeanalysis

package api

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	withdraw "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/withdraw"
	withdrawaddress "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/withdrawaddress"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) SubmitUserWithdraw(ctx context.Context, in *npool.SubmitUserWithdrawRequest) (*npool.SubmitUserWithdrawResponse, error) {
	resp, err := withdraw.Create(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("create user withdraw error: %v", err)
		return &npool.SubmitUserWithdrawResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) GetUserWithdrawsByAppUser(ctx context.Context, in *npool.GetUserWithdrawsByAppUserRequest) (*npool.GetUserWithdrawsByAppUserResponse, error) {
	resp, err := withdraw.GetByAppUser(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get user withdraw error: %v", err)
		return &npool.GetUserWithdrawsByAppUserResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) SetWithdrawAddress(ctx context.Context, in *npool.SetWithdrawAddressRequest) (*npool.SetWithdrawAddressResponse, error) {
	resp, err := withdrawaddress.Set(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("create user withdraw address error: %v", err)
		return &npool.SetWithdrawAddressResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) DeleteWithdrawAddress(ctx context.Context, in *npool.DeleteWithdrawAddressRequest) (*npool.DeleteWithdrawAddressResponse, error) {
	resp, err := withdrawaddress.Delete(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("delete user withdraw address error: %v", err)
		return &npool.DeleteWithdrawAddressResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) GetWithdrawAddressesByAppUser(ctx context.Context, in *npool.GetWithdrawAddressesByAppUserRequest) (*npool.GetWithdrawAddressesByAppUserResponse, error) {
	resp, err := withdrawaddress.GetByAppUser(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get user withdraw address error: %v", err)
		return &npool.GetWithdrawAddressesByAppUserResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

// func (s *Server) UpdateWithdrawAddressReview(ctx context.Context, in *npool.UpdateWithdrawAddressReviewRequest) (*npool.UpdateWithdrawAddressReviewResponse, error) {
//	resp, err := withdrawaddress.UpdateWithdrawUpdateAddressReview(ctx, in)
//	if err != nil {
//		logger.Sugar().Errorf("update withdraw addresses review error: %w", err)
//		return &npool.UpdateWithdrawAddressReviewResponse{}, status.Error(codes.Internal, err.Error())
//	}
//	return resp, nil
// }
