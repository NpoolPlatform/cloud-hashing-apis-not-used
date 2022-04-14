//go:build !codeanalysis
// +build !codeanalysis

package api

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"go.opentelemetry.io/otel"

	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	constant "github.com/NpoolPlatform/cloud-hashing-apis/pkg/message/const"
	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/coupon"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetCouponsByAppUser(ctx context.Context, in *npool.GetCouponsByAppUserRequest) (*npool.GetCouponsByAppUserResponse, error) {
	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetCouponsByAppUser")
	defer span.End()

	resp, err := coupon.GetCouponsByAppUser(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get coupons error: %w", err)
		return &npool.GetCouponsByAppUserResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}
