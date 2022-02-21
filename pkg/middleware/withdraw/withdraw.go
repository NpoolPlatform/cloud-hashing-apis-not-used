package withdraw

import (
	"context"
	"fmt"
	"time"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	constant "github.com/NpoolPlatform/cloud-hashing-apis/pkg/const"
	review "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/review"
	currency "github.com/NpoolPlatform/cloud-hashing-staker/pkg/middleware/currency"
	redis2 "github.com/NpoolPlatform/go-service-framework/pkg/redis"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	appusermgrconst "github.com/NpoolPlatform/appuser-manager/pkg/const"
	billingstate "github.com/NpoolPlatform/cloud-hashing-billing/pkg/const"
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

	"github.com/google/uuid"
	"golang.org/x/xerrors"
)

func Create(ctx context.Context, in *npool.SubmitUserWithdrawRequest) (*npool.SubmitUserWithdrawResponse, error) { //nolint
	if in.GetInfo().GetAmount() <= 0 {
		return nil, xerrors.Errorf("invalid amount")
	}

	user, err := grpc2.GetAppUserByAppUser(ctx, &appusermgrpb.GetAppUserByAppUserRequest{
		AppID:  in.GetInfo().GetAppID(),
		UserID: in.GetInfo().GetUserID(),
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
			AppID:   in.GetInfo().GetAppID(),
			PhoneNO: phoneNO,
			UsedFor: thirdgwconst.UsedForWithdraw,
			Code:    in.GetVerificationCode(),
		})
	} else if in.GetAccountType() == appusermgrconst.SignupByEmail {
		emailAddr := in.GetAccount()

		if user.Info.EmailAddress != emailAddr {
			return nil, xerrors.Errorf("invalid email address")
		}

		_, err = grpc2.VerifyEmailCode(ctx, &thirdgwpb.VerifyEmailCodeRequest{
			AppID:        in.GetInfo().GetAppID(),
			EmailAddress: emailAddr,
			UsedFor:      thirdgwconst.UsedForWithdraw,
			Code:         in.GetVerificationCode(),
		})
	} else {
		_, err = grpc2.VerifyGoogleAuthentication(ctx, &thirdgwpb.VerifyGoogleAuthenticationRequest{
			AppID:  in.GetInfo().GetAppID(),
			UserID: in.GetInfo().GetUserID(),
			Code:   in.GetVerificationCode(),
		})
	}
	if err != nil {
		return nil, xerrors.Errorf("fail verify code: %v", err)
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
	if coin.Info.PreSale {
		return nil, xerrors.Errorf("cannot withdraw presale coin")
	}

	lockKey := fmt.Sprintf("withdraw:%v:%v", in.GetInfo().GetAppID(), in.GetInfo().GetUserID())
	err = redis2.TryLock(fmt.Sprintf("withdraw:%v:%v", in.GetInfo().GetAppID(), in.GetInfo().GetUserID()), 10*time.Minute)
	if err != nil {
		return nil, xerrors.Errorf("lock withdraw fail: %v", err)
	}
	defer func() {
		if err := redis2.Unlock(lockKey); err != nil {
			logger.Sugar().Errorf("unlock withdraw fail: %v", err)
		}
	}()

	benefits, err := grpc2.GetUserBenefitsByAppUserCoin(ctx, &billingpb.GetUserBenefitsByAppUserCoinRequest{
		AppID:      in.GetInfo().GetAppID(),
		UserID:     in.GetInfo().GetUserID(),
		CoinTypeID: in.GetInfo().GetCoinTypeID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get user benefits: %v", err)
	}

	withdrawAddrs, err := grpc2.GetUserWithdrawsByAppUserCoin(ctx, &billingpb.GetUserWithdrawsByAppUserCoinRequest{
		AppID:      in.GetInfo().GetAppID(),
		UserID:     in.GetInfo().GetUserID(),
		CoinTypeID: in.GetInfo().GetCoinTypeID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get user withdraws: %v", err)
	}

	txs, err := grpc2.GetCoinAccountTransactionsByAppUserCoin(ctx, &billingpb.GetCoinAccountTransactionsByAppUserCoinRequest{
		AppID:      in.GetInfo().GetAppID(),
		UserID:     in.GetInfo().GetUserID(),
		CoinTypeID: in.GetInfo().GetCoinTypeID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get account transactions: %v", err)
	}

	incoming := 0.0
	outcoming := 0.0
	for _, info := range benefits.Infos {
		incoming += info.Amount
	}
	for _, info := range txs.Infos {
		withdraw := false
		for _, addr := range withdrawAddrs.Infos {
			if addr.AccountID == info.ToAddressID {
				withdraw = true
				break
			}
		}

		if !withdraw {
			continue
		}

		if info.State == billingstate.CoinTransactionStateFail ||
			info.State == billingstate.CoinTransactionStateRejected {
			continue
		}

		outcoming += info.Amount
	}

	if incoming < outcoming {
		return nil, xerrors.Errorf("invalid billing input")
	}
	if incoming-outcoming < in.GetInfo().GetAmount() {
		return nil, xerrors.Errorf("not sufficient funds")
	}

	// TODO: check waiting transaction: only one wait transaction is allowed
	// TODO: check reviewing withdraw: only one reviewing withdraw is allowed

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
		return nil, xerrors.Errorf("acount is not belong to user")
	}

	reviewState, _, err := review.GetReviewState(ctx, &reviewpb.GetReviewsByAppDomainObjectTypeIDRequest{
		AppID:      in.GetInfo().GetAppID(),
		Domain:     billingconst.ServiceName,
		ObjectType: constant.ReviewObjectUserWithdrawAddress,
		ObjectID:   withdrawAccount.Info.ID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get review: %v", err)
	}
	if reviewState != reviewconst.StateApproved {
		return nil, xerrors.Errorf("invalid account: not approved")
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

	_review, err := grpc2.CreateReview(ctx, &reviewpb.CreateReviewRequest{
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

	// TODO: check 24 hours total withdraw amount: if overflow, goto review
	reviewState = reviewconst.StateWait

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
				ChainTransactionID: "",
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

		_review.Info.State = reviewconst.StateApproved
		_, err = grpc2.UpdateReview(ctx, &reviewpb.UpdateReviewRequest{
			Info: _review.Info,
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

func Update(ctx context.Context, in *npool.UpdateUserWithdrawReviewRequest) (*npool.UpdateUserWithdrawReviewResponse, error) { //nolint
	// TODO: check permission of reviewer

	if in.GetUserID() != in.GetInfo().GetReviewerID() {
		return nil, xerrors.Errorf("mismatch reviewer id")
	}

	user, err := grpc2.GetAppUserInfoByAppUser(ctx, &appusermgrpb.GetAppUserInfoByAppUserRequest{
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
		ID: in.GetInfo().GetID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get review: %v", err)
	}
	if resp.Info == nil {
		return nil, xerrors.Errorf("fail get review")
	}
	if resp.Info.State == reviewconst.StateApproved {
		return nil, xerrors.Errorf("already approved")
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

	invalidID := uuid.UUID{}.String()
	if resp1.Info.PlatformTransactionID != invalidID {
		return nil, xerrors.Errorf("withdraw already processed")
	}

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

	// TODO: here should hold withdraw lock
	benefits, err := grpc2.GetUserBenefitsByAppUserCoin(ctx, &billingpb.GetUserBenefitsByAppUserCoinRequest{
		AppID:      resp1.Info.AppID,
		UserID:     resp1.Info.UserID,
		CoinTypeID: resp1.Info.CoinTypeID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get user benefits: %v", err)
	}

	withdrawAddrs, err := grpc2.GetUserWithdrawsByAppUserCoin(ctx, &billingpb.GetUserWithdrawsByAppUserCoinRequest{
		AppID:      resp1.Info.AppID,
		UserID:     resp1.Info.UserID,
		CoinTypeID: resp1.Info.CoinTypeID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get user withdraws: %v", err)
	}

	txs, err := grpc2.GetCoinAccountTransactionsByAppUserCoin(ctx, &billingpb.GetCoinAccountTransactionsByAppUserCoinRequest{
		AppID:      resp1.Info.AppID,
		UserID:     resp1.Info.UserID,
		CoinTypeID: resp1.Info.CoinTypeID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get account transactions: %v", err)
	}

	incoming := 0.0
	outcoming := 0.0
	for _, info := range benefits.Infos {
		incoming += info.Amount
	}
	for _, info := range txs.Infos {
		withdraw := false
		for _, addr := range withdrawAddrs.Infos {
			if addr.AccountID == info.ToAddressID {
				withdraw = true
				break
			}
		}

		if !withdraw {
			continue
		}

		if info.State == billingstate.CoinTransactionStateFail ||
			info.State == billingstate.CoinTransactionStateRejected {
			continue
		}

		outcoming += info.Amount
	}

	if incoming < outcoming {
		return nil, xerrors.Errorf("invalid billing input")
	}
	if incoming-outcoming < resp1.Info.Amount {
		return nil, xerrors.Errorf("not sufficient funds")
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

	// TODO: here should hold transfer lock
	balance, err := grpc2.GetBalance(ctx, &sphinxproxypb.GetBalanceRequest{
		Name:    coin.Info.Name,
		Address: account.Info.Address,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get wallet balance: %v", err)
	}

	if balance.Info.Balance < resp1.Info.Amount+coin.Info.ReservedAmount {
		return nil, xerrors.Errorf("not sufficient funds")
	}

	account, err = grpc2.GetBillingAccount(ctx, &billingpb.GetCoinAccountRequest{
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

	reviewState, _, err := review.GetReviewState(ctx, &reviewpb.GetReviewsByAppDomainObjectTypeIDRequest{
		AppID:      in.GetInfo().GetAppID(),
		Domain:     billingconst.ServiceName,
		ObjectType: constant.ReviewObjectUserWithdrawAddress,
		ObjectID:   withdrawAccount.Info.ID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get review: %v", err)
	}
	if reviewState != reviewconst.StateApproved {
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

	if resp1.Info.AppID != in.GetInfo().GetAppID() {
		return nil, xerrors.Errorf("invalid request")
	}

	if resp.Info.State != reviewconst.StateWait {
		return nil, xerrors.Errorf("already reviewed")
	}

	if in.GetInfo().GetState() == reviewconst.StateApproved {
		account, err := grpc2.GetBillingAccount(ctx, &billingpb.GetCoinAccountRequest{
			ID: coinsetting.Info.UserOnlineAccountID,
		})
		if err != nil {
			return nil, xerrors.Errorf("fail get account info: %v", err)
		}
		if account.Info == nil {
			return nil, xerrors.Errorf("fail get account info")
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
				ChainTransactionID: "",
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
	}

	resp.Info.State = in.GetInfo().GetState()
	resp2, err := grpc2.UpdateReview(ctx, &reviewpb.UpdateReviewRequest{
		Info: resp.Info,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail update review state: %v", err)
	}

	return &npool.UpdateUserWithdrawReviewResponse{
		Info: &npool.WithdrawReview{
			Withdraw: resp1.Info,
			Review:   resp2.Info,
			User:     user.Info,
		},
	}, nil
}

func GetByAppUser(ctx context.Context, in *npool.GetUserWithdrawsByAppUserRequest) (*npool.GetUserWithdrawsByAppUserResponse, error) { //nolint
	resp, err := grpc2.GetUserWithdrawItemsByAppUser(ctx, &billingpb.GetUserWithdrawItemsByAppUserRequest{
		AppID:  in.GetAppID(),
		UserID: in.GetUserID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get user withdraw items: %v", err)
	}

	withdraws := []*npool.UserWithdraw{}
	for _, info := range resp.Infos {
		reviewState, reviewMessage, err := review.GetReviewState(ctx, &reviewpb.GetReviewsByAppDomainObjectTypeIDRequest{
			AppID:      info.AppID,
			Domain:     billingconst.ServiceName,
			ObjectType: constant.ReviewObjectWithdraw,
			ObjectID:   info.ID,
		})
		if err != nil {
			return nil, xerrors.Errorf("fail get review state: %v", err)
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
