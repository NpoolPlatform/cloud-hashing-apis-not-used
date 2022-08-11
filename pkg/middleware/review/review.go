package review

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	billingconst "github.com/NpoolPlatform/cloud-hashing-billing/pkg/message/const"
	goodsconst "github.com/NpoolPlatform/cloud-hashing-goods/pkg/message/const"
	kycconst "github.com/NpoolPlatform/kyc-management/pkg/message/const"
	appusermgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v1"
	billingpb "github.com/NpoolPlatform/message/npool/cloud-hashing-billing"
	kycmgrpb "github.com/NpoolPlatform/message/npool/kyc"
	reviewpb "github.com/NpoolPlatform/message/npool/review-service"
	reviewconst "github.com/NpoolPlatform/review-service/pkg/const"
)

func GetKycReviews(ctx context.Context, in *npool.GetKycReviewsRequest) (*npool.GetKycReviewsResponse, error) {
	infos, err := grpc2.GetReviewsByAppDomain(ctx, &reviewpb.GetReviewsByAppDomainRequest{
		AppID:  in.GetAppID(),
		Domain: kycconst.ServiceName,
	})
	if err != nil {
		return nil, fmt.Errorf("fail get kyc reviews: %v", err)
	}
	// TODO: Expand reviewer

	reviews := []*npool.KycReview{}
	for _, info := range infos {
		kycs, err := grpc2.GetKycByIDs(ctx, &kycmgrpb.GetKycByKycIDsRequest{
			KycIDs: []string{
				info.ObjectID,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("fail get kyc info for %v: %v", info.ID, err)
		}
		if len(kycs) == 0 {
			logger.Sugar().Warnf("empty kyc info for %v", info.ObjectID)
			continue
		}

		user, err := grpc2.GetAppUserInfoByAppUser(ctx, &appusermgrpb.GetAppUserInfoByAppUserRequest{
			AppID:  in.GetAppID(),
			UserID: kycs[0].UserID,
		})
		if err != nil {
			logger.Sugar().Errorf("fail get user %v info for %v: %v", kycs[0].UserID, info.ObjectID, err)
			continue
		}

		reviews = append(reviews, &npool.KycReview{
			Review: info,
			User:   user,
			Kyc:    kycs[0],
		})
	}

	return &npool.GetKycReviewsResponse{
		Infos: reviews,
	}, nil
}

func GetGoodReviews(ctx context.Context, in *npool.GetGoodReviewsRequest) (*npool.GetGoodReviewsResponse, error) {
	infos, err := grpc2.GetReviewsByDomain(ctx, &reviewpb.GetReviewsByDomainRequest{
		Domain: goodsconst.ServiceName,
	})
	if err != nil {
		return nil, fmt.Errorf("fail get good reviews: %v", err)
	}
	// TODO: Expand reviewer
	// TODO: Expand good

	reviews := []*npool.GoodReview{}
	for _, info := range infos {
		reviews = append(reviews, &npool.GoodReview{
			Review: info,
		})
	}

	return &npool.GetGoodReviewsResponse{
		Infos: reviews,
	}, nil
}

func GetWithdrawReviews(ctx context.Context, in *npool.GetWithdrawReviewsRequest) (*npool.GetWithdrawReviewsResponse, error) {
	infos, err := grpc2.GetReviewsByAppDomain(ctx, &reviewpb.GetReviewsByAppDomainRequest{
		AppID:  in.GetAppID(),
		Domain: billingconst.ServiceName,
	})
	if err != nil {
		return nil, fmt.Errorf("fail get withdraw reviews: %v", err)
	}

	reviews := []*npool.WithdrawReview{}
	for _, info := range infos {
		item, err := grpc2.GetUserWithdrawItem(ctx, &billingpb.GetUserWithdrawItemRequest{
			ID: info.ObjectID,
		})
		if err != nil {
			return nil, fmt.Errorf("fail get user withdraw info for %v: %v", info.ID, err)
		}
		if item == nil {
			logger.Sugar().Warnf("fail get user withdraw info for %v", info.ObjectID)
			continue
		}

		user, err := grpc2.GetAppUserInfoByAppUser(ctx, &appusermgrpb.GetAppUserInfoByAppUserRequest{
			AppID:  in.GetAppID(),
			UserID: item.UserID,
		})
		if err != nil {
			logger.Sugar().Errorf("fail get user %v info for %v: %v", item.UserID, info.ObjectID, err)
			continue
		}

		reviews = append(reviews, &npool.WithdrawReview{
			Review:   info,
			User:     user,
			Withdraw: item,
		})
	}

	return &npool.GetWithdrawReviewsResponse{
		Infos: reviews,
	}, nil
}

func GetWithdrawAddressReviews(ctx context.Context, in *npool.GetWithdrawAddressReviewsRequest) (*npool.GetWithdrawAddressReviewsResponse, error) {
	infos, err := grpc2.GetReviewsByAppDomain(ctx, &reviewpb.GetReviewsByAppDomainRequest{
		AppID:  in.GetAppID(),
		Domain: billingconst.ServiceName,
	})
	if err != nil {
		return nil, fmt.Errorf("fail get withdraw reviews: %v", err)
	}

	reviews := []*npool.WithdrawAddressReview{}
	for _, info := range infos {
		address, err := grpc2.GetUserWithdraw(ctx, &billingpb.GetUserWithdrawRequest{
			ID: info.ObjectID,
		})
		if err != nil {
			return nil, fmt.Errorf("fail get user withdraw info for %v: %v", info.ID, err)
		}
		if address == nil {
			logger.Sugar().Warnf("fail get user withdraw info for %v", info.ObjectID)
			continue
		}

		account, err := grpc2.GetBillingAccount(ctx, &billingpb.GetCoinAccountRequest{
			ID: address.AccountID,
		})
		if err != nil || account == nil {
			return nil, fmt.Errorf("fail get account: %v", err)
		}

		user, err := grpc2.GetAppUserInfoByAppUser(ctx, &appusermgrpb.GetAppUserInfoByAppUserRequest{
			AppID:  in.GetAppID(),
			UserID: address.UserID,
		})
		if err != nil {
			logger.Sugar().Errorf("fail get user %v info for %v: %v", address.UserID, info.ObjectID, err)
			continue
		}

		reviews = append(reviews, &npool.WithdrawAddressReview{
			Review:  info,
			User:    user,
			Address: address,
			Account: account,
		})
	}

	return &npool.GetWithdrawAddressReviewsResponse{
		Infos: reviews,
	}, nil
}

func GetReviewState(ctx context.Context, in *reviewpb.GetReviewsByAppDomainObjectTypeIDRequest) (string, string, error) { //nolint
	infos, err := grpc2.GetReviewsByAppDomainObjectTypeID(ctx, in)
	if err != nil {
		return "", "", fmt.Errorf("fail get review: %v", err)
	}

	reviewState := reviewconst.StateRejected
	reviewMessage := ""
	messageTime := uint32(0)

	for _, info := range infos {
		if info.State == reviewconst.StateWait {
			reviewState = reviewconst.StateWait
			break
		}
	}

	for _, info := range infos {
		if info.State == reviewconst.StateApproved {
			reviewState = reviewconst.StateApproved
			break
		}
	}

	if reviewState == reviewconst.StateRejected {
		for _, info := range infos {
			if info.State == reviewconst.StateRejected {
				if messageTime < info.CreateAt {
					reviewMessage = info.Message
					messageTime = info.CreateAt
				}
			}
		}
	}

	return reviewState, reviewMessage, nil
}
