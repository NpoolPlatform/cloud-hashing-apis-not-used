package verify

import (
	"context"

	appusermgrconst "github.com/NpoolPlatform/appuser-manager/pkg/const"
	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	appusermgrpb "github.com/NpoolPlatform/message/npool/appusermgr"
	thirdgwpb "github.com/NpoolPlatform/message/npool/thirdgateway"

	"golang.org/x/xerrors"
)

func verifyByMobile(ctx context.Context, appID, userID, phoneNO, code, usedFor string) (int32, error) {
	user, err := grpc2.GetAppUserByAppUser(ctx, &appusermgrpb.GetAppUserByAppUserRequest{
		AppID:  appID,
		UserID: userID,
	})
	if err != nil {
		return -1, xerrors.Errorf("fail get app user: %v", err)
	}

	if user.Info.PhoneNO != phoneNO {
		return -1, xerrors.Errorf("invalid mobile")
	}

	resp, err := grpc2.VerifySMSCode(ctx, &thirdgwpb.VerifySMSCodeRequest{
		AppID:   appID,
		PhoneNO: phoneNO,
		UsedFor: usedFor,
		Code:    code,
	})
	if err != nil {
		return -1, xerrors.Errorf("fail verify sms code: %v", err)
	}

	return resp.Code, nil
}

func verifyByEmail(ctx context.Context, appID, userID, emailAddr, code, usedFor string) (int32, error) {
	user, err := grpc2.GetAppUserByAppUser(ctx, &appusermgrpb.GetAppUserByAppUserRequest{
		AppID:  appID,
		UserID: userID,
	})
	if err != nil {
		return -1, xerrors.Errorf("fail get app user: %v", err)
	}

	if user.Info.EmailAddress != emailAddr {
		return -1, xerrors.Errorf("invalid email address")
	}

	resp, err := grpc2.VerifyEmailCode(ctx, &thirdgwpb.VerifyEmailCodeRequest{
		AppID:        appID,
		EmailAddress: emailAddr,
		UsedFor:      usedFor,
		Code:         code,
	})
	if err != nil {
		return -1, xerrors.Errorf("fail verify email code: %v", err)
	}

	return resp.Code, nil
}

func verifyByGoogle(ctx context.Context, appID, userID, code string) (int32, error) {
	_, err := grpc2.GetAppUserByAppUser(ctx, &appusermgrpb.GetAppUserByAppUserRequest{
		AppID:  appID,
		UserID: userID,
	})
	if err != nil {
		return -1, xerrors.Errorf("fail get app user: %v", err)
	}

	resp, err := grpc2.VerifyGoogleAuthentication(ctx, &thirdgwpb.VerifyGoogleAuthenticationRequest{
		AppID:  appID,
		UserID: userID,
		Code:   code,
	})
	if err != nil {
		return -1, xerrors.Errorf("fail verify google code: %v", err)
	}

	return resp.Code, nil
}

func VerifyCode(ctx context.Context, appID, userID, account, accountType, code, usedFor string) error { //nolint
	var verified int32
	var err error

	switch accountType {
	case appusermgrconst.SignupByMobile:
		verified, err = verifyByMobile(ctx, appID, userID, account, code, usedFor)
	case appusermgrconst.SignupByEmail:
		verified, err = verifyByEmail(ctx, appID, userID, account, code, usedFor)
	default:
		verified, err = verifyByGoogle(ctx, appID, userID, code)
	}

	if err != nil {
		return xerrors.Errorf("fail verify code: %v", err)
	}
	if verified < 0 {
		return xerrors.Errorf("fail verify code")
	}

	return nil
}
