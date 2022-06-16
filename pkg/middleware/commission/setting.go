package commission

import (
	"context"
	"fmt"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/referral"
	rsetting "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/referral/setting"
	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"
)

func GetAmountSettings(ctx context.Context, appID, userID string) ([]*inspirepb.AppPurchaseAmountSetting, error) {
	settings, err := rsetting.GetAmountSettingsByAppUser(ctx, appID, userID)
	if err != nil {
		return nil, fmt.Errorf("fail get amount settings: %v", err)
	}

	invitees, err := referral.GetLayeredInvitees(ctx, appID, userID)
	if err != nil {
		return nil, fmt.Errorf("fail get layered invitees: %v", err)
	}

	for _, iv := range invitees {
		if iv.InviterID != userID {
			continue
		}

		ivsettings, err := rsetting.GetAmountSettingsByAppUser(ctx, appID, iv.InviteeID)
		if err != nil {
			return nil, fmt.Errorf("fail get amount settings: %v", err)
		}

		settings = append(settings, ivsettings...)
	}

	return settings, nil
}

func CreateAmountSetting(ctx context.Context, appID, userID, targetUserID string, setting *inspirepb.AppPurchaseAmountSetting) ([]*inspirepb.AppPurchaseAmountSetting, error) {
	iv, err := grpc2.GetRegistrationInvitationByAppInvitee(ctx, &inspirepb.GetRegistrationInvitationByAppInviteeRequest{
		AppID:     appID,
		InviteeID: targetUserID,
	})
	if err != nil || iv == nil {
		return nil, fmt.Errorf("fail get registration invitation: %v", err)
	}
	if iv.InviterID != userID {
		return nil, fmt.Errorf("permission denied")
	}

	// Create amount settings
	// Send email

	return nil, nil
}
