package user

import (
	"context"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	referral "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/referral"
	verifymw "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/verify"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	appusermgrconst "github.com/NpoolPlatform/appuser-manager/pkg/const"
	appusermgrpb "github.com/NpoolPlatform/message/npool/appusermgr"
	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"
	logingwpb "github.com/NpoolPlatform/message/npool/logingateway"
	thirdgwpb "github.com/NpoolPlatform/message/npool/thirdgateway"
	thirdgwconst "github.com/NpoolPlatform/third-gateway/pkg/const"

	"golang.org/x/xerrors"
)

func Signup(ctx context.Context, in *npool.SignupRequest) (*npool.SignupResponse, error) { //nolint
	appUser, err := grpc2.GetAppUserByAppAccount(ctx, &appusermgrpb.GetAppUserByAppAccountRequest{
		AppID:   in.GetAppID(),
		Account: in.GetAccount(),
	})
	if err != nil || appUser != nil {
		return nil, xerrors.Errorf("fail get app user by app account: %v", err)
	}

	invitationCode := in.GetInvitationCode()
	inviterID := ""

	app, err := grpc2.GetApp(ctx, &appusermgrpb.GetAppInfoRequest{
		ID: in.GetAppID(),
	})
	if err != nil || app == nil {
		return nil, xerrors.Errorf("fail get app: %v", err)
	}

	if app.Ctrl != nil && app.Ctrl.InvitationCodeMust {
		if invitationCode == "" {
			return nil, xerrors.Errorf("invitation code is must")
		}
	}

	if invitationCode != "" {
		code, err := grpc2.GetUserInvitationCodeByCode(ctx, &inspirepb.GetUserInvitationCodeByCodeRequest{
			Code: invitationCode,
		})
		if err != nil {
			return nil, xerrors.Errorf("fail get user invitation code: %v", err)
		}

		if code == nil {
			if app.Ctrl != nil && app.Ctrl.InvitationCodeMust {
				return nil, xerrors.Errorf("fail get invitation code")
			}
		} else {
			if code.AppID != in.GetAppID() {
				return nil, xerrors.Errorf("invalid invitation code for app")
			}
			inviterID = code.UserID
		}
	}

	err = verifymw.VerifyCode(
		ctx,
		in.GetAppID(),
		"",
		in.GetAccount(),
		in.GetAccountType(),
		in.GetVerificationCode(),
		thirdgwconst.UsedForSignup,
		false,
	)
	if err != nil {
		return nil, xerrors.Errorf("fail verify code: %v", err)
	}

	emailAddress := ""
	phoneNO := ""

	if in.GetAccountType() == appusermgrconst.SignupByMobile {
		phoneNO = in.GetAccount()
	} else if in.GetAccountType() == appusermgrconst.SignupByEmail {
		emailAddress = in.GetAccount()
	}

	appUser, err = grpc2.Signup(ctx, &appusermgrpb.CreateAppUserWithSecretRequest{
		User: &appusermgrpb.AppUser{
			AppID:        in.GetAppID(),
			EmailAddress: emailAddress,
			PhoneNO:      phoneNO,
		},
		Secret: &appusermgrpb.AppUserSecret{
			AppID:        in.GetAppID(),
			PasswordHash: in.GetPasswordHash(),
		},
	})
	if err != nil || appUser == nil {
		return nil, xerrors.Errorf("fail signup: %v", err)
	}

	if invitationCode != "" && inviterID != "" {
		_, err = grpc2.CreateRegistrationInvitation(ctx, &inspirepb.CreateRegistrationInvitationRequest{
			Info: &inspirepb.RegistrationInvitation{
				AppID:     in.GetAppID(),
				InviterID: inviterID,
				InviteeID: appUser.ID,
			},
		})
		if err != nil {
			return nil, xerrors.Errorf("fail create registration invitation: %v", err)
		}

		referral.GetInviteesRT(ctx, in.GetAppID(), inviterID)
	}

	return &npool.SignupResponse{
		Info: appUser,
	}, nil
}

