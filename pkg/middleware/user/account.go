package user

import (
	"context"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	verifymw "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/verify"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	appusermgrconst "github.com/NpoolPlatform/appuser-manager/pkg/const"
	appusermgrpb "github.com/NpoolPlatform/message/npool/appusermgr"
	thirdgwconst "github.com/NpoolPlatform/third-gateway/pkg/const"

	"golang.org/x/xerrors"
)

const (
	EmailCodeVerified  = uint32(0x0001)
	MobileCodeVerified = uint32(0x0002)
	GoogleCodeVerified = uint32(0x0004)
)

func verifyCode(ctx context.Context, user *appusermgrpb.AppUserInfo, code *npool.VerificationCode) (uint32, error) {
	switch code.GetAccountType() {
	case appusermgrconst.SignupByMobile:
		if code.GetAccount() != user.User.PhoneNO {
			return 0, xerrors.Errorf("invalid phone NO")
		}
	case appusermgrconst.SignupByEmail:
		if code.GetAccount() != user.User.EmailAddress {
			return 0, xerrors.Errorf("invalid phone NO")
		}
	}

	err := verifymw.VerifyCode(
		ctx,
		user.User.AppID,
		user.User.ID,
		code.GetAccount(),
		code.GetAccountType(),
		code.GetVerificationCode(),
		thirdgwconst.UsedForUpdate,
		true,
	)
	if err != nil {
		return 0, xerrors.Errorf("fail verify code: %v", err)
	}

	switch code.GetAccountType() {
	case appusermgrconst.SignupByMobile:
		return MobileCodeVerified, nil
	case appusermgrconst.SignupByEmail:
		return EmailCodeVerified, nil
	}

	return GoogleCodeVerified, nil
}

func verifyCodes(ctx context.Context, user *appusermgrpb.AppUserInfo, codes []*npool.VerificationCode) (uint32, error) {
	verified := uint32(0)
	for _, code := range codes {
		v, err := verifyCode(ctx, user, code)
		if err != nil {
			return 0, xerrors.Errorf("fail verify code: %v", err)
		}
		verified |= v
	}
	return verified, nil
}

func codesVerified(accountType string, v uint32) bool {
	switch accountType {
	case appusermgrconst.SignupByMobile:
		fallthrough //nolint
	case appusermgrconst.SignupByEmail:
		if v&GoogleCodeVerified == 0 {
			return false
		}
		if v&MobileCodeVerified == 0 {
			return false
		}
		fallthrough //nolint
	default:
		if v&EmailCodeVerified == 0 {
			return false
		}
	}

	return true
}

func updateAccount(ctx context.Context, user *appusermgrpb.AppUserInfo, in *npool.UpdateAccountRequest) (*appusermgrpb.AppUserInfo, error) {
	switch in.GetAccountType() {
	case appusermgrconst.SignupByMobile:
		user.User.PhoneNO = in.GetAccount()
	case appusermgrconst.SignupByEmail:
		user.User.EmailAddress = in.GetAccount()
	default:
		return user, nil
	}

	_, err := grpc2.UpdateAppUser(ctx, &appusermgrpb.UpdateAppUserRequest{
		Info: user.User,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail update app user: %v", err)
	}

	return user, nil
}

func UpdateAccount(ctx context.Context, in *npool.UpdateAccountRequest) (*npool.UpdateAccountResponse, error) { //nolint
	old, err := grpc2.GetAppUserByAppAccount(ctx, &appusermgrpb.GetAppUserByAppAccountRequest{
		AppID:   in.GetAppID(),
		Account: in.GetAccount(),
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

	verified, err := verifyCodes(ctx, info, in.GetVerificationCodes())
	if err != nil {
		return nil, xerrors.Errorf("fail verify codes: %v", err)
	}

	if !codesVerified(in.GetAccountType(), verified) {
		return nil, xerrors.Errorf("miss required verification code")
	}

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

	info, err = updateAccount(ctx, info, in)
	if err != nil {
		return nil, xerrors.Errorf("fail update account: %v", err)
	}

	return &npool.UpdateAccountResponse{
		Info: info,
	}, nil
}
