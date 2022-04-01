package kyc

import (
	"context"

	"entgo.io/ent/entc/integration/edgefield/ent/info"
	appusermgrpb "github.com/NpoolPlatform/message/npool/appusermgr"

	notificationconstant "github.com/NpoolPlatform/notification/pkg/const"

	constant "github.com/NpoolPlatform/cloud-hashing-apis/pkg/const"
	review "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/review"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	kycmgrconst "github.com/NpoolPlatform/kyc-management/pkg/message/const"
	kycmgrpb "github.com/NpoolPlatform/message/npool/kyc"

	notificationpbpb "github.com/NpoolPlatform/message/npool/notification"

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

func UpdateKycReview(ctx context.Context, in *npool.UpdateKycReviewRequest) (*npool.UpdateKycReviewResponse, error) {
	reviewInfo := in.GetInfo()

	kycs, err := grpc2.GetKycByIDs(ctx, &kycmgrpb.GetKycByKycIDsRequest{
		KycIDs: []string{
			reviewInfo.GetObjectID(),
		},
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get kyc info for %v: %v", info.ID, err)
	}
	if len(kycs) == 0 {
		return nil, xerrors.Errorf("empty kyc info for %v", reviewInfo.GetObjectID())
	}

	user, err := grpc2.GetAppUserInfoByAppUser(ctx, &appusermgrpb.GetAppUserInfoByAppUserRequest{
		AppID:  reviewInfo.GetAppID(),
		UserID: kycs[0].GetUserID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get app user: %v", err)
	}
	if user == nil {
		return nil, xerrors.Errorf("fail get app user")
	}

	reviewResp, err := grpc2.GetReview(ctx, &reviewpb.GetReviewRequest{
		ID: reviewInfo.GetID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get review: %v", err)
	}
	if reviewResp == nil {
		return nil, xerrors.Errorf("fail get review")
	}
	reviewUpResp, err := grpc2.UpdateReview(ctx, &reviewpb.UpdateReviewRequest{Info: in.GetInfo()})
	if err != nil {
		return nil, err
	}
	if reviewInfo.GetState() == reviewconst.StateApproved {
		_, err = grpc2.CreateNotification(ctx, &notificationpbpb.CreateNotificationRequest{
			Info: &notificationpbpb.UserNotification{
				AppID:  reviewInfo.GetAppID(),
				UserID: kycs[0].GetUserID(),
			},
			Message:  in.GetInfo().GetMessage(),
			LangID:   in.GetLangID(),
			UsedFor:  notificationconstant.UsedForKycReviewApprovedNotification,
			UserName: user.GetExtra().GetUsername(),
		})
		if err != nil {
			return nil, xerrors.Errorf("fail create notification: %v", err)
		}
	}

	if reviewInfo.GetState() == reviewconst.StateRejected {
		_, err = grpc2.CreateNotification(ctx, &notificationpbpb.CreateNotificationRequest{
			Info: &notificationpbpb.UserNotification{
				AppID:  reviewInfo.GetAppID(),
				UserID: kycs[0].GetUserID(),
			},
			Message:  in.GetInfo().GetMessage(),
			LangID:   in.GetLangID(),
			UsedFor:  notificationconstant.UsedForKycReviewRejectedNotification,
			UserName: user.GetExtra().GetUsername(),
		})
		if err != nil {
			return nil, xerrors.Errorf("fail create notification: %v", err)
		}
	}

	return &npool.UpdateKycReviewResponse{
		Info: &npool.KycReview{
			Review: reviewUpResp,
			User:   user,
			Kyc:    kycs[0],
		},
	}, nil
}
