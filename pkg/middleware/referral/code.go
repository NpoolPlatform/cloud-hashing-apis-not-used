package referral

import (
	"context"
	"fmt"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
	"github.com/NpoolPlatform/message/npool/third/mgr/v1/usedfor"
	thirdmwcli "github.com/NpoolPlatform/third-middleware/pkg/client/notify"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	inspirecli "github.com/NpoolPlatform/cloud-hashing-inspire/pkg/client"
	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"
)

func CreateInvitationCode(
	ctx context.Context,
	appID, userID, targetUserID, langID string,
	inviterName, inviteeName string,
	setting *inspirepb.UserInvitationCode,
) (*inspirepb.UserInvitationCode, error) {
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

	code, err := grpc2.GetUserInvitationCodeByAppUser(ctx, &inspirepb.GetUserInvitationCodeByAppUserRequest{
		AppID:  appID,
		UserID: targetUserID,
	})
	if err != nil {
		return nil, fmt.Errorf("fail get invitation code: %v", err)
	}

	if code == nil {
		code, err = inspirecli.CreateInvitationCode(ctx, &inspirepb.UserInvitationCode{
			AppID:  appID,
			UserID: targetUserID,
		})
		if err != nil {
			return nil, fmt.Errorf("fail create invitation code: %v", err)
		}
	}

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
		usedfor.UsedFor_CreateInvitationCode,
		targetUser.GetEmailAddress(),
		langID,
		inviterName,
		inviteeName,
	)
	if err != nil {
		logger.Sugar().Warnf("fail notify email: %v", err)
	}

	return code, nil
}
