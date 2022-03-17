package referral

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	"golang.org/x/xerrors"
)

func GetMyInvitations(ctx context.Context, in *npool.GetMyInvitationsRequest) (*npool.GetMyInvitationsResponse, error) { //nolint
	AddWatcher(in.GetAppID(), in.GetUserID())

	invitations, userInfo, err := getFullInvitations(in.GetAppID(), in.GetUserID())
	if err != nil {
		logger.Sugar().Errorf("fail get invitations: %v", err)
		return nil, xerrors.Errorf("fail get invitations: %v", err)
	}
	return &npool.GetMyInvitationsResponse{
		MySelf: userInfo,
		Infos:  invitations,
	}, nil
}

func GetMyDirectInvitations(ctx context.Context, in *npool.GetMyDirectInvitationsRequest) (*npool.GetMyDirectInvitationsResponse, error) { //nolint
	AddWatcher(in.GetAppID(), in.GetUserID())

	invitations, userInfo, err := getDirectInvitations(in.GetAppID(), in.GetUserID())
	if err != nil {
		logger.Sugar().Errorf("fail get invitations: %v", err)
		return nil, xerrors.Errorf("fail get invitations: %v", err)
	}
	return &npool.GetMyDirectInvitationsResponse{
		MySelf: userInfo,
		Infos:  invitations,
	}, nil
}
