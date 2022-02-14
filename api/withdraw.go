// +build !codeanalysis

package api

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	mw "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/withdraw"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) SubmitUserWithdraw(ctx context.Context, in *npool.SubmitUserWithdrawRequest) (*npool.SubmitUserWithdrawResponse, error) {
	resp, err := mw.Create(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("create user withdraw error: %w", err)
		return &npool.SubmitUserWithdrawResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) UpdateUserWithdrawReview(ctx context.Context, in *npool.UpdateUserWithdrawReviewRequest) (*npool.UpdateUserWithdrawReviewResponse, error) {
	resp, err := mw.Update(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("update user withdraw error: %w", err)
		return &npool.UpdateUserWithdrawReviewResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) GetUserWithdrawsByAppUser(ctx context.Context, in *npool.GetUserWithdrawsByAppUserRequest) (*npool.GetUserWithdrawsByAppUserResponse, error) {
	resp, err := mw.GetByAppUser(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get user withdraw error: %w", err)
		return &npool.GetUserWithdrawsByAppUserResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}
