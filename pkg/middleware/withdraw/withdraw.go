package withdraw

import (
	"context"
	"fmt"
	"time"

	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"

	"github.com/NpoolPlatform/message/npool/appuser/mgr/v2/signmethod"
	"github.com/NpoolPlatform/message/npool/third/mgr/v1/usedfor"
	thirdmwcli "github.com/NpoolPlatform/third-middleware/pkg/client/verify"

	notificationpbpb "github.com/NpoolPlatform/message/npool/notification"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	notificationconstant "github.com/NpoolPlatform/notification/pkg/const"

	constant "github.com/NpoolPlatform/cloud-hashing-apis/pkg/const"
	commissionmw "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/commission"
	fee "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/fee"
	commissionsettingmw "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/referral/setting"
	review "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/review"
	redis2 "github.com/NpoolPlatform/go-service-framework/pkg/redis"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"
	currency "github.com/NpoolPlatform/oracle-manager/pkg/middleware/currency"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	billingcli "github.com/NpoolPlatform/cloud-hashing-billing/pkg/client"
	billingstate "github.com/NpoolPlatform/cloud-hashing-billing/pkg/const"
	billingconst "github.com/NpoolPlatform/cloud-hashing-billing/pkg/message/const"
	appusermgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v1"
	billingpb "github.com/NpoolPlatform/message/npool/cloud-hashing-billing"
	coininfopb "github.com/NpoolPlatform/message/npool/coininfo"
	reviewpb "github.com/NpoolPlatform/message/npool/review-service"
	sphinxproxypb "github.com/NpoolPlatform/message/npool/sphinxproxy"
	reviewconst "github.com/NpoolPlatform/review-service/pkg/const"

	"github.com/google/uuid"
)

// nolint
func Outcoming(ctx context.Context, appID, userID, coinTypeID, withdrawType string, includeReviewing bool) (float64, error) {
	withdraws := []*billingpb.UserWithdrawItem{}
	var err error

	switch withdrawType {
	case billingstate.WithdrawTypeBenefit:
		withdraws, err = grpc2.GetUserWithdrawItemsByAppUserCoinWithdrawType(ctx, &billingpb.GetUserWithdrawItemsByAppUserCoinWithdrawTypeRequest{
			AppID:        appID,
			UserID:       userID,
			CoinTypeID:   coinTypeID,
			WithdrawType: withdrawType,
		})
		if err != nil {
			return 0, fmt.Errorf("fail get user withdraws: %v", err)
		}
	case billingstate.WithdrawTypeCommission:
		fallthrough //nolint
	case billingstate.WithdrawTypeUserPaymentBalance:
		withdraws, err = grpc2.GetUserWithdrawItemsByAppUserWithdrawType(ctx, &billingpb.GetUserWithdrawItemsByAppUserWithdrawTypeRequest{
			AppID:        appID,
			UserID:       userID,
			WithdrawType: withdrawType,
		})
		if err != nil {
			return 0, fmt.Errorf("fail get user withdraws: %v", err)
		}
	}

	txs, err := grpc2.GetCoinAccountTransactionsByAppUserCoin(ctx, &billingpb.GetCoinAccountTransactionsByAppUserCoinRequest{
		AppID:      appID,
		UserID:     userID,
		CoinTypeID: coinTypeID,
	})
	if err != nil {
		return 0, fmt.Errorf("fail get account transactions: %v", err)
	}

	outcoming := 0.0

	for _, info := range txs {
		myWithdraw := false

		for _, withdraw := range withdraws {
			if withdraw.PlatformTransactionID == info.ID {
				myWithdraw = true
				break
			}
		}

		if !myWithdraw {
			continue
		}

		if info.State == billingstate.CoinTransactionStateFail ||
			info.State == billingstate.CoinTransactionStateRejected {
			continue
		}

		outcoming += info.Amount
	}

	if includeReviewing {
		states, err := GetByAppUser(ctx, &npool.GetUserWithdrawsByAppUserRequest{
			AppID:  appID,
			UserID: userID,
		})
		if err != nil {
			return 0, fmt.Errorf("fail get user withdraws: %v", err)
		}

		for _, s := range states.Infos {
			if s.Withdraw.WithdrawType != withdrawType {
				continue
			}
			if s.Withdraw.CoinTypeID != coinTypeID {
				continue
			}
			if s.State == reviewconst.StateWait {
				outcoming += s.Withdraw.Amount
			}
		}
	}

	return outcoming, nil
}

