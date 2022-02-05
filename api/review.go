package api

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	mw "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/review"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetKycReviews(ctx context.Context, in *npool.GetKycReviewsRequest) (*npool.GetKycReviewsResponse, error) {
	resp, err := mw.GetKycReviews(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get kyc reviews error: %v", err)
		return &npool.GetKycReviewsResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) GetKycReviewsByOtherApp(ctx context.Context, in *npool.GetKycReviewsByOtherAppRequest) (*npool.GetKycReviewsByOtherAppResponse, error) {
	resp, err := mw.GetKycReviews(ctx, &npool.GetKycReviewsRequest{
		AppID: in.GetTargetAppID(),
	})
	if err != nil {
		logger.Sugar().Errorf("get kyc reviews by other app error: %v", err)
		return &npool.GetKycReviewsByOtherAppResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &npool.GetKycReviewsByOtherAppResponse{
		Infos: resp.Infos,
	}, nil
}

func (s *Server) GetGoodReviews(ctx context.Context, in *npool.GetGoodReviewsRequest) (*npool.GetGoodReviewsResponse, error) {
	resp, err := mw.GetGoodReviews(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get good reviews error: %v", err)
		return &npool.GetGoodReviewsResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}
