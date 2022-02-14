package withdraw

import (
	"context"
	"fmt"
	"time"

	constant "github.com/NpoolPlatform/cloud-hashing-apis/pkg/const"
	currency "github.com/NpoolPlatform/cloud-hashing-staker/pkg/middleware/currency"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	billingconst "github.com/NpoolPlatform/cloud-hashing-billing/pkg/message/const"
	billingpb "github.com/NpoolPlatform/message/npool/cloud-hashing-billing"
	coininfopb "github.com/NpoolPlatform/message/npool/coininfo"
	reviewpb "github.com/NpoolPlatform/message/npool/review-service"
	sphinxproxypb "github.com/NpoolPlatform/message/npool/sphinxproxy"
	reviewconst "github.com/NpoolPlatform/review-service/pkg/const"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"

	"github.com/google/uuid"
	"golang.org/x/xerrors"
)

func Create(ctx context.Context, in *npool.SubmitUserWithdrawRequest) (*npool.SubmitUserWithdrawResponse, error) { //nolint
	if in.GetInfo().GetAmount() <= 0 {
		return nil, xerrors.Errorf("invalid amount")
	}

	coin, err := grpc2.GetCoinInfo(ctx, &coininfopb.GetCoinInfoRequest{
		ID: in.GetInfo().GetCoinTypeID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get coin info: %v", err)
	}
	if coin.Info == nil {
		return nil, xerrors.Errorf("fail get coin info")
	}

	account, err := grpc2.GetBillingAccount(ctx, &billingpb.GetCoinAccountRequest{
		ID: in.GetInfo().GetWithdrawToAccountID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get account info: %v", err)
	}
	if account.Info == nil {
		return nil, xerrors.Errorf("fail get account info")
	}

	withdrawAccount, err := grpc2.GetUserWithdrawByAccount(ctx, &billingpb.GetUserWithdrawByAccountRequest{
		AccountID: in.GetInfo().GetWithdrawToAccountID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get withdraw account")
	}

	if withdrawAccount.Info.AppID != in.GetInfo().GetAppID() || withdrawAccount.Info.UserID != in.GetInfo().GetUserID() {
		return nil, xerrors.Errorf("invalid account")
	}

	autoReviewCoinAmount := 0

	setting, err := grpc2.GetAppWithdrawSettingByAppCoin(ctx, &billingpb.GetAppWithdrawSettingByAppCoinRequest{
		AppID:      in.GetInfo().GetAppID(),
		CoinTypeID: in.GetInfo().GetCoinTypeID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get app withdraw setting: %v", err)
	}
	if setting.Info != nil {
		autoReviewCoinAmount = int(setting.Info.WithdrawAutoReviewCoinAmount)
	} else {
		setting, err := grpc2.GetPlatformSetting(ctx, &billingpb.GetPlatformSettingRequest{})
		if err != nil {
			return nil, xerrors.Errorf("fail get platform setting: %v", err)
		}
		price, err := currency.USDPrice(ctx, coin.Info.Name)
		if err != nil {
			return nil, xerrors.Errorf("fail get coin price: %v", err)
		}
		autoReviewCoinAmount = int(setting.Info.WithdrawAutoReviewUSDAmount / price)
	}

	coinsetting, err := grpc2.GetCoinSettingByCoin(ctx, &billingpb.GetCoinSettingByCoinRequest{
		CoinTypeID: in.GetInfo().GetCoinTypeID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get coin setting: %v", err)
	}

	account, err = grpc2.GetBillingAccount(ctx, &billingpb.GetCoinAccountRequest{
		ID: coinsetting.Info.UserOnlineAccountID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get account info: %v", err)
	}
	if account.Info == nil {
		return nil, xerrors.Errorf("fail get account info")
	}

	balance, err := grpc2.GetBalance(ctx, &sphinxproxypb.GetBalanceRequest{
		Name:    coin.Info.Name,
		Address: account.Info.Address,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get wallet balance: %v", err)
	}

	resp, err := grpc2.CreateUserWithdrawItem(ctx, &billingpb.CreateUserWithdrawItemRequest{
		Info: in.GetInfo(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail create user withdraw item: %v", err)
	}

	reason := "auto review"
	autoReview := true

	if balance.Info.Balance < in.GetInfo().GetAmount()+coin.Info.ReservedAmount {
		reason = "insufficient"
		autoReview = false
	} else if float64(autoReviewCoinAmount) < in.GetInfo().GetAmount() {
		reason = "large amount"
		autoReview = false
	}

	review, err := grpc2.CreateReview(ctx, &reviewpb.CreateReviewRequest{
		Info: &reviewpb.Review{
			AppID:      in.GetInfo().GetAppID(),
			Domain:     billingconst.ServiceName,
			ObjectType: constant.ReviewObjectWithdraw,
			ObjectID:   resp.Info.ID,
			Trigger:    reason,
		},
	})
	if err != nil {
		// TODO: rollback user withdraw item database
		return nil, xerrors.Errorf("fail create user withdraw item review: %v", err)
	}

	reviewState := reviewconst.StateWait

	if autoReview {
		_, err := grpc2.CreateCoinAccountTransaction(ctx, &billingpb.CreateCoinAccountTransactionRequest{
			Info: &billingpb.CoinAccountTransaction{
				AppID:              in.GetInfo().GetAppID(),
				UserID:             in.GetInfo().GetUserID(),
				FromAddressID:      coinsetting.Info.UserOnlineAccountID,
				ToAddressID:        in.GetInfo().GetWithdrawToAccountID(),
				CoinTypeID:         in.GetInfo().GetCoinTypeID(),
				Amount:             in.GetInfo().GetAmount(),
				Message:            fmt.Sprintf("user withdraw at %v", time.Now()),
				ChainTransactionID: uuid.New().String(),
			},
		})
		if err != nil {
			return nil, xerrors.Errorf("fail create coin account transaction: %v", err)
		}

		review.Info.State = reviewconst.StateApproved
		_, err = grpc2.UpdateReview(ctx, &reviewpb.UpdateReviewRequest{
			Info: review.Info,
		})
		if err != nil {
			return nil, xerrors.Errorf("fail update review state: %v", err)
		}

		reviewState = reviewconst.StateApproved
	}

	return &npool.SubmitUserWithdrawResponse{
		Info: &npool.UserWithdraw{
			Withdraw: resp.Info,
			State:    reviewState,
		},
	}, nil
}

func Update(ctx context.Context, in *npool.UpdateUserWithdrawReviewRequest) (*npool.UpdateUserWithdrawReviewResponse, error) {
	return nil, nil
}

func GetByAppUser(ctx context.Context, in *npool.GetUserWithdrawsByAppUserRequest) (*npool.GetUserWithdrawsByAppUserResponse, error) {
	return nil, nil
}
