package review

import (
	"context"

	"github.com/NpoolPlatform/cloud-hashing-apis/message/npool"
	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"

	goodsconst "github.com/NpoolPlatform/cloud-hashing-goods/pkg/message/const" //nolint
	reviewpb "github.com/NpoolPlatform/review-service/message/npool"            //nolint

	"golang.org/x/xerrors"
)

func GetKycReviews(ctx context.Context, in *npool.GetKycReviewsRequest) (*npool.GetKycReviewsResponse, error) {
	return nil, nil
}

func GetGoodReviews(ctx context.Context, in *npool.GetGoodReviewsRequest) (*npool.GetGoodReviewsResponse, error) {
	// Get all good reviews
	resp, err := grpc2.GetReviewsByDomain(ctx, &reviewpb.GetReviewsByDomainRequest{
		Domain: goodsconst.ServiceName,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get good reviews: %v", err)
	}
	// Expand reviewer
	// Expand good

	reviews := []*npool.GoodReview{}
	for _, info := range resp.Infos {
		reviews = append(reviews, &npool.GoodReview{
			Review: &npool.Review{
				ID:         info.ID,
				ObjectType: info.ObjectType,
				AppID:      info.AppID,
				ReviewerID: info.ReviewerID,
				State:      info.State,
				Message:    info.Message,
				ObjectID:   info.ObjectID,
				Domain:     info.Domain,
			},
		})
	}

	return &npool.GetGoodReviewsResponse{
		Infos: reviews,
	}, nil
}
