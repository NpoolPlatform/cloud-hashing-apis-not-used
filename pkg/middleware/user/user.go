package user

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	"github.com/NpoolPlatform/cloud-hashing-apis/message/npool"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"

	inspirepb "github.com/NpoolPlatform/cloud-hashing-inspire/message/npool"
	usermgrpb "github.com/NpoolPlatform/user-management/message/npool"

	"golang.org/x/xerrors"
)

func convertUserinfo(info *usermgrpb.UserBasicInfo) *npool.UserInfo {
	return &npool.UserInfo{
		UserID:         info.UserID,
		Username:       info.Username,
		Avatar:         info.Avatar,
		Age:            info.Age,
		Gender:         info.Gender,
		Region:         info.Region,
		Birthday:       info.Birthday,
		Country:        info.Country,
		Province:       info.Province,
		City:           info.City,
		PhoneNumber:    info.PhoneNumber,
		EmailAddress:   info.EmailAddress,
		CreateAt:       info.CreateAt,
		LoginTimes:     info.LoginTimes,
		KycVerify:      info.KycVerify,
		GaVerify:       info.GaVerify,
		GaLogin:        info.GaLogin,
		SignupMethod:   info.SignupMethod,
		Career:         info.Career,
		DisplayName:    info.DisplayName,
		FirstName:      info.FirstName,
		LastName:       info.LastName,
		StreetAddress1: info.StreetAddress1,
		StreetAddress2: info.StreetAddress2,
	}
}

func Signup(ctx context.Context, in *npool.SignupRequest) (*npool.SignupResponse, error) {
	signupResp, err := grpc2.Signup(ctx, &usermgrpb.SignupRequest{
		Username:     in.GetUsername(),
		Password:     in.GetPassword(),
		EmailAddress: in.GetEmailAddress(),
		PhoneNumber:  in.GetPhoneNumber(),
		Code:         in.GetVerificationCode(),
		AppID:        in.GetAppID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail signup: %v", err)
	}

	invitationCode := in.GetInvitationCode()
	if invitationCode != "" {
		getByCodeResp, err := grpc2.GetUserInvitationCodeByCode(ctx, &inspirepb.GetUserInvitationCodeByCodeRequest{
			Code: invitationCode,
		})
		if err != nil {
			logger.Sugar().Errorf("fail get user invitation code: %v", err)
			return &npool.SignupResponse{
				Info: convertUserinfo(signupResp.Info),
			}, nil
		}

		if getByCodeResp.Info.AppID != in.GetAppID() {
			logger.Sugar().Errorf("invalid invitation code for app")
			return &npool.SignupResponse{
				Info: convertUserinfo(signupResp.Info),
			}, nil
		}

		_, err = grpc2.CreateRegistrationInvitation(ctx, &inspirepb.CreateRegistrationInvitationRequest{
			Info: &inspirepb.RegistrationInvitation{
				AppID:     in.GetAppID(),
				InviterID: getByCodeResp.Info.UserID,
				InviteeID: signupResp.Info.UserID,
			},
		})
		if err != nil {
			logger.Sugar().Errorf("fail create registration invitation: %v", err)
			return &npool.SignupResponse{
				Info: convertUserinfo(signupResp.Info),
			}, nil
		}
	}

	return &npool.SignupResponse{
		Info: convertUserinfo(signupResp.Info),
	}, nil
}
