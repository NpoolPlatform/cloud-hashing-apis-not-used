package withdrawaddress

import (
	"context"

	constant "github.com/NpoolPlatform/cloud-hashing-apis/pkg/const"
	account "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/account"
	review "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/review"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	appusermgrconst "github.com/NpoolPlatform/appuser-manager/pkg/const"
	billingconst "github.com/NpoolPlatform/cloud-hashing-billing/pkg/message/const"
	appusermgrpb "github.com/NpoolPlatform/message/npool/appusermgr"
	billingpb "github.com/NpoolPlatform/message/npool/cloud-hashing-billing"
	coininfopb "github.com/NpoolPlatform/message/npool/coininfo"
	reviewpb "github.com/NpoolPlatform/message/npool/review-service"
	sphinxproxypb "github.com/NpoolPlatform/message/npool/sphinxproxy"
	thirdgwpb "github.com/NpoolPlatform/message/npool/thirdgateway"
	reviewconst "github.com/NpoolPlatform/review-service/pkg/const"
	thirdgwconst "github.com/NpoolPlatform/third-gateway/pkg/const"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"

	"golang.org/x/xerrors"
)

func Set(ctx context.Context, in *npool.SetWithdrawAddressRequest) (*npool.SetWithdrawAddressResponse, error) {
	user, err := grpc2.GetAppUserByAppUser(ctx, &appusermgrpb.GetAppUserByAppUserRequest{
		AppID:  in.GetAppID(),
		UserID: in.GetUserID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get app user: %v", err)
	}

	if in.GetAccountType() == appusermgrconst.SignupByMobile {
		phoneNO := in.GetAccount()

		if user.Info.PhoneNO != phoneNO {
			return nil, xerrors.Errorf("invalid mobile")
		}

		_, err = grpc2.VerifySMSCode(ctx, &thirdgwpb.VerifySMSCodeRequest{
			AppID:   in.GetAppID(),
			PhoneNO: phoneNO,
			UsedFor: thirdgwconst.UsedForSetWithdrawAddress,
			Code:    in.GetVerificationCode(),
		})
	} else if in.GetAccountType() == appusermgrconst.SignupByEmail {
		emailAddr := in.GetAccount()

		if user.Info.EmailAddress != emailAddr {
			return nil, xerrors.Errorf("invalid email address")
		}

		_, err = grpc2.VerifyEmailCode(ctx, &thirdgwpb.VerifyEmailCodeRequest{
			AppID:        in.GetAppID(),
			EmailAddress: emailAddr,
			UsedFor:      thirdgwconst.UsedForSetWithdrawAddress,
			Code:         in.GetVerificationCode(),
		})
	} else {
		_, err = grpc2.VerifyGoogleAuthentication(ctx, &thirdgwpb.VerifyGoogleAuthenticationRequest{
			AppID:  in.GetAppID(),
			UserID: in.GetUserID(),
			Code:   in.GetVerificationCode(),
		})
	}
	if err != nil {
		return nil, xerrors.Errorf("fail verify signup code: %v", err)
	}

	coin, err := grpc2.GetCoinInfo(ctx, &coininfopb.GetCoinInfoRequest{
		ID: in.GetCoinTypeID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get coin info: %v", err)
	}
	if coin.Info == nil {
		return nil, xerrors.Errorf("fail get coin info")
	}

	_, err = grpc2.GetBalance(ctx, &sphinxproxypb.GetBalanceRequest{
		Name:    coin.Info.Name,
		Address: in.GetAddress(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get wallet balance: %v", err)
	}

	_account, err := account.CreateUserCoinAccount(ctx, &npool.CreateUserCoinAccountRequest{
		Info: &billingpb.CoinAccountInfo{
			CoinTypeID:             in.GetCoinTypeID(),
			Address:                in.GetAddress(),
			PlatformHoldPrivateKey: false,
		},
	})
	if err != nil {
		return nil, xerrors.Errorf("fail create account: %v", err)
	}

	// TODO: rollback create user coin account
	address, err := grpc2.CreateUserWithdraw(ctx, &billingpb.CreateUserWithdrawRequest{
		Info: &billingpb.UserWithdraw{
			AppID:      in.GetAppID(),
			UserID:     in.GetUserID(),
			CoinTypeID: in.GetCoinTypeID(),
			AccountID:  _account.Info.ID,
			Name:       in.GetName(),
			Message:    in.GetMessage(),
		},
	})
	if err != nil {
		return nil, xerrors.Errorf("fail create address: %v", err)
	}

	_, err = grpc2.CreateReview(ctx, &reviewpb.CreateReviewRequest{
		Info: &reviewpb.Review{
			AppID:      in.GetAppID(),
			Domain:     billingconst.ServiceName,
			ObjectType: constant.ReviewObjectUserWithdrawAddress,
			ObjectID:   address.Info.ID,
		},
	})
	if err != nil {
		return nil, xerrors.Errorf("fail create review: %v", err)
	}

	return &npool.SetWithdrawAddressResponse{
		Info: &npool.WithdrawAddress{
			Address: address.Info,
			Account: _account.Info,
			State:   reviewconst.StateWait,
		},
	}, nil
}

func GetByAppUser(ctx context.Context, in *npool.GetWithdrawAddressesByAppUserRequest) (*npool.GetWithdrawAddressesByAppUserResponse, error) {
	_, err := grpc2.GetAppUserByAppUser(ctx, &appusermgrpb.GetAppUserByAppUserRequest{
		AppID:  in.GetAppID(),
		UserID: in.GetUserID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get app user: %v", err)
	}

	resp, err := grpc2.GetUserWithdrawsByAppUser(ctx, &billingpb.GetUserWithdrawsByAppUserRequest{
		AppID:  in.GetAppID(),
		UserID: in.GetUserID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get addresses: %v", err)
	}

	addresses := []*npool.WithdrawAddress{}

	for _, info := range resp.Infos {
		_account, err := grpc2.GetBillingAccount(ctx, &billingpb.GetCoinAccountRequest{
			ID: info.AccountID,
		})
		if err != nil {
			return nil, xerrors.Errorf("fail get account: %v", err)
		}

		reviewState, reviewMessage, err := review.GetReviewState(ctx, &reviewpb.GetReviewsByAppDomainObjectTypeIDRequest{
			AppID:      in.GetAppID(),
			Domain:     billingconst.ServiceName,
			ObjectType: constant.ReviewObjectUserWithdrawAddress,
			ObjectID:   info.ID,
		})
		if err != nil {
			return nil, xerrors.Errorf("fail get review state: %v", err)
		}

		addresses = append(addresses, &npool.WithdrawAddress{
			Address: info,
			Account: _account.Info,
			State:   reviewState,
			Message: reviewMessage,
		})
	}

	return &npool.GetWithdrawAddressesByAppUserResponse{
		Infos: addresses,
	}, nil
}
