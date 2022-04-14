//go:build !codeanalysis
// +build !codeanalysis

package api

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"go.opentelemetry.io/otel"

	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	constant "github.com/NpoolPlatform/cloud-hashing-apis/pkg/message/const"
	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/order" //nolint

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetOrder(ctx context.Context, in *npool.GetOrderRequest) (*npool.GetOrderResponse, error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetOrder")
	defer span.End()

	resp, err := order.GetOrder(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get order detail error: %v", err)
		return &npool.GetOrderResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) GetOrdersByAppUser(ctx context.Context, in *npool.GetOrdersByAppUserRequest) (*npool.GetOrdersByAppUserResponse, error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetOrdersByAppUser")
	defer span.End()

	resp, err := order.GetOrdersByAppUser(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get order detail by app user error: %v", err)
		return &npool.GetOrdersByAppUserResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) GetOrdersByApp(ctx context.Context, in *npool.GetOrdersByAppRequest) (*npool.GetOrdersByAppResponse, error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetOrdersByApp")
	defer span.End()

	resp, err := order.GetOrdersByApp(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get order detail by app error: %v", err)
		return &npool.GetOrdersByAppResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) GetOrdersByGood(ctx context.Context, in *npool.GetOrdersByGoodRequest) (*npool.GetOrdersByGoodResponse, error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetOrdersByGood")
	defer span.End()

	resp, err := order.GetOrdersByGood(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get order detail by good error: %v", err)
		return &npool.GetOrdersByGoodResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) SubmitOrder(ctx context.Context, in *npool.SubmitOrderRequest) (*npool.SubmitOrderResponse, error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "SubmitOrder")
	defer span.End()

	resp, err := order.SubmitOrder(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("submit order error: %v", err)
		return &npool.SubmitOrderResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *Server) CreateOrderPayment(ctx context.Context, in *npool.CreateOrderPaymentRequest) (*npool.CreateOrderPaymentResponse, error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateOrderPayment")
	defer span.End()

	resp, err := order.CreateOrderPayment(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("create order payment error: %v", err)
		return &npool.CreateOrderPaymentResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}
