package withdraw

import (
	"context"
	"fmt"
	"time"

	constant "github.com/NpoolPlatform/cloud-hashing-apis/pkg/const"
	currency "github.com/NpoolPlatform/cloud-hashing-staker/pkg/middleware/currency"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	billingconst "github.com/NpoolPlatform/cloud-hashing-billing/pkg/message/const"
	appusermgrpb "github.com/NpoolPlatform/message/npool/appusermgr"
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
		resp1, err := grpc2.CreateCoinAccountTransaction(ctx, &billingpb.CreateCoinAccountTransactionRequest{
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

		resp.Info.PlatformTransactionID = resp1.Info.ID
		_, err = grpc2.UpdateUserWithdrawItem(ctx, &billingpb.UpdateUserWithdrawItemRequest{
			Info: resp.Info,
		})
		if err != nil {
			return nil, xerrors.Errorf("fail update user withdraw item: %v", err)
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

// This API should not convert app id and user id in header to body and body.Info
func Update(ctx context.Context, in *npool.UpdateUserWithdrawReviewRequest) (*npool.UpdateUserWithdrawReviewResponse, error) { //nolint
	user, err := grpc2.GetAppUserByAppUser(ctx, &appusermgrpb.GetAppUserByAppUserRequest{
		AppID:  in.GetAppID(),
		UserID: in.GetUserID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get app user: %v", err)
	}
	if user.Info == nil {
		return nil, xerrors.Errorf("fail get app user")
	}

	resp, err := grpc2.GetReview(ctx, &reviewpb.GetReviewRequest{
		ID: in.GetReview().GetID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get review: %v", err)
	}
	if resp.Info == nil {
		return nil, xerrors.Errorf("fail get review")
	}

	resp1, err := grpc2.GetUserWithdrawItem(ctx, &billingpb.GetUserWithdrawItemRequest{
		ID: resp.Info.ObjectID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get object: %v", err)
	}
	if resp1.Info == nil {
		return nil, xerrors.Errorf("fail get object")
	}

	account, err := grpc2.GetBillingAccount(ctx, &billingpb.GetCoinAccountRequest{
		ID: resp1.Info.WithdrawToAccountID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get account info: %v", err)
	}
	if account.Info == nil {
		return nil, xerrors.Errorf("fail get account info")
	}

	withdrawAccount, err := grpc2.GetUserWithdrawByAccount(ctx, &billingpb.GetUserWithdrawByAccountRequest{
		AccountID: resp1.Info.WithdrawToAccountID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get withdraw account")
	}

	if withdrawAccount.Info.AppID != resp1.Info.AppID || withdrawAccount.Info.UserID != resp1.Info.UserID {
		return nil, xerrors.Errorf("invalid account")
	}

	_, err = grpc2.GetAppUserByAppUser(ctx, &appusermgrpb.GetAppUserByAppUserRequest{
		AppID:  resp1.Info.AppID,
		UserID: resp1.Info.UserID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get app user: %v", err)
	}
	if user.Info == nil {
		return nil, xerrors.Errorf("fail get app user")
	}

	if resp1.Info.AppID != in.GetReview().GetAppID() {
		return nil, xerrors.Errorf("invalid request")
	}

	if resp.Info.State != reviewconst.StateWait {
		return nil, xerrors.Errorf("already reviewed")
	}

	reviewState := resp.Info.State

	resp.Info.State = in.GetReview().GetState()
	_, err = grpc2.UpdateReview(ctx, &reviewpb.UpdateReviewRequest{
		Info: resp.Info,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail update review state: %v", err)
	}

	if in.GetReview().GetState() == reviewconst.StateApproved {
		coinsetting, err := grpc2.GetCoinSettingByCoin(ctx, &billingpb.GetCoinSettingByCoinRequest{
			CoinTypeID: resp1.Info.CoinTypeID,
		})
		if err != nil {
			return nil, xerrors.Errorf("fail get coin setting: %v", err)
		}

		account, err := grpc2.GetBillingAccount(ctx, &billingpb.GetCoinAccountRequest{
			ID: coinsetting.Info.UserOnlineAccountID,
		})
		if err != nil {
			return nil, xerrors.Errorf("fail get account info: %v", err)
		}
		if account.Info == nil {
			return nil, xerrors.Errorf("fail get account info")
		}

		coin, err := grpc2.GetCoinInfo(ctx, &coininfopb.GetCoinInfoRequest{
			ID: resp1.Info.CoinTypeID,
		})
		if err != nil {
			return nil, xerrors.Errorf("fail get coin info: %v", err)
		}
		if coin.Info == nil {
			return nil, xerrors.Errorf("fail get coin info")
		}

		balance, err := grpc2.GetBalance(ctx, &sphinxproxypb.GetBalanceRequest{
			Name:    coin.Info.Name,
			Address: account.Info.Address,
		})
		if err != nil {
			return nil, xerrors.Errorf("fail get wallet balance: %v", err)
		}

		if balance.Info.Balance < resp1.Info.Amount {
			return nil, xerrors.Errorf("insufficient funds")
		}

		resp2, err := grpc2.CreateCoinAccountTransaction(ctx, &billingpb.CreateCoinAccountTransactionRequest{
			Info: &billingpb.CoinAccountTransaction{
				AppID:              resp1.Info.AppID,
				UserID:             resp1.Info.UserID,
				FromAddressID:      coinsetting.Info.UserOnlineAccountID,
				ToAddressID:        resp1.Info.WithdrawToAccountID,
				CoinTypeID:         resp1.Info.CoinTypeID,
				Amount:             resp1.Info.Amount,
				Message:            fmt.Sprintf("user withdraw at %v", time.Now()),
				ChainTransactionID: uuid.New().String(),
			},
		})
		if err != nil {
			return nil, xerrors.Errorf("fail create coin account transaction: %v", err)
		}

		resp1.Info.PlatformTransactionID = resp2.Info.ID
		_, err = grpc2.UpdateUserWithdrawItem(ctx, &billingpb.UpdateUserWithdrawItemRequest{
			Info: resp1.Info,
		})
		if err != nil {
			return nil, xerrors.Errorf("fail update user withdraw item: %v", err)
		}

		resp.Info.State = reviewconst.StateApproved
		_, err = grpc2.UpdateReview(ctx, &reviewpb.UpdateReviewRequest{
			Info: resp.Info,
		})
		if err != nil {
			return nil, xerrors.Errorf("fail update review state: %v", err)
		}

		reviewState = reviewconst.StateApproved
	}

	return &npool.UpdateUserWithdrawReviewResponse{
		Info: &npool.UserWithdraw{
			Withdraw: resp1.Info,
			State:    reviewState,
		},
	}, nil
}

func GetByAppUser(ctx context.Context, in *npool.GetUserWithdrawsByAppUserRequest) (*npool.GetUserWithdrawsByAppUserResponse, error) {
	return nil, nil
}
