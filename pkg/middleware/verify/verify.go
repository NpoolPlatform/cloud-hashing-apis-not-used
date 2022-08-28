package verify

import (
	"context"
	"fmt"

	ga "github.com/NpoolPlatform/appuser-gateway/pkg/ga"
	appusermgrconst "github.com/NpoolPlatform/appuser-manager/pkg/const"
	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	appusermgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v1"
	thirdgwpb "github.com/NpoolPlatform/message/npool/thirdgateway"
)

func verifyByMobile(ctx context.Context, appID, phoneNO, code, usedFor string) (int32, error) {
	resp, err := grpc2.VerifySMSCode(ctx, &thirdgwpb.VerifySMSCodeRequest{
		AppID:   appID,
		PhoneNO: phoneNO,
		UsedFor: usedFor,
		Code:    code,
	})
	if err != nil {
		return -1, fmt.Errorf("fail verify sms code: %v", err)
	}

	return resp.Code, nil
}

func verifyByEmail(ctx context.Context, appID, emailAddr, code, usedFor string) (int32, error) {
	resp, err := grpc2.VerifyEmailCode(ctx, &thirdgwpb.VerifyEmailCodeRequest{
		AppID:        appID,
		EmailAddress: emailAddr,
		UsedFor:      usedFor,
		Code:         code,
	})
	if err != nil {
		return -1, fmt.Errorf("fail verify email code: %v", err)
	}

	return resp.Code, nil
}

func verifyByGoogle(ctx context.Context, appID, userID, code string) (int32, error) {
	if _, err := ga.VerifyGoogleAuth(ctx, appID, userID, code); err != nil {
		return -1, err
	}
	return 0, nil
}

func VerifyCode(ctx context.Context, appID, userID, account, accountType, code, usedFor string, accountMatch bool) error { //nolint
	var verified int32
	var err error

	if accountMatch {
		user, err := grpc2.GetAppUserByAppUser(ctx, &appusermgrpb.GetAppUserByAppUserRequest{
			AppID:  appID,
			UserID: userID,
		})
		if err != nil {
			return fmt.Errorf("fail get app user: %v", err)
		}

		switch accountType {
		case appusermgrconst.SignupByMobile:
			if user.PhoneNO != account {
				return fmt.Errorf("invalid mobile")
			}
		case appusermgrconst.SignupByEmail:
			if user.EmailAddress != account {
				return fmt.Errorf("invalid email")
			}
		}
	}

	switch accountType {
	case appusermgrconst.SignupByMobile:
		verified, err = verifyByMobile(ctx, appID, account, code, usedFor)
	case appusermgrconst.SignupByEmail:
		verified, err = verifyByEmail(ctx, appID, account, code, usedFor)
	default:
		verified, err = verifyByGoogle(ctx, appID, userID, code)
	}

	if err != nil {
		return fmt.Errorf("fail verify code: %v", err)
	}
	if verified < 0 {
		return fmt.Errorf("fail verify code")
	}

	return nil
}
