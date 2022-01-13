package user

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"

	appmgrpb "github.com/NpoolPlatform/message/npool/application"
	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"
	usermgrpb "github.com/NpoolPlatform/message/npool/user"

	"golang.org/x/xerrors"
)

func Signup(ctx context.Context, in *npool.SignupRequest) (*npool.SignupResponse, error) { //nolint
	invitationCode := in.GetInvitationCode()
	inviterID := ""

	appResp, err := grpc2.GetApp(ctx, &appmgrpb.GetApplicationRequest{
		AppID: in.GetAppID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get app: %v", err)
	}

	if appResp.Info.InvitationCodeMust {
		if invitationCode == "" {
			return nil, xerrors.Errorf("invitation code is must")
		}
	}

	if invitationCode != "" {
		getByCodeResp, err := grpc2.GetUserInvitationCodeByCode(ctx, &inspirepb.GetUserInvitationCodeByCodeRequest{
			Code: invitationCode,
		})
		if err != nil {
			return nil, xerrors.Errorf("fail get user invitation code: %v", err)
		}

		if getByCodeResp.Info == nil {
			if appResp.Info.InvitationCodeMust {
				return nil, xerrors.Errorf("fail get invitation code")
			}
		} else {
			if getByCodeResp.Info.AppID != in.GetAppID() {
				return nil, xerrors.Errorf("invalid invitation code for app")
			}
			inviterID = getByCodeResp.Info.UserID
		}
	}

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

	if invitationCode != "" && inviterID != "" {
		_, err = grpc2.CreateRegistrationInvitation(ctx, &inspirepb.CreateRegistrationInvitationRequest{
			Info: &inspirepb.RegistrationInvitation{
				AppID:     in.GetAppID(),
				InviterID: inviterID,
				InviteeID: signupResp.Info.UserID,
			},
		})
		if err != nil {
			return nil, xerrors.Errorf("fail create registration invitation: %v", err)
		}
	}

	return &npool.SignupResponse{
		Info: signupResp.Info,
	}, nil
}

func GetMyInvitations(ctx context.Context, in *npool.GetMyInvitationsRequest) (*npool.GetMyInvitationsResponse, error) { //nolint
	addWatcher(in.GetAppID(), in.GetInviterID())

	invitations, userInfo, err := getFullInvitations(in.GetAppID(), in.GetInviterID())
	if err != nil {
		return nil, xerrors.Errorf("fail get invitations: %v", err)
	}
	return &npool.GetMyInvitationsResponse{
		MySelf: userInfo,
		Infos:  invitations,
	}, nil
}

func GetMyDirectInvitations(ctx context.Context, in *npool.GetMyDirectInvitationsRequest) (*npool.GetMyDirectInvitationsResponse, error) { //nolint
	addWatcher(in.GetAppID(), in.GetInviterID())

	invitations, userInfo, err := getDirectInvitations(in.GetAppID(), in.GetInviterID())
	if err != nil {
		return nil, xerrors.Errorf("fail get invitations: %v", err)
	}
	return &npool.GetMyDirectInvitationsResponse{
		MySelf: userInfo,
		Infos:  invitations,
	}, nil
}
