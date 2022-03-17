package referral

import (
	"context"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"
)

func getInvitees(ctx context.Context, appID, userID string) ([]*inspirepb.RegistrationInvitation, error) {
	return grpc2.GetRegistrationInvitationsByAppInviter(ctx, &inspirepb.GetRegistrationInvitationsByAppInviterRequest{
		AppID:     appID,
		InviterID: userID,
	})
}

func getInviter(ctx context.Context, appID, userID string) (*inspirepb.RegistrationInvitation, error) {
	return grpc2.GetRegistrationInvitationByAppInvitee(ctx, &inspirepb.GetRegistrationInvitationByAppInviteeRequest{
		AppID:     appID,
		InviteeID: userID,
	})
}
