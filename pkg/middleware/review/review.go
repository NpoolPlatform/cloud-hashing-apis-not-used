package review

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	goodsconst "github.com/NpoolPlatform/cloud-hashing-goods/pkg/message/const"
	kycconst "github.com/NpoolPlatform/kyc-management/pkg/message/const"
	appusermgrpb "github.com/NpoolPlatform/message/npool/appusermgr"
	kycmgrpb "github.com/NpoolPlatform/message/npool/kyc"
	reviewpb "github.com/NpoolPlatform/message/npool/review-service"

	"golang.org/x/xerrors"
)

func GetKycReviews(ctx context.Context, in *npool.GetKycReviewsRequest) (*npool.GetKycReviewsResponse, error) {
	resp, err := grpc2.GetReviewsByAppDomain(ctx, &reviewpb.GetReviewsByAppDomainRequest{
		AppID:  in.GetAppID(),
		Domain: kycconst.ServiceName,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get kyc reviews: %v", err)
	}
	// TODO: Expand reviewer

	reviews := []*npool.KycReview{}
	for _, info := range resp.Infos {
		kyc, err := grpc2.GetKycByIDs(ctx, &kycmgrpb.GetKycByKycIDsRequest{
			KycIDs: []string{
				info.ObjectID,
			},
		})
		if err != nil {
			return nil, xerrors.Errorf("fail get kyc info for %v: %v", info.ID, err)
		}
		if len(kyc.Infos) == 0 {
			logger.Sugar().Warnf("empty kyc info for %v", info.ObjectID)
			continue
		}

		user, err := grpc2.GetAppUserInfoByAppUser(ctx, &appusermgrpb.GetAppUserInfoByAppUserRequest{
			AppID:  in.GetAppID(),
			UserID: kyc.Infos[0].UserID,
		})
		if err != nil {
			logger.Sugar().Errorf("fail get user %v info for %v: %v", kyc.Infos[0].UserID, info.ObjectID, err)
			continue
		}

		reviews = append(reviews, &npool.KycReview{
			Review: info,
			User:   user.Info,
			Kyc:    kyc.Infos[0],
		})
	}

	return &npool.GetKycReviewsResponse{
		Infos: reviews,
	}, nil
}

func GetGoodReviews(ctx context.Context, in *npool.GetGoodReviewsRequest) (*npool.GetGoodReviewsResponse, error) {
	resp, err := grpc2.GetReviewsByDomain(ctx, &reviewpb.GetReviewsByDomainRequest{
		Domain: goodsconst.ServiceName,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get good reviews: %v", err)
	}
	// TODO: Expand reviewer
	// TODO: Expand good

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
