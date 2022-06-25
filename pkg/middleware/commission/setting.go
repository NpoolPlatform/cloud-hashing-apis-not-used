package commission

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/referral"
	rsetting "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/referral/setting"
	inspirecli "github.com/NpoolPlatform/cloud-hashing-inspire/pkg/client"
	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"
	thirdgwpb "github.com/NpoolPlatform/message/npool/thirdgateway"
	thirdgwcli "github.com/NpoolPlatform/third-gateway/pkg/client"
	thirdgwconst "github.com/NpoolPlatform/third-gateway/pkg/const"
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

func CreateAmountSetting(
	ctx context.Context,
	appID, userID, targetUserID, langID string,
	inviterName, inviteeName string,
	setting *inspirepb.AppPurchaseAmountSetting,
) ([]*inspirepb.AppPurchaseAmountSetting, error) {
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

	settings, err := grpc2.GetAppPurchaseAmountSettingsByAppUser(ctx, &inspirepb.GetAppPurchaseAmountSettingsByAppUserRequest{
		AppID:  appID,
		UserID: iv.InviterID,
	})
	if err != nil {
		return nil, fmt.Errorf("fail get amount settings: %v", err)
	}

	percent := uint32(0)
	for _, s := range settings {
		if s.End != 0 || s.GoodID != setting.GoodID {
			continue
		}
		percent = s.Percent
		break
	}

	if setting.Percent > percent {
		return nil, fmt.Errorf("overflow percent")
	}

	setting.AppID = appID
	setting.UserID = targetUserID

	_, err = inspirecli.CreateAmountSetting(ctx, setting)
	if err != nil {
		return nil, fmt.Errorf("fail create amount setting: %v", err)
	}

	settings, err = grpc2.GetAppPurchaseAmountSettingsByAppUser(ctx, &inspirepb.GetAppPurchaseAmountSettingsByAppUserRequest{
		AppID:  appID,
		UserID: targetUserID,
	})
	if err != nil {
		return nil, fmt.Errorf("fail get amount settings: %v", err)
	}

	rsetting.UpdateAmountSettingsCache(ctx, appID, targetUserID, settings)
	summaries, err := referral.GetGoodSummaries(ctx, appID, targetUserID)
	if err != nil {
		logger.Sugar().Warnf("fail get good summaries: %v", err)
	}
	if summaries != nil {
		for _, sum := range summaries {
			if sum.GoodID == setting.GoodID {
				sum.Percent = setting.Percent
				break
			}
		}
		referral.UpdateGoodSummariesCache(ctx, appID, userID, summaries)
	}

	err = thirdgwcli.NotifyEmail(ctx, &thirdgwpb.NotifyEmailRequest{
		AppID:        appID,
		UserID:       userID,
		ReceiverID:   targetUserID,
		LangID:       langID,
		SenderName:   inviterName,
		ReceiverName: inviteeName,
		UsedFor:      thirdgwconst.UsedForSetCommission,
	})
	if err != nil {
		logger.Sugar().Warnf("fail notify email: %v", err)
	}

	return settings, nil
}
