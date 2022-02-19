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

func (s *Server) UpdateUserWithdrawReview(ctx context.Context, in *npool.UpdateUserWithdrawReviewRequest) (*npool.UpdateUserWithdrawReviewResponse, error) {
	resp, err := withdraw.Update(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("update user withdraw error: %v", err)
		return &npool.UpdateUserWithdrawReviewResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) UpdateUserWithdrawReviewForOtherAppUser(ctx context.Context, in *npool.UpdateUserWithdrawReviewForOtherAppUserRequest) (*npool.UpdateUserWithdrawReviewForOtherAppUserResponse, error) {
	info := in.GetInfo()
	info.AppID = in.GetTargetAppID()

	resp, err := withdraw.Update(ctx, &npool.UpdateUserWithdrawReviewRequest{
		AppID:  in.GetAppID(),
		UserID: in.GetUserID(),
		Info:   info,
	})
	if err != nil {
		logger.Sugar().Errorf("update user withdraw error: %v", err)
		return &npool.UpdateUserWithdrawReviewForOtherAppUserResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &npool.UpdateUserWithdrawReviewForOtherAppUserResponse{
		Info: resp.Info,
	}, nil
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

func (s *Server) GetWithdrawAddressesByAppUser(ctx context.Context, in *npool.GetWithdrawAddressesByAppUserRequest) (*npool.GetWithdrawAddressesByAppUserResponse, error) {
	resp, err := withdrawaddress.GetByAppUser(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get user withdraw address error: %v", err)
		return &npool.GetWithdrawAddressesByAppUserResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}
