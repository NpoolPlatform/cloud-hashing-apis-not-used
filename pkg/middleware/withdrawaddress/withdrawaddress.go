package withdrawaddress

import (
	"context"

	constant "github.com/NpoolPlatform/cloud-hashing-apis/pkg/const"
	account "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/account"
	review "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/review"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	billingconst "github.com/NpoolPlatform/cloud-hashing-billing/pkg/message/const"
	appusermgrpb "github.com/NpoolPlatform/message/npool/appusermgr"
	billingpb "github.com/NpoolPlatform/message/npool/cloud-hashing-billing"
	reviewpb "github.com/NpoolPlatform/message/npool/review-service"
	reviewconst "github.com/NpoolPlatform/review-service/pkg/const"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"

	"golang.org/x/xerrors"
)

func Set(ctx context.Context, in *npool.SetWithdrawAddressRequest) (*npool.SetWithdrawAddressResponse, error) {
	_, err := grpc2.GetAppUserByAppUser(ctx, &appusermgrpb.GetAppUserByAppUserRequest{
		AppID:  in.GetAppID(),
		UserID: in.GetUserID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get app user: %v", err)
	}

	account, err := account.CreateUserCoinAccount(ctx, &npool.CreateUserCoinAccountRequest{
		Info: &billingpb.CoinAccountInfo{
			CoinTypeID:             in.GetCoinTypeID(),
			Address:                in.GetAddress(),
			PlatformHoldPrivateKey: false,
		},
	})
	if err != nil {
		return nil, xerrors.Errorf("fail create account: %v", err)
	}

	address, err := grpc2.CreateUserWithdraw(ctx, &billingpb.CreateUserWithdrawRequest{
		Info: &billingpb.UserWithdraw{
			AppID:      in.GetAppID(),
			UserID:     in.GetUserID(),
			CoinTypeID: in.GetCoinTypeID(),
			AccountID:  account.Info.ID,
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
			Account: account.Info,
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
		account, err := grpc2.GetBillingAccount(ctx, &billingpb.GetCoinAccountRequest{
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
			Account: account.Info,
			State:   reviewState,
			Message: reviewMessage,
		})
	}

	return &npool.GetWithdrawAddressesByAppUserResponse{
		Infos: addresses,
	}, nil
}
