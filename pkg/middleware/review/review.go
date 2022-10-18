package review

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	billingconst "github.com/NpoolPlatform/cloud-hashing-billing/pkg/message/const"
	billingpb "github.com/NpoolPlatform/message/npool/cloud-hashing-billing"
	reviewpb "github.com/NpoolPlatform/message/npool/review-service"
	reviewconst "github.com/NpoolPlatform/review-service/pkg/const"

	appusermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"
)

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

		user, err := appusermwcli.GetUser(ctx, in.GetAppID(), address.UserID)
		if err != nil {
			logger.Sugar().Errorf("fail get user %v info for %v: %v", address.UserID, info.ObjectID, err)
			continue
		}

		if user == nil {
			logger.Sugar().Errorf("fail get user info ")
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
