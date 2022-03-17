package referral

import (
	"context"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	cache "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/cache"
	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"

	"golang.org/x/xerrors"
)

func getInvitees(ctx context.Context, appID, userID string) ([]*inspirepb.RegistrationInvitation, error) {
	cacheFor := "invitees"

	invitees := cache.GetEntry(cacheKey(appID, userID, cacheFor))
	if invitees != nil {
		return invitees.([]*inspirepb.RegistrationInvitation), nil
	}

	invitees, err := grpc2.GetRegistrationInvitationsByAppInviter(ctx, &inspirepb.GetRegistrationInvitationsByAppInviterRequest{
		AppID:     appID,
		InviterID: userID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get invitations: %v", err)
	}

	cache.AddEntry(cacheKey(appID, userID, cacheFor), invitees)
	return invitees.([]*inspirepb.RegistrationInvitation), nil

}

func getInviter(ctx context.Context, appID, userID string) (*inspirepb.RegistrationInvitation, error) {
	cacheFor := "inviter"

	inviter := cache.GetEntry(cacheKey(appID, userID, cacheFor))
	if inviter != nil {
		return inviter.(*inspirepb.RegistrationInvitation), nil
	}

	inviter, err := grpc2.GetRegistrationInvitationByAppInvitee(ctx, &inspirepb.GetRegistrationInvitationByAppInviteeRequest{
		AppID:     appID,
		InviteeID: userID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get inviter: %v", err)
	}

	cache.AddEntry(cacheKey(appID, userID, cacheFor), inviter)
	return inviter.(*inspirepb.RegistrationInvitation), nil
}

func getNextLayerInvitees(ctx context.Context, curLayer []*inspirepb.RegistrationInvitation) ([]*inspirepb.RegistrationInvitation, error) {
	invitees := []*inspirepb.RegistrationInvitation{}

	for _, iv := range curLayer {
		ivs, err := getInvitees(ctx, iv.AppID, iv.InviteeID)
		if err != nil {
			return nil, xerrors.Errorf("fail get invitees: %v", err)
		}
		invitees = append(invitees, ivs...)
	}

	return invitees, nil
}

func getLayeredInvitees(ctx context.Context, appID, userID string) ([]*inspirepb.RegistrationInvitation, error) {
	invitees, err := getInvitees(ctx, appID, userID)
	if err != nil {
		return nil, xerrors.Errorf("fail get invitees: %v", err)
	}

	curLayer := invitees

	for {
		nextLayer, err := getNextLayerInvitees(ctx, curLayer)
		if err != nil {
			return nil, xerrors.Errorf("fail get invitees: %v", err)
		}

		if len(nextLayer) == 0 {
			break
		}

		invitees = append(invitees, nextLayer...)
		curLayer = nextLayer
	}

	return invitees, nil
}
