// +build !codeanalysis

package api

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	"github.com/NpoolPlatform/cloud-hashing-apis/message/npool"

	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/good-detail" //nolint

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetGoodsDetail(ctx context.Context, in *npool.GetGoodsDetailRequest) (*npool.GetGoodsDetailResponse, error) {
	resp, err := gooddetail.GetAll(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get good detail all error: %w", err)
		return &npool.GetGoodsDetailResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) GetGoodDetail(ctx context.Context, in *npool.GetGoodDetailRequest) (*npool.GetGoodDetailResponse, error) {
	resp, err := gooddetail.Get(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get good detail all error: %w", err)
		return &npool.GetGoodDetailResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func GetRecommendGoodsDetailByApp(ctx context.Context, in *npool.GetRecommendGoodsDetailByAppRequest) (*npool.GetRecommendGoodsDetailByAppResponse, error) {
	resp, err := gooddetail.GetRecommendsByApp(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get recommend good by app error: %v", err)
		return &npool.GetRecommendGoodsDetailByAppResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}
