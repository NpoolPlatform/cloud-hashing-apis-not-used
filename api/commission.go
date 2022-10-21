//go:build !codeanalysis
// +build !codeanalysis

package api

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	commission "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/commission"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetAmountSettings(ctx context.Context, in *npool.GetAmountSettingsRequest) (*npool.GetAmountSettingsResponse, error) {
	settings, err := commission.GetAmountSettings(ctx, in.GetAppID(), in.GetUserID())
	if err != nil {
		logger.Sugar().Errorf("get amount settings error: %v", err)
		return &npool.GetAmountSettingsResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &npool.GetAmountSettingsResponse{
		Infos: settings,
	}, nil
}

func (s *Server) CreateAmountSetting(ctx context.Context, in *npool.CreateAmountSettingRequest) (*npool.CreateAmountSettingResponse, error) {
	settings, err := commission.CreateAmountSetting(
		ctx,
		in.GetAppID(), in.GetUserID(), in.GetTargetUserID(), in.GetLangID(),
		in.GetInviterName(), in.GetInviteeName(),
		in.GetInfo())
	if err != nil {
		logger.Sugar().Errorf("create amount settings error: %v", err)
		return &npool.CreateAmountSettingResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &npool.CreateAmountSettingResponse{
		Infos: settings,
	}, nil
}
