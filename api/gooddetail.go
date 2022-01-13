// +build !codeanalysis

package api

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/good-detail" //nolint

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetGoods(ctx context.Context, in *npool.GetGoodsRequest) (*npool.GetGoodsResponse, error) {
	resp, err := gooddetail.GetAll(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get good detail all error: %w", err)
		return &npool.GetGoodsResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) GetGood(ctx context.Context, in *npool.GetGoodRequest) (*npool.GetGoodResponse, error) {
	resp, err := gooddetail.Get(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get good detail all error: %w", err)
		return &npool.GetGoodResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func GetRecommendGoodsByApp(ctx context.Context, in *npool.GetRecommendGoodsByAppRequest) (*npool.GetRecommendGoodsByAppResponse, error) {
	resp, err := gooddetail.GetRecommendsByApp(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get recommend good by app error: %v", err)
		return &npool.GetRecommendGoodsByAppResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}
