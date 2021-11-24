// +build !codeanalysis

package api

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	"github.com/NpoolPlatform/cloud-hashing-apis/message/npool"

	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/order" //nolint

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetOrderDetail(ctx context.Context, in *npool.GetOrderDetailRequest) (*npool.GetOrderDetailResponse, error) {
	resp, err := order.GetOrderDetail(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get order detail error: %w", err)
		return &npool.GetOrderDetailResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) GetOrdersDetailByAppUser(ctx context.Context, in *npool.GetOrdersDetailByAppUserRequest) (*npool.GetOrdersDetailByAppUserResponse, error) {
	resp, err := order.GetOrdersDetailByAppUser(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get order detail by app user error: %w", err)
		return &npool.GetOrdersDetailByAppUserResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) GetOrdersDetailByApp(ctx context.Context, in *npool.GetOrdersDetailByAppRequest) (*npool.GetOrdersDetailByAppResponse, error) {
	resp, err := order.GetOrdersDetailByApp(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get order detail by app error: %w", err)
		return &npool.GetOrdersDetailByAppResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) GetOrdersDetailByGood(ctx context.Context, in *npool.GetOrdersDetailByGoodRequest) (*npool.GetOrdersDetailByGoodResponse, error) {
	resp, err := order.GetOrdersDetailByGood(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get order detail by good error: %w", err)
		return &npool.GetOrdersDetailByGoodResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) SubmitOrder(ctx context.Context, in *npool.SubmitOrderRequest) (*npool.SubmitOrderResponse, error) {
	resp, err := order.SubmitOrder(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("submit order error: %w", err)
		return &npool.SubmitOrderResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}