func CommissionCoinTypeID(ctx context.Context) (string, error) {
	coin, err := commissionsettingmw.GetUsingCoin(ctx)
	if err != nil {
		return "", fmt.Errorf("fail get using coin: %v", err)
	}

	return coin.CoinTypeID, nil
}

func userPaymentBalanceWithdrawable(ctx context.Context, appID, userID, withdrawType string, amount float64, includeReviewing bool, exemptFee bool) (bool, float64, float64, error) { //nolint
	balances, err := billingcli.GetUserPaymentBalances(ctx, appID, userID)
	if err != nil {
		return false, 0, 0, fmt.Errorf("fail get user payment balances: %v", err)
	}

	paymentBalance := 0.0
	invalidUUID := uuid.UUID{}.String()

	for _, balance := range balances {
		if balance.UsedByPaymentID != "" && balance.UsedByPaymentID != invalidUUID {
			continue
		}
		paymentBalance += balance.Amount * balance.CoinUSDCurrency
	}

	coinTypeID, err := CommissionCoinTypeID(ctx)
	if err != nil {
		return false, 0, 0, fmt.Errorf("fail get coin type id: %v", err)
	}

	outcoming, err := Outcoming(ctx, appID, userID, coinTypeID, withdrawType, includeReviewing)
	if err != nil {
		return false, 0, 0, fmt.Errorf("fail get withdraw outcoming: %v", err)
	}

	feeAmount := 0.0
	if !exemptFee {
		feeAmount, err = fee.Amount(ctx, coinTypeID)
		if err != nil {
			return false, 0, 0, fmt.Errorf("fail get fee amount: %v", err)
		}
	}

	if amount <= feeAmount {
		return false, 0, 0, fmt.Errorf("transfer payment balance amount is not enough for fee %v <= %v | %v", amount, feeAmount, exemptFee)
	}

	if paymentBalance-outcoming < amount {
		return false, 0, 0, fmt.Errorf("not sufficient funds (%v - %v < %v)", paymentBalance, outcoming, amount)
	}

	return true, amount, 0, nil
}

func commissionWithdrawable(ctx context.Context, appID, userID, withdrawType string, amount float64, includeReviewing bool) (bool, float64, float64, error) { //nolint
	myCommission, err := commissionmw.GetCommission(ctx, appID, userID)
	if err != nil {
		return false, 0, 0, fmt.Errorf("fail get total amount: %v", err)
	}

	coinTypeID, err := CommissionCoinTypeID(ctx)
	if err != nil {
		return false, 0, 0, fmt.Errorf("fail get coin type id: %v", err)
	}

	outcoming, err := Outcoming(ctx, appID, userID, coinTypeID, withdrawType, includeReviewing)
	if err != nil {
		return false, 0, 0, fmt.Errorf("fail get withdraw outcoming: %v", err)
	}

	if myCommission <= outcoming {
		return false, 0, 0, fmt.Errorf("invalid billing input")
	}

	feeAmount, err := fee.Amount(ctx, coinTypeID)
	if err != nil {
		return false, 0, 0, fmt.Errorf("fail get fee amount: %v", err)
	}

	if amount <= feeAmount {
		return false, 0, 0, fmt.Errorf("transfer commission amount is not enough for fee %v <= %v", amount, feeAmount)
	}

	amount1 := amount
	amount2 := 0.0

	if myCommission-outcoming < amount {
		amount1 = myCommission - outcoming
		amount2 = amount - amount1

		able, _, _, err := userPaymentBalanceWithdrawable(ctx, appID, userID, billingstate.WithdrawTypeUserPaymentBalance, amount-(myCommission-outcoming), includeReviewing, true)
		if err != nil {
			return false, 0, 0, fmt.Errorf("fail check user payment balance (%v - %v < %v): %v", myCommission, outcoming, amount, err)
		}
		if !able {
			return false, 0, 0, fmt.Errorf("not sufficient funds")
		}
	}

	return true, amount1, amount2, nil
}

