package review

import (
	"context"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	goodsconst "github.com/NpoolPlatform/cloud-hashing-goods/pkg/message/const" //nolint
	reviewpb "github.com/NpoolPlatform/message/npool/review-service"            //nolint

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
			Review: info,
		})
	}

	return &npool.GetGoodReviewsResponse{
		Infos: reviews,
	}, nil
}
