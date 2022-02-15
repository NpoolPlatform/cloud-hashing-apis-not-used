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

func (s *Server) GetKycReviewsByApp(ctx context.Context, in *npool.GetKycReviewsByAppRequest) (*npool.GetKycReviewsByAppResponse, error) {
	resp, err := mw.GetKycReviews(ctx, &npool.GetKycReviewsRequest{
		AppID: in.GetAppID(),
	})
	if err != nil {
		logger.Sugar().Errorf("get kyc reviews by app error: %v", err)
		return &npool.GetKycReviewsByAppResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &npool.GetKycReviewsByAppResponse{
		Infos: resp.Infos,
	}, nil
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

func (s *Server) GetWithdrawReviews(ctx context.Context, in *npool.GetWithdrawReviewsRequest) (*npool.GetWithdrawReviewsResponse, error) {
	resp, err := mw.GetWithdrawReviews(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get kyc reviews error: %v", err)
		return &npool.GetWithdrawReviewsResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) GetWithdrawReviewsByApp(ctx context.Context, in *npool.GetWithdrawReviewsByAppRequest) (*npool.GetWithdrawReviewsByAppResponse, error) {
	resp, err := mw.GetWithdrawReviews(ctx, &npool.GetWithdrawReviewsRequest{
		AppID: in.GetAppID(),
	})
	if err != nil {
		logger.Sugar().Errorf("get kyc reviews by app error: %v", err)
		return &npool.GetWithdrawReviewsByAppResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &npool.GetWithdrawReviewsByAppResponse{
		Infos: resp.Infos,
	}, nil
}

func (s *Server) GetWithdrawReviewsByOtherApp(ctx context.Context, in *npool.GetWithdrawReviewsByOtherAppRequest) (*npool.GetWithdrawReviewsByOtherAppResponse, error) {
	resp, err := mw.GetWithdrawReviews(ctx, &npool.GetWithdrawReviewsRequest{
		AppID: in.GetTargetAppID(),
	})
	if err != nil {
		logger.Sugar().Errorf("get kyc reviews by other app error: %v", err)
		return &npool.GetWithdrawReviewsByOtherAppResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &npool.GetWithdrawReviewsByOtherAppResponse{
		Infos: resp.Infos,
	}, nil
}

func (s *Server) GetWithdrawAddressReviews(ctx context.Context, in *npool.GetWithdrawAddressReviewsRequest) (*npool.GetWithdrawAddressReviewsResponse, error) {
	return nil, nil
}

func (s *Server) GetWithdrawAddressReviewsByApp(ctx context.Context, in *npool.GetWithdrawAddressReviewsByAppRequest) (*npool.GetWithdrawAddressReviewsByAppResponse, error) {
	return nil, nil
}

func (s *Server) GetWithdrawAddressReviewsByOtherApp(ctx context.Context, in *npool.GetWithdrawAddressReviewsByOtherAppRequest) (*npool.GetWithdrawAddressReviewsByOtherAppResponse, error) {
	return nil, nil
}
