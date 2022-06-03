// +build !codeanalysis

package api

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/order" //nolint

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetOrder(ctx context.Context, in *npool.GetOrderRequest) (*npool.GetOrderResponse, error) {
	resp, err := order.GetOrder(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get order detail error: %v", err)
		return &npool.GetOrderResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) GetOrdersByAppUser(ctx context.Context, in *npool.GetOrdersByAppUserRequest) (*npool.GetOrdersByAppUserResponse, error) {
	resp, err := order.GetOrdersByAppUser(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get order detail by app user error: %v", err)
		return &npool.GetOrdersByAppUserResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) GetOrdersByApp(ctx context.Context, in *npool.GetOrdersByAppRequest) (*npool.GetOrdersByAppResponse, error) {
	resp, err := order.GetOrdersByApp(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get order detail by app error: %v", err)
		return &npool.GetOrdersByAppResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) GetOrdersByOtherApp(ctx context.Context, in *npool.GetOrdersByOtherAppRequest) (*npool.GetOrdersByOtherAppResponse, error) {
	resp, err := order.GetOrdersByApp(ctx, &npool.GetOrdersByAppRequest{
		AppID: in.GetTargetAppID(),
	})
	if err != nil {
		logger.Sugar().Errorf("get order detail by app error: %v", err)
		return &npool.GetOrdersByOtherAppResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &npool.GetOrdersByOtherAppResponse{
		Infos: resp.Infos,
	}, nil
}

func (s *Server) GetOrdersByGood(ctx context.Context, in *npool.GetOrdersByGoodRequest) (*npool.GetOrdersByGoodResponse, error) {
	resp, err := order.GetOrdersByGood(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get order detail by good error: %v", err)
		return &npool.GetOrdersByGoodResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) SubmitOrder(ctx context.Context, in *npool.SubmitOrderRequest) (*npool.SubmitOrderResponse, error) {
	resp, err := order.SubmitOrder(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("submit order error: %v", err)
		return &npool.SubmitOrderResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) CreateOrderPayment(ctx context.Context, in *npool.CreateOrderPaymentRequest) (*npool.CreateOrderPaymentResponse, error) {
	resp, err := order.CreateOrderPayment(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("create order payment error: %v", err)
		return &npool.CreateOrderPaymentResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}
