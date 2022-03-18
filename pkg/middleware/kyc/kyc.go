package kyc

import (
	"context"

	constant "github.com/NpoolPlatform/cloud-hashing-apis/pkg/const"
	review "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/review"
	kycmgrconst "github.com/NpoolPlatform/kyc-management/pkg/message/const"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"
	kycmgrpb "github.com/NpoolPlatform/message/npool/kyc"

	notificationpbpb "github.com/NpoolPlatform/message/npool/notification"

	reviewpb "github.com/NpoolPlatform/message/npool/review-service"
	reviewconst "github.com/NpoolPlatform/review-service/pkg/const"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"

	"golang.org/x/xerrors"
)

func Create(ctx context.Context, in *npool.CreateKycRequest) (*npool.CreateKycResponse, error) {
	resp, err := grpc2.CreateKyc(ctx, &kycmgrpb.CreateKycRequest{
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
			ObjectID:   resp.Info.ID,
		},
	})
	if err != nil {
		// TODO: rollback kyc database
		return nil, xerrors.Errorf("fail create kyc review: %v", err)
	}

	return &npool.CreateKycResponse{
		Info: &npool.Kyc{
			Kyc:   resp.Info,
			State: reviewconst.StateWait,
		},
	}, nil
}

func Update(ctx context.Context, in *npool.UpdateKycRequest) (*npool.UpdateKycResponse, error) {
	allowed := true
	reviewing := false

	_review, err := grpc2.GetReviewsByAppDomainObjectTypeID(ctx, &reviewpb.GetReviewsByAppDomainObjectTypeIDRequest{
		AppID:      in.GetInfo().GetAppID(),
		Domain:     kycmgrconst.ServiceName,
		ObjectType: constant.ReviewObjectKyc,
		ObjectID:   in.GetInfo().GetID(),
	})
	if err == nil {
		for _, info := range _review.Infos {
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

	resp, err := grpc2.UpdateKyc(ctx, &kycmgrpb.UpdateKycRequest{
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
				ObjectID:   resp.Info.ID,
			},
		})
		if err != nil {
			// TODO: rollback kyc database
			return nil, xerrors.Errorf("fail create kyc review: %v", err)
		}
	}

	return &npool.UpdateKycResponse{
		Info: &npool.Kyc{
			Kyc:   resp.Info,
			State: reviewconst.StateWait,
		},
	}, nil
}

func GetByAppUser(ctx context.Context, in *npool.GetKycByAppUserRequest) (*npool.GetKycByAppUserResponse, error) {
	resp, err := grpc2.GetKycByUserID(ctx, &kycmgrpb.GetKycByUserIDRequest{
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
		ObjectID:   resp.Info.ID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get review state: %v", err)
	}

	return &npool.GetKycByAppUserResponse{
		Info: &npool.Kyc{
			Kyc:     resp.Info,
			State:   reviewState,
			Message: reviewMessage,
		},
	}, nil
}

func UpdateKycReview(ctx context.Context, in *npool.UpdateKycReviewRequest) (*npool.UpdateKycReviewResponse, error) {

	reviewInfo := in.GetReview().GetInfo()
	reviewResp, err := grpc2.GetReview(ctx, &reviewpb.GetReviewRequest{
		ID: reviewInfo.GetID(),
	})
	if err != nil {
		// TODO: rollback kyc database
		return nil, xerrors.Errorf("fail get review: %v", err)
	}
	resp, err := grpc2.UpdateReview(ctx, &reviewpb.UpdateReviewRequest{
		Info: reviewInfo,
	})

	if reviewResp.GetInfo().GetState() == reviewconst.StateWait && reviewInfo.GetState() == reviewconst.StateApproved {
		_, err := grpc2.CreateNotification(ctx, &notificationpbpb.CreateNotificationRequest{
			Info: &notificationpbpb.UserNotification{
				AppID:   reviewInfo.GetAppID(),
				UserID:  reviewInfo.GetObjectID(),
				Title:   "kyc消息通知",
				Content: "测试kyc消息通知",
			},
		})
		if err != nil {
			return nil, err
		}
	}

	return &npool.UpdateKycReviewResponse{Info: resp}, nil
}
