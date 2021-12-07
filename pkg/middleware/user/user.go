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

func GetMyInvitations(ctx context.Context, in *npool.GetMyInvitationsRequest) (*npool.GetMyInvitationsResponse, error) { //nolint
	userResp, err := grpc2.GetUser(ctx, &usermgrpb.GetUserRequest{
		AppID:  in.GetAppID(),
		UserID: in.GetInviterID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get inviter user information: %v", err)
	}

	layer := 0
	goon := true

	type myInvitee struct {
		layer   int
		inviter *npool.UserInfo
		invitee *npool.UserInfo
	}
	tmpInvitees := []*myInvitee{&myInvitee{ //nolint
		layer:   layer,
		inviter: nil,
		invitee: &npool.UserInfo{
			UserID:   userResp.Info.UserID,
			Username: userResp.Info.Username,
			Avatar:   userResp.Info.Avatar,
		},
	}}
	layer++

	// TODO: process deadloop
	for goon {
		goon = false

		for _, curInvitee := range tmpInvitees {
			if curInvitee.layer != layer-1 {
				continue
			}

			resp, err := grpc2.GetRegistrationInvitationsByAppInviter(ctx, &inspirepb.GetRegistrationInvitationsByAppInviterRequest{
				AppID:     in.GetAppID(),
				InviterID: curInvitee.invitee.UserID,
			})
			if err != nil {
				logger.Sugar().Errorf("fail get invitations by inviter: %v", err)
				continue
			}

			for _, info := range resp.Infos {
				if info.AppID != in.GetAppID() ||
					info.InviterID != curInvitee.invitee.UserID {
					logger.Sugar().Errorf("invalid inviter id or app id")
					continue
				}

				inviteeResp, err := grpc2.GetUser(ctx, &usermgrpb.GetUserRequest{
					AppID:  in.GetAppID(),
					UserID: info.InviteeID,
				})
				if err != nil {
					logger.Sugar().Errorf("fail get invitee user info: %v", err)
					continue
				}

				tmpInvitees = append(tmpInvitees, &myInvitee{
					layer:   layer,
					inviter: curInvitee.invitee,
					invitee: &npool.UserInfo{
						UserID:   inviteeResp.Info.UserID,
						Username: inviteeResp.Info.Username,
						Avatar:   inviteeResp.Info.Avatar,
					},
				})

				goon = true
			}
		}

		layer++
	}

	invitations := []*npool.Invitation{}
	for _, curInvitee := range tmpInvitees {
		if curInvitee.inviter == nil {
			continue
		}

		inserted := false

		for _, invitation := range invitations {
			if invitation.Inviter.UserID == curInvitee.inviter.UserID {
				invitation.Invitees = append(invitation.Invitees, curInvitee.invitee)
				inserted = true
				break
			}
		}

		if !inserted {
			invitations = append(invitations, &npool.Invitation{
				Inviter:  curInvitee.inviter,
				Invitees: []*npool.UserInfo{curInvitee.invitee},
			})
		}
	}

	return &npool.GetMyInvitationsResponse{
		Infos: invitations,
	}, nil
}
