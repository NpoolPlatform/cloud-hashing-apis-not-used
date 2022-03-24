package api

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	mw "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/kyc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateKyc(ctx context.Context, in *npool.CreateKycRequest) (*npool.CreateKycResponse, error) {
	resp, err := mw.Create(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("create kyc error: %w", err)
		return &npool.CreateKycResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) UpdateKyc(ctx context.Context, in *npool.UpdateKycRequest) (*npool.UpdateKycResponse, error) {
	resp, err := mw.Update(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("update kyc error: %w", err)
		return &npool.UpdateKycResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) GetKycByAppUser(ctx context.Context, in *npool.GetKycByAppUserRequest) (*npool.GetKycByAppUserResponse, error) {
	resp, err := mw.GetByAppUser(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get kyc error: %w", err)
		return &npool.GetKycByAppUserResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) UpdateKycReview(ctx context.Context, in *npool.UpdateKycReviewRequest) (*npool.UpdateKycReviewResponse, error) {
	resp, err := mw.UpdateKycReview(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("update kyc review error: %w", err)
		return &npool.UpdateKycReviewResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}
