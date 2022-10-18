package api

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	mw "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/review"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetWithdrawAddressReviews(ctx context.Context, in *npool.GetWithdrawAddressReviewsRequest) (*npool.GetWithdrawAddressReviewsResponse, error) {
	resp, err := mw.GetWithdrawAddressReviews(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get withdraw address reviews error: %v", err)
		return &npool.GetWithdrawAddressReviewsResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) GetWithdrawAddressReviewsByApp(ctx context.Context, in *npool.GetWithdrawAddressReviewsByAppRequest) (*npool.GetWithdrawAddressReviewsByAppResponse, error) {
	resp, err := mw.GetWithdrawAddressReviews(ctx, &npool.GetWithdrawAddressReviewsRequest{
		AppID: in.GetAppID(),
	})
	if err != nil {
		logger.Sugar().Errorf("get withdraw address reviews by app error: %v", err)
		return &npool.GetWithdrawAddressReviewsByAppResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &npool.GetWithdrawAddressReviewsByAppResponse{
		Infos: resp.Infos,
	}, nil
}

func (s *Server) GetWithdrawAddressReviewsByOtherApp(ctx context.Context, in *npool.GetWithdrawAddressReviewsByOtherAppRequest) (*npool.GetWithdrawAddressReviewsByOtherAppResponse, error) {
	resp, err := mw.GetWithdrawAddressReviews(ctx, &npool.GetWithdrawAddressReviewsRequest{
		AppID: in.GetTargetAppID(),
	})
	if err != nil {
		logger.Sugar().Errorf("get withdraw address reviews by other app error: %v", err)
		return &npool.GetWithdrawAddressReviewsByOtherAppResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &npool.GetWithdrawAddressReviewsByOtherAppResponse{
		Infos: resp.Infos,
	}, nil
}
