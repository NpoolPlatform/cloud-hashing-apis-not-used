package commission

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	"github.com/NpoolPlatform/message/npool/third/mgr/v1/usedfor"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	"github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/referral"
	rsetting "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/referral/setting"
	inspirecli "github.com/NpoolPlatform/cloud-hashing-inspire/pkg/client"
	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"
	thirdmwcli "github.com/NpoolPlatform/third-middleware/pkg/client/notify"

	goodcli "github.com/NpoolPlatform/good-middleware/pkg/client/good"
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

// nolint
func CreateAmountSetting(
	ctx context.Context,
	appID, userID, targetUserID, langID string,
	inviterName, inviteeName string,
	setting *inspirepb.AppPurchaseAmountSetting,
) ([]*inspirepb.AppPurchaseAmountSetting, error) {
	inCode, err := grpc2.GetUserInvitationCodeByAppUser(ctx, &inspirepb.GetUserInvitationCodeByAppUserRequest{
		AppID:  appID,
		UserID: targetUserID,
	})
	if err != nil {
		return nil, err
	}
	if inCode == nil {
		return nil, fmt.Errorf("user is not KOL")
	}
	good, err := goodcli.GetGood(ctx, setting.GoodID)
	if err != nil {
		return nil, err
	}
	if good == nil {
		return nil, fmt.Errorf("good not exist")
	}
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
		if s == nil {
			continue
		}
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

	//rsetting.UpdateAmountSettingsCache(ctx, appID, targetUserID, settings)
	//summaries, err := referral.GetGoodSummaries(ctx, appID, targetUserID)
	//if err != nil {
	//	logger.Sugar().Warnf("fail get good summaries: %v", err)
	//}
	//if summaries != nil {
	//	for _, sum := range summaries {
	//		if sum.GoodID == setting.GoodID {
	//			sum.Percent = setting.Percent
	//			break
	//		}
	//	}
	//	referral.UpdateGoodSummariesCache(ctx, appID, userID, summaries)
	//}

	user, err := usermwcli.GetUser(ctx, appID, userID)
	if err != nil {
		return nil, err
	}

	targetUser, err := usermwcli.GetUser(ctx, appID, targetUserID)
	if err != nil {
		return nil, err
	}

	err = thirdmwcli.NotifyEmail(
		ctx,
		appID,
		user.GetEmailAddress(),
		usedfor.UsedFor_SetCommission,
		targetUser.GetEmailAddress(),
		langID,
		inviterName,
		inviteeName,
	)
	if err != nil {
		logger.Sugar().Warnf("fail notify email: %v", err)
	}

	return settings, nil
}
