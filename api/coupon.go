// +build !codeanalysis

package api

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/coupon"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetCouponsByAppUser(ctx context.Context, in *npool.GetCouponsByAppUserRequest) (*npool.GetCouponsByAppUserResponse, error) {
	resp, err := coupon.GetCouponsByAppUser(ctx, in)
	if err != nil {
		logger.Sugar().Errorf("get coupons error: %w", err)
		return &npool.GetCouponsByAppUserResponse{}, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}
