package referral

import (
	"context"

	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"

	"golang.org/x/xerrors"
)

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