func benefitWithdrawable(ctx context.Context, appID, userID, coinTypeID, withdrawType string, amount float64, includeReviewing bool) (bool, float64, float64, error) { //nolint
	benefits, err := grpc2.GetUserBenefitsByAppUserCoin(ctx, &billingpb.GetUserBenefitsByAppUserCoinRequest{
		AppID:      appID,
		UserID:     userID,
		CoinTypeID: coinTypeID,
	})
	if err != nil {
		return false, 0, 0, fmt.Errorf("fail get user benefits: %v", err)
	}

	incoming := 0.0
	for _, info := range benefits {
		incoming += info.Amount
	}

	outcoming, err := Outcoming(ctx, appID, userID, coinTypeID, withdrawType, includeReviewing)
	if err != nil {
		return false, 0, 0, fmt.Errorf("fail get withdraw outcoming: %v", err)
	}

	if incoming < outcoming {
		return false, 0, 0, fmt.Errorf("invalid billing input")
	}

	feeAmount, err := fee.Amount(ctx, coinTypeID)
	if err != nil {
		return false, 0, 0, fmt.Errorf("fail get fee amount: %v", err)
	}

	if amount <= feeAmount {
		return false, 0, 0, fmt.Errorf("transfer benefit amount is not enough for fee %v < %v", amount, feeAmount)
	}

	if incoming-outcoming < amount {
		return false, 0, 0, fmt.Errorf("not sufficient funds %v - %v < %v", incoming, outcoming, amount)
	}

	return true, amount, 0, nil
}

func withdrawable(ctx context.Context, appID, userID, coinTypeID, withdrawType string, amount float64, includeReviewing bool, exemptFee bool) (bool, float64, float64, error) { //nolint
	switch withdrawType {
	case billingstate.WithdrawTypeBenefit:
		return benefitWithdrawable(ctx, appID, userID, coinTypeID, withdrawType, amount, includeReviewing)
	case billingstate.WithdrawTypeCommission:
		return commissionWithdrawable(ctx, appID, userID, withdrawType, amount, includeReviewing)
	case billingstate.WithdrawTypeUserPaymentBalance:
		return userPaymentBalanceWithdrawable(ctx, appID, userID, withdrawType, amount, includeReviewing, exemptFee)
	}
	return false, 0, 0, fmt.Errorf("invalid withdraw type")
}

