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

func (s *Server) SubmitOrder(ctx context.Context, in *npool.SubmitOrderRequest) (*npool.SubmitOrderResponse, error) {
	resp, err := order.SubmitOrder(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("submit order error: %w", err)
		return &npool.SubmitOrderResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}
