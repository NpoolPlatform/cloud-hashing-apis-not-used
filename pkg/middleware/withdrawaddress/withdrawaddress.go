package withdrawaddress

import (
	"context"
	"fmt"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"

	"github.com/NpoolPlatform/message/npool/appuser/mgr/v2/signmethod"
	"github.com/NpoolPlatform/message/npool/third/mgr/v1/usedfor"
	thirdmwcli "github.com/NpoolPlatform/third-middleware/pkg/client/verify"

	constant "github.com/NpoolPlatform/cloud-hashing-apis/pkg/const"
	account "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/account"
	review "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/review"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	billingconst "github.com/NpoolPlatform/cloud-hashing-billing/pkg/message/const"
	appusermgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v1"
	billingpb "github.com/NpoolPlatform/message/npool/cloud-hashing-billing"
	coininfopb "github.com/NpoolPlatform/message/npool/coininfo"
	reviewpb "github.com/NpoolPlatform/message/npool/review-service"
	sphinxproxypb "github.com/NpoolPlatform/message/npool/sphinxproxy"
	reviewconst "github.com/NpoolPlatform/review-service/pkg/const"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
)

func Set(ctx context.Context, in *npool.SetWithdrawAddressRequest) (*npool.SetWithdrawAddressResponse, error) {
	user, err := usermwcli.GetUser(ctx, in.GetAppID(), in.GetUserID())
	if err != nil {
		return nil, err
	}

	accountN := in.GetAccount()

	accountType := signmethod.SignMethodType(signmethod.SignMethodType_value[in.GetAccountType()])
	if accountType == signmethod.SignMethodType_Google {
		accountN = user.GetGoogleSecret()
	}

	err = thirdmwcli.VerifyCode(
		ctx,
		in.GetAppID(),
		accountN,
		in.GetVerificationCode(),
		signmethod.SignMethodType(signmethod.SignMethodType_value[in.GetAccountType()]),
		usedfor.UsedFor_SetWithdrawAddress,
	)
	if err != nil {
		return nil, fmt.Errorf("fail verify code: %v", err)
	}

	coin, err := grpc2.GetCoinInfo(ctx, &coininfopb.GetCoinInfoRequest{
		ID: in.GetCoinTypeID(),
	})
	if err != nil || coin == nil {
		return nil, fmt.Errorf("fail get coin info: %v", err)
	}

	_, err = grpc2.GetBalance(ctx, &sphinxproxypb.GetBalanceRequest{
		Name:    coin.Name,
		Address: in.GetAddress(),
	})
	if err != nil {
		return nil, fmt.Errorf("fail get wallet balance: %v", err)
	}

	_account, err := account.CreateUserCoinAccount(ctx, &npool.CreateUserCoinAccountRequest{
		Info: &billingpb.CoinAccountInfo{
			CoinTypeID:             in.GetCoinTypeID(),
			Address:                in.GetAddress(),
			PlatformHoldPrivateKey: false,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("fail create account: %v", err)
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
			Labels:     in.GetLabels(),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("fail create address: %v", err)
	}

	_review, err := grpc2.CreateReview(ctx, &reviewpb.CreateReviewRequest{
		Info: &reviewpb.Review{
			AppID:      in.GetAppID(),
			Domain:     billingconst.ServiceName,
			ObjectType: constant.ReviewObjectUserWithdrawAddress,
			ObjectID:   address.ID,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("fail create review: %v", err)
	}

	_review.State = reviewconst.StateApproved
	_, err = grpc2.UpdateReview(ctx, &reviewpb.UpdateReviewRequest{
		Info: _review,
	})
	if err != nil {
		return nil, fmt.Errorf("fail update review: %v", err)
	}

	return &npool.SetWithdrawAddressResponse{
		Info: &npool.WithdrawAddress{
			Address: address,
			Account: _account.Info,
			State:   reviewconst.StateApproved,
		},
	}, nil
}

func Delete(ctx context.Context, in *npool.DeleteWithdrawAddressRequest) (*npool.DeleteWithdrawAddressResponse, error) {
	reviewState, _, err := review.GetReviewState(ctx, &reviewpb.GetReviewsByAppDomainObjectTypeIDRequest{
		AppID:      in.GetAppID(),
		Domain:     billingconst.ServiceName,
		ObjectType: constant.ReviewObjectUserWithdrawAddress,
		ObjectID:   in.GetID(),
	})
	if err != nil {
		return nil, fmt.Errorf("fail get withdraw address review state: %v", err)
	}
	if reviewState == reviewconst.StateWait {
		return nil, fmt.Errorf("fail delete reviewing withdarw address")
	}

	info, err := grpc2.DeleteUserWithdraw(ctx, &billingpb.DeleteUserWithdrawRequest{
		ID: in.GetID(),
	})
	if err != nil {
		return nil, fmt.Errorf("fail delete withdraw address: %v", err)
	}

	return &npool.DeleteWithdrawAddressResponse{
		Info: info,
	}, nil
}

func GetByAppUser(ctx context.Context, in *npool.GetWithdrawAddressesByAppUserRequest) (*npool.GetWithdrawAddressesByAppUserResponse, error) {
	_, err := grpc2.GetAppUserByAppUser(ctx, &appusermgrpb.GetAppUserByAppUserRequest{
		AppID:  in.GetAppID(),
		UserID: in.GetUserID(),
	})
	if err != nil {
		return nil, fmt.Errorf("fail get app user: %v", err)
	}

	infos, err := grpc2.GetUserWithdrawsByAppUser(ctx, &billingpb.GetUserWithdrawsByAppUserRequest{
		AppID:  in.GetAppID(),
		UserID: in.GetUserID(),
	})
	if err != nil {
		return nil, fmt.Errorf("fail get addresses: %v", err)
	}

	addresses := []*npool.WithdrawAddress{}

	for _, info := range infos {
		_account, err := grpc2.GetBillingAccount(ctx, &billingpb.GetCoinAccountRequest{
			ID: info.AccountID,
		})
		if err != nil {
			return nil, fmt.Errorf("fail get account: %v", err)
		}

		reviewState, reviewMessage, err := review.GetReviewState(ctx, &reviewpb.GetReviewsByAppDomainObjectTypeIDRequest{
			AppID:      in.GetAppID(),
			Domain:     billingconst.ServiceName,
			ObjectType: constant.ReviewObjectUserWithdrawAddress,
			ObjectID:   info.ID,
		})
		if err != nil {
			return nil, fmt.Errorf("fail get review state: %v", err)
		}

		addresses = append(addresses, &npool.WithdrawAddress{
			Address: info,
			Account: _account,
			State:   reviewState,
			Message: reviewMessage,
		})
	}

	return &npool.GetWithdrawAddressesByAppUserResponse{
		Infos: addresses,
	}, nil
}

// func UpdateWithdrawUpdateAddressReview(ctx context.Context, in *npool.UpdateWithdrawAddressReviewRequest) (*npool.UpdateWithdrawAddressReviewResponse, error) {
//	reviewInfo := in.GetInfo()
//
//	userWithdraw, err := grpc2.GetUserWithdraw(ctx, &billingpb.GetUserWithdrawRequest{
//		ID: reviewInfo.GetObjectID(),
//	})
//	if err != nil {
//		return nil, fmt.Errorf("fail get user withdraw: %v", err)
//	}
//	if userWithdraw == nil {
//		return nil, fmt.Errorf("fail get user withdraw")
//	}
//
//	billingAccount, err := grpc2.GetBillingAccount(ctx, &billingpb.GetCoinAccountRequest{
//		ID: userWithdraw.GetAccountID(),
//	})
//	if err != nil {
//		return nil, fmt.Errorf("fail get billing account: %v", err)
//	}
//	if billingAccount == nil {
//		return nil, fmt.Errorf("fail get billing account")
//	}
//
//	user, err := grpc2.GetAppUserInfoByAppUser(ctx, &appusermgrpb.GetAppUserInfoByAppUserRequest{
//		AppID:  reviewInfo.GetAppID(),
//		UserID: userWithdraw.GetUserID(),
//	})
//	if err != nil {
//		return nil, fmt.Errorf("fail get app user: %v", err)
//	}
//	if user == nil {
//		return nil, fmt.Errorf("fail get app user")
//	}
//
//	reviewResp, err := grpc2.GetReview(ctx, &reviewpb.GetReviewRequest{
//		ID: reviewInfo.GetID(),
//	})
//	if err != nil {
//		return nil, fmt.Errorf("fail get review: %v", err)
//	}
//	if reviewResp == nil {
//		return nil, fmt.Errorf("fail get review")
//	}
//
//	reviewUpResp, err := grpc2.UpdateReview(ctx,&reviewpb.UpdateReviewRequest{Info: in.GetInfo()})
//	if err != nil {
//		return nil, err
//	}
//
//	if reviewInfo.GetState() == reviewconst.StateApproved {
//		_, err = grpc2.CreateNotification(ctx, &notificationpbpb.CreateNotificationRequest{
//			Info: &notificationpbpb.UserNotification{
//				AppID:  reviewInfo.GetAppID(),
//				UserID: userWithdraw.GetUserID(),
//			},
//			Message:  in.GetInfo().GetMessage(),
//			LangID:   in.GetLangID(),
//			UsedFor:  notificationconstant.UsedForWithdrawAddressReviewApprovedNotification,
//			UserName: user.GetExtra().GetUsername(),
//		})
//		if err != nil {
//			return nil, fmt.Errorf("fail create notification: %v", err)
//		}
//	}
//	if reviewInfo.GetState() == reviewconst.StateRejected {
//		_, err = grpc2.CreateNotification(ctx, &notificationpbpb.CreateNotificationRequest{
//			Info: &notificationpbpb.UserNotification{
//				AppID:  reviewInfo.GetAppID(),
//				UserID: userWithdraw.GetUserID(),
//			},
//			Message:  in.GetInfo().GetMessage(),
//			LangID:   in.GetLangID(),
//			UsedFor:  notificationconstant.UsedForWithdrawAddressReviewRejectedNotification,
//			UserName: user.GetExtra().GetUsername(),
//		})
//		if err != nil {
//			return nil, fmt.Errorf("fail create notification: %v", err)
//		}
//	}
//
//	return &npool.UpdateWithdrawAddressReviewResponse{
//		Info: &npool.WithdrawAddressReview{
//			Review:  reviewUpResp,
//			User:    user,
//			Address: userWithdraw,
//			Account: billingAccount,
//		},
//	}, nil
// }