func Create(ctx context.Context, in *npool.SubmitUserWithdrawRequest) (*npool.SubmitUserWithdrawResponse, error) { //nolint
	if in.GetInfo().GetAmount() <= 0 {
		return nil, fmt.Errorf("invalid amount")
	}

	user, err := usermwcli.GetUser(ctx, in.GetInfo().GetAppID(), in.GetInfo().GetUserID())
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
		in.GetInfo().GetAppID(),
		accountN,
		in.GetVerificationCode(),
		signmethod.SignMethodType(signmethod.SignMethodType_value[in.GetAccountType()]),
		usedfor.UsedFor_Withdraw,
	)
	if err != nil {
		return nil, fmt.Errorf("fail verify code: %v", err)
	}

	coin, err := grpc2.GetCoinInfo(ctx, &coininfopb.GetCoinInfoRequest{
		ID: in.GetInfo().GetCoinTypeID(),
	})
	if err != nil || coin == nil {
		return nil, fmt.Errorf("fail get coin info: %v", err)
	}
	if coin.PreSale {
		return nil, fmt.Errorf("cannot withdraw presale coin")
	}

	lockKey := fmt.Sprintf("withdraw:%v:%v", in.GetInfo().GetAppID(), in.GetInfo().GetUserID())
	err = redis2.TryLock(lockKey, 10*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("lock withdraw fail: %v", err)
	}
	defer func() {
		if err := redis2.Unlock(lockKey); err != nil {
			logger.Sugar().Errorf("unlock withdraw fail: %v", err)
		}
	}()

	coinTypeID := in.GetInfo().GetCoinTypeID()
	if in.GetInfo().GetWithdrawType() == billingstate.WithdrawTypeCommission {
		coinTypeID, err = CommissionCoinTypeID(ctx)
		if err != nil {
			return nil, fmt.Errorf("fail get coin type id: %v", err)
		}
	}

	ok, amount1, amount2, err := withdrawable(
		ctx,
		in.GetInfo().GetAppID(),
		in.GetInfo().GetUserID(),
		coinTypeID,
		in.GetInfo().GetWithdrawType(),
		in.GetInfo().GetAmount(),
		true,
		false,
	)
	if !ok || err != nil {
		return nil, fmt.Errorf("user not withdrawable: %v", err)
	}

	account, err := grpc2.GetBillingAccount(ctx, &billingpb.GetCoinAccountRequest{
		ID: in.GetInfo().GetWithdrawToAccountID(),
	})
	if err != nil || account == nil {
		return nil, fmt.Errorf("fail get account info: %v", err)
	}

	withdrawAccount, err := grpc2.GetUserWithdrawByAccount(ctx, &billingpb.GetUserWithdrawByAccountRequest{
		AccountID: in.GetInfo().GetWithdrawToAccountID(),
	})
	if err != nil || withdrawAccount == nil {
		return nil, fmt.Errorf("fail get withdraw account")
	}

	if withdrawAccount.AppID != in.GetInfo().GetAppID() || withdrawAccount.UserID != in.GetInfo().GetUserID() {
		return nil, fmt.Errorf("acount is not belong to user")
	}

	reviewState, _, err := review.GetReviewState(ctx, &reviewpb.GetReviewsByAppDomainObjectTypeIDRequest{
		AppID:      in.GetInfo().GetAppID(),
		Domain:     billingconst.ServiceName,
		ObjectType: constant.ReviewObjectUserWithdrawAddress,
		ObjectID:   withdrawAccount.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("fail get review: %v", err)
	}
	if reviewState != reviewconst.StateApproved {
		return nil, fmt.Errorf("invalid account: not approved")
	}

	autoReviewCoinAmount := 0

	setting, err := grpc2.GetAppWithdrawSettingByAppCoin(ctx, &billingpb.GetAppWithdrawSettingByAppCoinRequest{
		AppID:      in.GetInfo().GetAppID(),
		CoinTypeID: coinTypeID,
	})
	if err != nil {
		return nil, fmt.Errorf("fail get app withdraw setting: %v", err)
	}
	if setting != nil {
		autoReviewCoinAmount = int(setting.WithdrawAutoReviewCoinAmount)
	} else {
		setting, err := grpc2.GetPlatformSetting(ctx, &billingpb.GetPlatformSettingRequest{})
		if err != nil || setting == nil {
			return nil, fmt.Errorf("fail get platform setting: %v", err)
		}
		price, err := currency.USDPrice(ctx, coin.Name)
		if err != nil {
			return nil, fmt.Errorf("fail get coin price: %v", err)
		}
		autoReviewCoinAmount = int(setting.WithdrawAutoReviewUSDAmount / price)
	}

	coinsetting, err := grpc2.GetCoinSettingByCoin(ctx, &billingpb.GetCoinSettingByCoinRequest{
		CoinTypeID: coinTypeID,
	})
	if err != nil {
		return nil, fmt.Errorf("fail get coin setting: %v", err)
	}

	account, err = grpc2.GetBillingAccount(ctx, &billingpb.GetCoinAccountRequest{
		ID: coinsetting.UserOnlineAccountID,
	})
	if err != nil || account == nil {
		return nil, fmt.Errorf("fail get account info: %v", err)
	}

	balance, err := grpc2.GetBalance(ctx, &sphinxproxypb.GetBalanceRequest{
		Name:    coin.Name,
		Address: account.Address,
	})
	if err != nil {
		return nil, fmt.Errorf("fail get wallet balance: %v", err)
	}

	reason := "auto review"
	autoReview := true

	if balance.Balance < in.GetInfo().GetAmount()+coin.ReservedAmount {
		reason = "insufficient"
		autoReview = false
	} else if float64(autoReviewCoinAmount) < in.GetInfo().GetAmount() {
		reason = "large amount"
		autoReview = false
	}

	reviewLockKey := fmt.Sprintf("withdraw-review:%v:%v", in.GetInfo().GetAppID(), in.GetInfo().GetUserID())
	err = redis2.TryLock(reviewLockKey, 0)
	if err != nil {
		return nil, fmt.Errorf("fail lock withdraw review: %v", err)
	}

	wInfo := in.GetInfo()

	var _review1 *reviewpb.Review
	var withdrawItem1 *billingpb.UserWithdrawItem
	var _review2 *reviewpb.Review
	var withdrawItem2 *billingpb.UserWithdrawItem

	if amount1 > 0 {
		wInfo.Amount = amount1

		withdrawItem1, err = grpc2.CreateUserWithdrawItem(ctx, &billingpb.CreateUserWithdrawItemRequest{
			Info: wInfo,
		})
		if err != nil {
			return nil, fmt.Errorf("fail create user withdraw item: %v", err)
		}

		_review1, err = grpc2.CreateReview(ctx, &reviewpb.CreateReviewRequest{
			Info: &reviewpb.Review{
				AppID:      in.GetInfo().GetAppID(),
				Domain:     billingconst.ServiceName,
				ObjectType: constant.ReviewObjectWithdraw,
				ObjectID:   withdrawItem1.ID,
				Trigger:    reason,
			},
		})
		if err != nil {
			// TODO: rollback user withdraw item database
			return nil, fmt.Errorf("fail create user withdraw item review: %v", err)
		}
	}

	if amount2 > 0 && wInfo.GetWithdrawType() == billingstate.WithdrawTypeCommission {
		wInfo.Amount = amount2
		wInfo.WithdrawType = billingstate.WithdrawTypeUserPaymentBalance
		wInfo.ExemptFee = true

		withdrawItem2, err = grpc2.CreateUserWithdrawItem(ctx, &billingpb.CreateUserWithdrawItemRequest{
			Info: wInfo,
		})
		if err != nil {
			return nil, fmt.Errorf("fail create user withdraw item: %v", err)
		}

		_review2, err = grpc2.CreateReview(ctx, &reviewpb.CreateReviewRequest{
			Info: &reviewpb.Review{
				AppID:      in.GetInfo().GetAppID(),
				Domain:     billingconst.ServiceName,
				ObjectType: constant.ReviewObjectWithdraw,
				ObjectID:   withdrawItem2.ID,
				Trigger:    reason,
			},
		})
		if err != nil {
			// TODO: rollback user withdraw item database
			return nil, fmt.Errorf("fail create user withdraw item review: %v", err)
		}
	}

	// TODO: check 24 hours total withdraw amount: if overflow, goto review
	reviewState = reviewconst.StateWait

	if autoReview {
		if err := redis2.Unlock(reviewLockKey); err != nil {
			logger.Sugar().Errorf("unlock withdraw review fail: %v", err)
		}

		feeAmount, err := fee.Amount(ctx, coinTypeID)
		if err != nil {
			return nil, fmt.Errorf("fail get fee amount: %v", err)
		}

		tx, err := grpc2.CreateCoinAccountTransaction(ctx, &billingpb.CreateCoinAccountTransactionRequest{
			Info: &billingpb.CoinAccountTransaction{
				AppID:              in.GetInfo().GetAppID(),
				UserID:             in.GetInfo().GetUserID(),
				GoodID:             uuid.UUID{}.String(),
				FromAddressID:      coinsetting.UserOnlineAccountID,
				ToAddressID:        in.GetInfo().GetWithdrawToAccountID(),
				CoinTypeID:         coinTypeID,
				Amount:             in.GetInfo().GetAmount(),
				TransactionFee:     feeAmount,
				Message:            fmt.Sprintf("user withdraw at %v", time.Now()),
				ChainTransactionID: "",
			},
		})
		if err != nil {
			return nil, fmt.Errorf("fail create coin account transaction: %v", err)
		}

		if withdrawItem1 != nil {
			withdrawItem1.PlatformTransactionID = tx.ID
			_, err = grpc2.UpdateUserWithdrawItem(ctx, &billingpb.UpdateUserWithdrawItemRequest{
				Info: withdrawItem1,
			})
			if err != nil {
				return nil, fmt.Errorf("fail update user withdraw item: %v", err)
			}
		}

		if withdrawItem2 != nil {
			withdrawItem2.PlatformTransactionID = tx.ID
			_, err = grpc2.UpdateUserWithdrawItem(ctx, &billingpb.UpdateUserWithdrawItemRequest{
				Info: withdrawItem2,
			})
			if err != nil {
				return nil, fmt.Errorf("fail update user withdraw item: %v", err)
			}
		}

		if _review1 != nil {
			_review1.State = reviewconst.StateApproved
			_, err = grpc2.UpdateReview(ctx, &reviewpb.UpdateReviewRequest{
				Info: _review1,
			})
			if err != nil {
				return nil, fmt.Errorf("fail update review state: %v", err)
			}
		}

		if _review2 != nil {
			_review2.State = reviewconst.StateApproved
			_, err = grpc2.UpdateReview(ctx, &reviewpb.UpdateReviewRequest{
				Info: _review2,
			})
			if err != nil {
				return nil, fmt.Errorf("fail update review state: %v", err)
			}
		}

		reviewState = reviewconst.StateApproved
	}

	withdrawItem := withdrawItem1
	if withdrawItem == nil {
		withdrawItem = withdrawItem2
	}

	return &npool.SubmitUserWithdrawResponse{
		Info: &npool.UserWithdraw{
			Withdraw: withdrawItem,
			State:    reviewState,
		},
	}, nil
}