func UpdatePasswordByAppUser(ctx context.Context, in *npool.UpdatePasswordByAppUserRequest, checkOldPassword bool) (*npool.UpdatePasswordByAppUserResponse, error) {
	var err error

	err = verifymw.VerifyCode(
		ctx,
		in.GetAppID(),
		in.GetUserID(),
		in.GetAccount(),
		in.GetAccountType(),
		in.GetVerificationCode(),
		thirdgwconst.UsedForUpdate,
		true,
	)
	if err != nil {
		return nil, xerrors.Errorf("fail verify code: %v", err)
	}

	secret, err := grpc2.GetAppUserSecretByAppUser(ctx, &appusermgrpb.GetAppUserSecretByAppUserRequest{
		AppID:  in.GetAppID(),
		UserID: in.GetUserID(),
	})
	if err != nil || secret == nil {
		return nil, xerrors.Errorf("fail get app user secret: %v", err)
	}

	if checkOldPassword {
		_, err = grpc2.VerifyAppUserByAppAccountPassword(ctx, &appusermgrpb.VerifyAppUserByAppAccountPasswordRequest{
			AppID:        in.GetAppID(),
			Account:      in.GetAccount(),
			PasswordHash: in.GetOldPasswordHash(),
		})
		if err != nil {
			return nil, xerrors.Errorf("fail verify username or password: %v", err)
		}
	}

	secret.PasswordHash = in.GetPasswordHash()
	secret.Salt = ""

	secret, err = grpc2.UpdateAppUserSecret(ctx, &appusermgrpb.UpdateAppUserSecretRequest{
		Info: secret,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail update app user secret: %v", err)
	}

	return &npool.UpdatePasswordByAppUserResponse{
		Info: secret,
	}, nil
}

func UpdatePassword(ctx context.Context, in *npool.UpdatePasswordRequest) (*npool.UpdatePasswordResponse, error) {
	appUser, err := grpc2.GetAppUserByAppAccount(ctx, &appusermgrpb.GetAppUserByAppAccountRequest{
		AppID:   in.GetAppID(),
		Account: in.GetAccount(),
	})
	if err != nil || appUser == nil {
		return nil, xerrors.Errorf("fail get app user by app account: %v", err)
	}
	resp, err := UpdatePasswordByAppUser(ctx, &npool.UpdatePasswordByAppUserRequest{
		AppID:            in.GetAppID(),
		UserID:           appUser.ID,
		Account:          in.GetAccount(),
		AccountType:      in.GetAccountType(),
		PasswordHash:     in.GetPasswordHash(),
		VerificationCode: in.GetVerificationCode(),
	}, false)
	if err != nil {
		return nil, xerrors.Errorf("fail update password: %v", err)
	}

	return &npool.UpdatePasswordResponse{
		Info: resp.Info,
	}, nil
}

func UpdateEmailAddress(ctx context.Context, in *npool.UpdateEmailAddressRequest) (*npool.UpdateEmailAddressResponse, error) { //nolint
	old, err := grpc2.GetAppUserByAppAccount(ctx, &appusermgrpb.GetAppUserByAppAccountRequest{
		AppID:   in.GetAppID(),
		Account: in.GetNewEmailAddress(),
	})
	if err != nil || old != nil {
		return nil, xerrors.Errorf("fail get app user by app account: %v", err)
	}

	info, err := grpc2.GetAppUserInfoByAppUser(ctx, &appusermgrpb.GetAppUserInfoByAppUserRequest{
		AppID:  in.GetAppID(),
		UserID: in.GetUserID(),
	})
	if err != nil || info == nil {
		return nil, xerrors.Errorf("fail get app user by app user: %v", err)
	}

	if in.GetOldAccountType() == appusermgrconst.SignupByMobile {
		if in.GetOldAccount() != info.User.PhoneNO {
			return nil, xerrors.Errorf("invalid account info")
		}
	} else if in.GetOldAccountType() == appusermgrconst.SignupByEmail {
		if in.GetOldAccount() != info.User.EmailAddress {
			return nil, xerrors.Errorf("invalid account info")
		}
	} else {
		return nil, xerrors.Errorf("invalid account type")
	}

	err = verifymw.VerifyCode(
		ctx,
		in.GetAppID(),
		in.GetUserID(),
		in.GetOldAccount(),
		in.GetOldAccountType(),
		in.GetOldVerificationCode(),
		thirdgwconst.UsedForUpdate,
		true,
	)
	if err != nil {
		return nil, xerrors.Errorf("fail verify code: %v", err)
	}

	_, err = grpc2.VerifyEmailCode(ctx, &thirdgwpb.VerifyEmailCodeRequest{
		AppID:        in.GetAppID(),
		EmailAddress: in.GetNewEmailAddress(),
		UsedFor:      thirdgwconst.UsedForUpdate,
		Code:         in.GetNewEmailVerificationCode(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail verify code: %v", err)
	}

	info.User.EmailAddress = in.GetNewEmailAddress()
	_, err = grpc2.UpdateAppUser(ctx, &appusermgrpb.UpdateAppUserRequest{
		Info: info.User,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail update app user: %v", err)
	}

	_, err = grpc2.UpdateCache(ctx, &logingwpb.UpdateCacheRequest{
		Info: info,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail update cache: %v", err)
	}

	return &npool.UpdateEmailAddressResponse{
		Info: info,
	}, nil
}

func UpdatePhoneNO(ctx context.Context, in *npool.UpdatePhoneNORequest) (*npool.UpdatePhoneNOResponse, error) { //nolint
	old, err := grpc2.GetAppUserByAppAccount(ctx, &appusermgrpb.GetAppUserByAppAccountRequest{
		AppID:   in.GetAppID(),
		Account: in.GetNewPhoneNO(),
	})
	if err != nil || old != nil {
		return nil, xerrors.Errorf("fail get app user by app account: %v", err)
	}

	info, err := grpc2.GetAppUserInfoByAppUser(ctx, &appusermgrpb.GetAppUserInfoByAppUserRequest{
		AppID:  in.GetAppID(),
		UserID: in.GetUserID(),
	})
	if err != nil || info == nil {
		return nil, xerrors.Errorf("fail get app user by app user: %v", err)
	}

	if in.GetOldAccountType() == appusermgrconst.SignupByMobile {
		if in.GetOldAccount() != info.User.PhoneNO {
			return nil, xerrors.Errorf("invalid account info")
		}
	} else if in.GetOldAccountType() == appusermgrconst.SignupByEmail {
		if in.GetOldAccount() != info.User.EmailAddress {
			return nil, xerrors.Errorf("invalid account info")
		}
	} else {
		return nil, xerrors.Errorf("invalid account type")
	}

	err = verifymw.VerifyCode(
		ctx,
		in.GetAppID(),
		in.GetUserID(),
		in.GetOldAccount(),
		in.GetOldAccountType(),
		in.GetOldVerificationCode(),
		thirdgwconst.UsedForUpdate,
		true,
	)
	if err != nil {
		return nil, xerrors.Errorf("fail verify code: %v", err)
	}

	resp1, err := grpc2.VerifySMSCode(ctx, &thirdgwpb.VerifySMSCodeRequest{
		AppID:   in.GetAppID(),
		PhoneNO: in.GetNewPhoneNO(),
		UsedFor: thirdgwconst.UsedForUpdate,
		Code:    in.GetNewPhoneVerificationCode(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail verify code: %v", err)
	}
	if resp1.Code < 0 {
		return nil, xerrors.Errorf("fail verify code")
	}

	info.User.PhoneNO = in.GetNewPhoneNO()
	_, err = grpc2.UpdateAppUser(ctx, &appusermgrpb.UpdateAppUserRequest{
		Info: info.User,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail update app user: %v", err)
	}

	_, err = grpc2.UpdateCache(ctx, &logingwpb.UpdateCacheRequest{
		Info: info,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail update cache: %v", err)
	}

	return &npool.UpdatePhoneNOResponse{
		Info: info,
	}, nil
}
