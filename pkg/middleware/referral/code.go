package referral

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	inspirecli "github.com/NpoolPlatform/cloud-hashing-inspire/pkg/client"
	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"
	thirdgwpb "github.com/NpoolPlatform/message/npool/thirdgateway"
	thirdgwcli "github.com/NpoolPlatform/third-gateway/pkg/client"
	thirdgwconst "github.com/NpoolPlatform/third-gateway/pkg/const"
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

	err = thirdgwcli.NotifyEmail(ctx, &thirdgwpb.NotifyEmailRequest{
		AppID:        appID,
		UserID:       userID,
		ReceiverID:   targetUserID,
		LangID:       langID,
		SenderName:   inviterName,
		ReceiverName: inviteeName,
		UsedFor:      thirdgwconst.UsedForCreateInvitationCode,
	})
	if err != nil {
		logger.Sugar().Warnf("fail notify email: %v", err)
	}

	return code, nil
}
