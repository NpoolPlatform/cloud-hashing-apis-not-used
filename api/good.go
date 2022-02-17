// +build !codeanalysis

package api

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	good "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/good"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetGoods(ctx context.Context, in *npool.GetGoodsRequest) (*npool.GetGoodsResponse, error) {
	resp, err := good.GetAll(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get good all error: %v", err)
		return &npool.GetGoodsResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) GetGoodsByApp(ctx context.Context, in *npool.GetGoodsByAppRequest) (*npool.GetGoodsByAppResponse, error) {
	resp, err := good.GetByApp(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get good all error: %v", err)
		return &npool.GetGoodsByAppResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) CreateGood(ctx context.Context, in *npool.CreateGoodRequest) (*npool.CreateGoodResponse, error) {
	resp, err := good.Create(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("create good error: %v", err)
		return &npool.CreateGoodResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) GetGood(ctx context.Context, in *npool.GetGoodRequest) (*npool.GetGoodResponse, error) {
	resp, err := good.Get(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get good error: %v", err)
		return &npool.GetGoodResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func GetRecommendGoodsByApp(ctx context.Context, in *npool.GetRecommendGoodsByAppRequest) (*npool.GetRecommendGoodsByAppResponse, error) {
	resp, err := good.GetRecommendsByApp(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get recommend good by app error: %v", err)
		return &npool.GetRecommendGoodsByAppResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}
