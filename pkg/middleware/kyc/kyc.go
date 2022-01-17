package kyc

import (
	"context"

	constant "github.com/NpoolPlatform/cloud-hashing-apis/pkg/const"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	kycmgrpb "github.com/NpoolPlatform/message/npool/kyc"
	reviewpb "github.com/NpoolPlatform/message/npool/review-service"
	reviewconst "github.com/NpoolPlatform/review-service/pkg/const"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"

	"golang.org/x/xerrors"
)

func GetByAppUser(ctx context.Context, in *npool.GetKycByAppUserRequest) (*npool.GetKycByAppUserResponse, error) {
	resp, err := grpc2.GetKycByUserID(ctx, &kycmgrpb.GetKycByUserIDRequest{
		AppID:  in.GetAppID(),
		UserID: in.GetUserID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get kyc: %v", err)
	}

	review, err := grpc2.GetReviewsByAppDomainObjectTypeID(ctx, &reviewpb.GetReviewsByAppDomainObjectTypeIDRequest{
		AppID:      in.GetAppID(),
		ObjectType: constant.ReviewObjectKyc,
		ObjectID:   resp.Info.ID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get kyc review for %v: %v", resp.Info.ID, err)
	}
	if len(review.Infos) == 0 {
		return nil, xerrors.Errorf("empty kyc review for %v", resp.Info.ID)
	}

	reviewState := reviewconst.StateRejected

	for _, info := range review.Infos {
		if info.State == reviewconst.StateWait {
			reviewState = reviewconst.StateWait
			break
		}
	}

	for _, info := range review.Infos {
		if info.State == reviewconst.StateApproved {
			reviewState = reviewconst.StateApproved
			break
		}
	}

	return &npool.GetKycByAppUserResponse{
		Info: &npool.Kyc{
			Kyc:   resp.Info,
			State: reviewState,
		},
	}, nil
}