func Update(ctx context.Context, in *npool.UpdateUserWithdrawReviewRequest) (*npool.UpdateUserWithdrawReviewResponse, error) { //nolint
	// TODO: check permission of reviewer

	if in.GetUserID() != in.GetInfo().GetReviewerID() {
		return nil, fmt.Errorf("mismatch reviewer id")
	}

	_review, err := grpc2.GetReview(ctx, &reviewpb.GetReviewRequest{
		ID: in.GetInfo().GetID(),
	})
	if err != nil || _review == nil {
		return nil, fmt.Errorf("fail get review: %v", err)
	}
	if _review.State == reviewconst.StateApproved {
		return nil, fmt.Errorf("already approved")
	}

	withdrawItem, err := grpc2.GetUserWithdrawItem(ctx, &billingpb.GetUserWithdrawItemRequest{
		ID: _review.ObjectID,
	})
	if err != nil || withdrawItem == nil {
		return nil, fmt.Errorf("fail get object: %v", err)
	}

	user, err := grpc2.GetAppUserInfoByAppUser(ctx, &appusermgrpb.GetAppUserInfoByAppUserRequest{
		AppID:  withdrawItem.AppID,
		UserID: withdrawItem.UserID,
	})
	if err != nil || user == nil {
		return nil, fmt.Errorf("fail get app user: %v", err)
	}

	invalidID := uuid.UUID{}.String()
	if withdrawItem.PlatformTransactionID != invalidID {
		return nil, fmt.Errorf("withdraw already processed")
	}

	coinsetting, err := grpc2.GetCoinSettingByCoin(ctx, &billingpb.GetCoinSettingByCoinRequest{
		CoinTypeID: withdrawItem.CoinTypeID,
	})
	if err != nil {
		return nil, fmt.Errorf("fail get coin setting: %v", err)
	}

	onlineAccount, err := grpc2.GetBillingAccount(ctx, &billingpb.GetCoinAccountRequest{
		ID: coinsetting.UserOnlineAccountID,
	})
	if err != nil || onlineAccount == nil {
		return nil, fmt.Errorf("fail get account info: %v", err)
	}

	ok, _, _, err := withdrawable(
		ctx,
		withdrawItem.AppID,
		withdrawItem.UserID,
		withdrawItem.CoinTypeID,
		withdrawItem.WithdrawType,
		withdrawItem.Amount,
		false,
		withdrawItem.ExemptFee,
	)
	if !ok || err != nil {
		return nil, fmt.Errorf("user not withdrawable: %v", err)
	}

	coin, err := grpc2.GetCoinInfo(ctx, &coininfopb.GetCoinInfoRequest{
		ID: withdrawItem.CoinTypeID,
	})
	if err != nil || coin == nil {
		return nil, fmt.Errorf("fail get coin info: %v", err)
	}

	account, err := grpc2.GetBillingAccount(ctx, &billingpb.GetCoinAccountRequest{
		ID: withdrawItem.WithdrawToAccountID,
	})
	if err != nil || account == nil {
		return nil, fmt.Errorf("fail get account info: %v", err)
	}

	withdrawAccount, err := grpc2.GetUserWithdrawByAccount(ctx, &billingpb.GetUserWithdrawByAccountRequest{
		AccountID: withdrawItem.WithdrawToAccountID,
	})
	if err != nil {
		return nil, fmt.Errorf("fail get withdraw account: %v", err)
	}

	var reviewState string

	if withdrawAccount == nil {
		if in.GetInfo().GetState() != reviewconst.StateRejected {
			return nil, fmt.Errorf("withdraw account may be deleted: you can only and have to reject this withdraw")
		}
		goto lUpdateReview
	}

	if withdrawAccount.AppID != withdrawItem.AppID || withdrawAccount.UserID != withdrawItem.UserID {
		return nil, fmt.Errorf("invalid account: mismatch app user")
	}

	reviewState, _, err = review.GetReviewState(ctx, &reviewpb.GetReviewsByAppDomainObjectTypeIDRequest{
		AppID:      withdrawItem.AppID,
		Domain:     billingconst.ServiceName,
		ObjectType: constant.ReviewObjectUserWithdrawAddress,
		ObjectID:   withdrawAccount.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("fail get review: %v", err)
	}
	if reviewState != reviewconst.StateApproved {
		return nil, fmt.Errorf("invalid account: invalid review state")
	}

	_, err = grpc2.GetAppUserByAppUser(ctx, &appusermgrpb.GetAppUserByAppUserRequest{
		AppID:  withdrawItem.AppID,
		UserID: withdrawItem.UserID,
	})
	if err != nil || user == nil {
		return nil, fmt.Errorf("fail get app user: %v", err)
	}

	if _review.State != reviewconst.StateWait {
		return nil, fmt.Errorf("already reviewed")
	}

	if in.GetInfo().GetState() == reviewconst.StateApproved {
		// TODO: here should hold transfer lock
		balance, err := grpc2.GetBalance(ctx, &sphinxproxypb.GetBalanceRequest{
			Name:    coin.Name,
			Address: onlineAccount.Address,
		})
		if err != nil {
			return nil, fmt.Errorf("fail get wallet balance: %v", err)
		}

		if balance.Balance < withdrawItem.Amount+coin.ReservedAmount {
			return nil, fmt.Errorf("not sufficient funds %v < %v + %v of coin %v address %v",
				balance.Balance, withdrawItem.Amount, coin.ReservedAmount, coin.Name, onlineAccount.Address)
		}

		account, err := grpc2.GetBillingAccount(ctx, &billingpb.GetCoinAccountRequest{
			ID: coinsetting.UserOnlineAccountID,
		})
		if err != nil || account == nil {
			return nil, fmt.Errorf("fail get account info: %v", err)
		}

		feeAmount := 0.0

		if !withdrawItem.ExemptFee {
			feeAmount, err = fee.Amount(ctx, withdrawItem.CoinTypeID)
			if err != nil {
				return nil, fmt.Errorf("fail get fee amount: %v", err)
			}
		}

		tx, err := grpc2.CreateCoinAccountTransaction(ctx, &billingpb.CreateCoinAccountTransactionRequest{
			Info: &billingpb.CoinAccountTransaction{
				AppID:              withdrawItem.AppID,
				UserID:             withdrawItem.UserID,
				GoodID:             uuid.UUID{}.String(),
				FromAddressID:      coinsetting.UserOnlineAccountID,
				ToAddressID:        withdrawItem.WithdrawToAccountID,
				CoinTypeID:         withdrawItem.CoinTypeID,
				Amount:             withdrawItem.Amount,
				TransactionFee:     feeAmount,
				Message:            fmt.Sprintf("user withdraw at %v", time.Now()),
				ChainTransactionID: "",
			},
		})
		if err != nil {
			return nil, fmt.Errorf("fail create coin account transaction: %v", err)
		}

		withdrawItem.PlatformTransactionID = tx.ID
		_, err = grpc2.UpdateUserWithdrawItem(ctx, &billingpb.UpdateUserWithdrawItemRequest{
			Info: withdrawItem,
		})
		if err != nil {
			return nil, fmt.Errorf("fail update user withdraw item: %v", err)
		}
	}

lUpdateReview:
	_review.State = in.GetInfo().GetState()
	_review.ReviewerID = in.GetInfo().GetReviewerID()

	_review, err = grpc2.UpdateReview(ctx, &reviewpb.UpdateReviewRequest{
		Info: _review,
	})
	if err != nil {
		return nil, fmt.Errorf("fail update review state: %v", err)
	}

	reviewLockKey := fmt.Sprintf("withdraw-review:%v:%v", withdrawItem.AppID, withdrawItem.UserID)
	err = redis2.Unlock(reviewLockKey)
	if err != nil {
		return nil, fmt.Errorf("fail unlock withdraw review: %v", err)
	}

	return &npool.UpdateUserWithdrawReviewResponse{
		Info: &npool.WithdrawReview{
			Withdraw: withdrawItem,
			Review:   _review,
			User:     user,
		},
	}, nil
}

func GetByAppUser(ctx context.Context, in *npool.GetUserWithdrawsByAppUserRequest) (*npool.GetUserWithdrawsByAppUserResponse, error) { //nolint
	withdrawItems, err := grpc2.GetUserWithdrawItemsByAppUser(ctx, &billingpb.GetUserWithdrawItemsByAppUserRequest{
		AppID:  in.GetAppID(),
		UserID: in.GetUserID(),
	})
	if err != nil {
		return nil, fmt.Errorf("fail get user withdraw items: %v", err)
	}

	withdraws := []*npool.UserWithdraw{}
	for _, info := range withdrawItems {
		reviewState, reviewMessage, err := review.GetReviewState(ctx, &reviewpb.GetReviewsByAppDomainObjectTypeIDRequest{
			AppID:      info.AppID,
			Domain:     billingconst.ServiceName,
			ObjectType: constant.ReviewObjectWithdraw,
			ObjectID:   info.ID,
		})
		if err != nil {
			return nil, fmt.Errorf("fail get review state: %v", err)
		}

		withdraws = append(withdraws, &npool.UserWithdraw{
			Withdraw: info,
			State:    reviewState,
			Message:  reviewMessage,
		})
	}

	return &npool.GetUserWithdrawsByAppUserResponse{
		Infos: withdraws,
	}, nil
}

func UpdateWithdrawReview(ctx context.Context, in *npool.UpdateWithdrawReviewRequest) (*npool.UpdateWithdrawReviewResponse, error) {
	reviewInfo := in.GetInfo()

	withdrawItem, err := grpc2.GetUserWithdrawItem(ctx, &billingpb.GetUserWithdrawItemRequest{
		ID: reviewInfo.GetObjectID(),
	})
	if err != nil {
		return nil, fmt.Errorf("fail get withdrawItem: %v", err)
	}
	if withdrawItem == nil {
		return nil, fmt.Errorf("fail get withdrawItem")
	}
	user, err := grpc2.GetAppUserInfoByAppUser(ctx, &appusermgrpb.GetAppUserInfoByAppUserRequest{
		AppID:  withdrawItem.GetAppID(),
		UserID: withdrawItem.GetUserID(),
	})
	if err != nil {
		return nil, fmt.Errorf("fail get app user: %v", err)
	}
	if user == nil {
		return nil, fmt.Errorf("fail get app user")
	}

	reviewResp, err := grpc2.GetReview(ctx, &reviewpb.GetReviewRequest{
		ID: reviewInfo.GetID(),
	})
	if err != nil {
		return nil, fmt.Errorf("fail get review: %v", err)
	}
	if reviewResp == nil {
		return nil, fmt.Errorf("fail get review")
	}

	reviewUpResp, err := grpc2.UpdateReview(ctx, &reviewpb.UpdateReviewRequest{Info: in.GetInfo()})
	if err != nil {
		return nil, err
	}

	if reviewInfo.GetState() == reviewconst.StateApproved && false {
		_, err = grpc2.CreateNotification(ctx, &notificationpbpb.CreateNotificationRequest{
			Info: &notificationpbpb.UserNotification{
				AppID:  withdrawItem.GetAppID(),
				UserID: withdrawItem.GetUserID(),
			},
			Message:  in.GetInfo().GetMessage(),
			LangID:   in.GetLangID(),
			UsedFor:  notificationconstant.UsedForWithdrawReviewApprovedNotification,
			UserName: user.GetExtra().GetUsername(),
		})
		if err != nil {
			return nil, fmt.Errorf("fail create notification: %v", err)
		}
	}

	if reviewInfo.GetState() == reviewconst.StateRejected && false {
		_, err = grpc2.CreateNotification(ctx, &notificationpbpb.CreateNotificationRequest{
			Info: &notificationpbpb.UserNotification{
				AppID:  withdrawItem.GetAppID(),
				UserID: withdrawItem.GetUserID(),
			},
			Message:  in.GetInfo().GetMessage(),
			LangID:   in.GetLangID(),
			UsedFor:  notificationconstant.UsedForWithdrawReviewRejectedNotification,
			UserName: user.GetExtra().GetUsername(),
		})
		if err != nil {
			return nil, fmt.Errorf("fail create notification: %v", err)
		}
	}

	return &npool.UpdateWithdrawReviewResponse{
		Info: &npool.WithdrawReview{
			Review:   reviewUpResp,
			User:     user,
			Withdraw: withdrawItem,
		},
	}, nil
}
