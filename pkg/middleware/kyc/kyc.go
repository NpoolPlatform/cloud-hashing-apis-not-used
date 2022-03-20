package kyc

import (
	"context"

	constant "github.com/NpoolPlatform/cloud-hashing-apis/pkg/const"
	review "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/review"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	kycmgrconst "github.com/NpoolPlatform/kyc-management/pkg/message/const"
	kycmgrpb "github.com/NpoolPlatform/message/npool/kyc"
	reviewpb "github.com/NpoolPlatform/message/npool/review-service"
	reviewconst "github.com/NpoolPlatform/review-service/pkg/const"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"

	"golang.org/x/xerrors"
)

func Create(ctx context.Context, in *npool.CreateKycRequest) (*npool.CreateKycResponse, error) {
	kyc, err := grpc2.CreateKyc(ctx, &kycmgrpb.CreateKycRequest{
		Info: in.GetInfo(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail create kyc: %v", err)
	}

	_, err = grpc2.CreateReview(ctx, &reviewpb.CreateReviewRequest{
		Info: &reviewpb.Review{
			AppID:      in.GetInfo().GetAppID(),
			Domain:     kycmgrconst.ServiceName,
			ObjectType: constant.ReviewObjectKyc,
			ObjectID:   kyc.ID,
		},
	})
	if err != nil {
		// TODO: rollback kyc database
		return nil, xerrors.Errorf("fail create kyc review: %v", err)
	}

	return &npool.CreateKycResponse{
		Info: &npool.Kyc{
			Kyc:   kyc,
			State: reviewconst.StateWait,
		},
	}, nil
}

func Update(ctx context.Context, in *npool.UpdateKycRequest) (*npool.UpdateKycResponse, error) {
	allowed := true
	reviewing := false

	reviews, err := grpc2.GetReviewsByAppDomainObjectTypeID(ctx, &reviewpb.GetReviewsByAppDomainObjectTypeIDRequest{
		AppID:      in.GetInfo().GetAppID(),
		Domain:     kycmgrconst.ServiceName,
		ObjectType: constant.ReviewObjectKyc,
		ObjectID:   in.GetInfo().GetID(),
	})
	if err == nil {
		for _, info := range reviews {
			if info.State == reviewconst.StateApproved {
				allowed = false
			}
			if info.State == reviewconst.StateWait {
				reviewing = true
			}
		}
	}

	if !allowed {
		return nil, xerrors.Errorf("not allowed update kyc")
	}

	kyc, err := grpc2.UpdateKyc(ctx, &kycmgrpb.UpdateKycRequest{
		Info: in.GetInfo(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail update kyc: %v", err)
	}

	if !reviewing {
		_, err = grpc2.CreateReview(ctx, &reviewpb.CreateReviewRequest{
			Info: &reviewpb.Review{
				AppID:      in.GetInfo().GetAppID(),
				Domain:     kycmgrconst.ServiceName,
				ObjectType: constant.ReviewObjectKyc,
				ObjectID:   kyc.ID,
			},
		})
		if err != nil {
			// TODO: rollback kyc database
			return nil, xerrors.Errorf("fail create kyc review: %v", err)
		}
	}

	return &npool.UpdateKycResponse{
		Info: &npool.Kyc{
			Kyc:   kyc,
			State: reviewconst.StateWait,
		},
	}, nil
}

func GetByAppUser(ctx context.Context, in *npool.GetKycByAppUserRequest) (*npool.GetKycByAppUserResponse, error) {
	kyc, err := grpc2.GetKycByUserID(ctx, &kycmgrpb.GetKycByUserIDRequest{
		AppID:  in.GetAppID(),
		UserID: in.GetUserID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get kyc: %v", err)
	}

	reviewState, reviewMessage, err := review.GetReviewState(ctx, &reviewpb.GetReviewsByAppDomainObjectTypeIDRequest{
		AppID:      in.GetAppID(),
		Domain:     kycmgrconst.ServiceName,
		ObjectType: constant.ReviewObjectKyc,
		ObjectID:   kyc.ID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get review state: %v", err)
	}

	return &npool.GetKycByAppUserResponse{
		Info: &npool.Kyc{
			Kyc:     kyc,
			State:   reviewState,
			Message: reviewMessage,
		},
	}, nil
}
