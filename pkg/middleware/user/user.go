package user

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"

	appusermgrconst "github.com/NpoolPlatform/appuser-manager/pkg/const"
	appusermgrpb "github.com/NpoolPlatform/message/npool/appusermgr"
	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"
	thirdgwpb "github.com/NpoolPlatform/message/npool/thirdgateway"
	thirdgwconst "github.com/NpoolPlatform/third-gateway/pkg/const"

	"golang.org/x/xerrors"
)

func Signup(ctx context.Context, in *npool.SignupRequest) (*npool.SignupResponse, error) { //nolint
	resp, err := grpc2.GetAppUserByAppAccount(ctx, &appusermgrpb.GetAppUserByAppAccountRequest{
		AppID:   in.GetAppID(),
		Account: in.GetAccount(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get app user by app account: %v", err)
	}
	if resp.Info != nil {
		return nil, xerrors.Errorf("fail check app user")
	}

	invitationCode := in.GetInvitationCode()
	inviterID := ""

	appResp, err := grpc2.GetApp(ctx, &appusermgrpb.GetAppInfoRequest{
		ID: in.GetAppID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get app: %v", err)
	}
	if appResp.Info == nil {
		return nil, xerrors.Errorf("fail get app")
	}

	if appResp.Info.Ctrl != nil && appResp.Info.Ctrl.InvitationCodeMust {
		if invitationCode == "" {
			return nil, xerrors.Errorf("invitation code is must")
		}
	}

	if invitationCode != "" {
		getByCodeResp, err := grpc2.GetUserInvitationCodeByCode(ctx, &inspirepb.GetUserInvitationCodeByCodeRequest{
			Code: invitationCode,
		})
		if err != nil {
			return nil, xerrors.Errorf("fail get user invitation code: %v", err)
		}

		if getByCodeResp.Info == nil {
			if appResp.Info.Ctrl != nil && appResp.Info.Ctrl.InvitationCodeMust {
				return nil, xerrors.Errorf("fail get invitation code")
			}
		} else {
			if getByCodeResp.Info.AppID != in.GetAppID() {
				return nil, xerrors.Errorf("invalid invitation code for app")
			}
			inviterID = getByCodeResp.Info.UserID
		}
	}

	emailAddr := ""
	phoneNO := ""

	if in.GetAccountType() == appusermgrconst.SignupByMobile {
		phoneNO = in.GetAccount()
		_, err = grpc2.VerifySMSCode(ctx, &thirdgwpb.VerifySMSCodeRequest{
			AppID:   in.GetAppID(),
			PhoneNO: phoneNO,
			UsedFor: thirdgwconst.UsedForSignup,
			Code:    in.GetVerificationCode(),
		})
	} else if in.GetAccountType() == appusermgrconst.SignupByEmail {
		emailAddr = in.GetAccount()
		_, err = grpc2.VerifyEmailCode(ctx, &thirdgwpb.VerifyEmailCodeRequest{
			AppID:        in.GetAppID(),
			EmailAddress: emailAddr,
			UsedFor:      thirdgwconst.UsedForSignup,
			Code:         in.GetVerificationCode(),
		})
	} else {
		return nil, xerrors.Errorf("invalid signup method")
	}
	if err != nil {
		return nil, xerrors.Errorf("fail verify signup code: %v", err)
	}

	signupResp, err := grpc2.Signup(ctx, &appusermgrpb.CreateAppUserWithSecretRequest{
		User: &appusermgrpb.AppUser{
			AppID:        in.GetAppID(),
			EmailAddress: emailAddr,
			PhoneNO:      phoneNO,
		},
		Secret: &appusermgrpb.AppUserSecret{
			AppID:        in.GetAppID(),
			PasswordHash: in.GetPasswordHash(),
		},
	})
	if err != nil {
		return nil, xerrors.Errorf("fail signup: %v", err)
	}

	if invitationCode != "" && inviterID != "" {
		_, err = grpc2.CreateRegistrationInvitation(ctx, &inspirepb.CreateRegistrationInvitationRequest{
			Info: &inspirepb.RegistrationInvitation{
				AppID:     in.GetAppID(),
				InviterID: inviterID,
				InviteeID: signupResp.Info.ID,
			},
		})
		if err != nil {
			return nil, xerrors.Errorf("fail create registration invitation: %v", err)
		}
	}

	return &npool.SignupResponse{
		Info: signupResp.Info,
	}, nil
}

func GetMyInvitations(ctx context.Context, in *npool.GetMyInvitationsRequest) (*npool.GetMyInvitationsResponse, error) { //nolint
	addWatcher(in.GetAppID(), in.GetUserID())

	invitations, userInfo, err := getFullInvitations(in.GetAppID(), in.GetUserID())
	if err != nil {
		return nil, xerrors.Errorf("fail get invitations: %v", err)
	}
	return &npool.GetMyInvitationsResponse{
		MySelf: userInfo,
		Infos:  invitations,
	}, nil
}

func GetMyDirectInvitations(ctx context.Context, in *npool.GetMyDirectInvitationsRequest) (*npool.GetMyDirectInvitationsResponse, error) { //nolint
	addWatcher(in.GetAppID(), in.GetUserID())

	invitations, userInfo, err := getDirectInvitations(in.GetAppID(), in.GetUserID())
	if err != nil {
		return nil, xerrors.Errorf("fail get invitations: %v", err)
	}
	return &npool.GetMyDirectInvitationsResponse{
		MySelf: userInfo,
		Infos:  invitations,
	}, nil
}

func UpdatePasswordByAppUser(ctx context.Context, in *npool.UpdatePasswordByAppUserRequest, checkOldPassword bool) (*npool.UpdatePasswordByAppUserResponse, error) {
	var err error
	emailAddr := ""
	phoneNO := ""

	if in.GetAccountType() == appusermgrconst.SignupByMobile {
		phoneNO = in.GetAccount()
		_, err = grpc2.VerifySMSCode(ctx, &thirdgwpb.VerifySMSCodeRequest{
			AppID:   in.GetAppID(),
			PhoneNO: phoneNO,
			UsedFor: thirdgwconst.UsedForSignup,
			Code:    in.GetVerificationCode(),
		})
	} else if in.GetAccountType() == appusermgrconst.SignupByEmail {
		emailAddr = in.GetAccount()
		_, err = grpc2.VerifyEmailCode(ctx, &thirdgwpb.VerifyEmailCodeRequest{
			AppID:        in.GetAppID(),
			EmailAddress: emailAddr,
			UsedFor:      thirdgwconst.UsedForSignup,
			Code:         in.GetVerificationCode(),
		})
	} else {
		return nil, xerrors.Errorf("invalid signup method")
	}
	if err != nil {
		return nil, xerrors.Errorf("fail verify signup code: %v", err)
	}

	resp, err := grpc2.GetAppUserSecretByAppUser(ctx, &appusermgrpb.GetAppUserSecretByAppUserRequest{
		AppID:  in.GetAppID(),
		UserID: in.GetUserID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get app user secret: %v", err)
	}
	if resp.Info == nil {
		return nil, xerrors.Errorf("fail get app user secret")
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

	resp.Info.PasswordHash = in.GetPasswordHash()
	resp.Info.Salt = ""

	resp1, err := grpc2.UpdateAppUserSecret(ctx, &appusermgrpb.UpdateAppUserSecretRequest{
		Info: resp.Info,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail update app user secret: %v", err)
	}

	return &npool.UpdatePasswordByAppUserResponse{
		Info: resp1.Info,
	}, nil
}

func UpdatePassword(ctx context.Context, in *npool.UpdatePasswordRequest) (*npool.UpdatePasswordResponse, error) {
	resp, err := grpc2.GetAppUserByAppAccount(ctx, &appusermgrpb.GetAppUserByAppAccountRequest{
		AppID:   in.GetAppID(),
		Account: in.GetAccount(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get app user by app account: %v", err)
	}
	if resp.Info == nil {
		return nil, xerrors.Errorf("fail get app user by app account")
	}

	resp1, err := UpdatePasswordByAppUser(ctx, &npool.UpdatePasswordByAppUserRequest{
		AppID:            in.GetAppID(),
		UserID:           resp.Info.ID,
		Account:          in.GetAccount(),
		AccountType:      in.GetAccountType(),
		PasswordHash:     in.GetPasswordHash(),
		VerificationCode: in.GetVerificationCode(),
	}, false)
	if err != nil {
		return nil, xerrors.Errorf("fail update password: %v", err)
	}

	return &npool.UpdatePasswordResponse{
		Info: resp1.Info,
	}, nil
}
