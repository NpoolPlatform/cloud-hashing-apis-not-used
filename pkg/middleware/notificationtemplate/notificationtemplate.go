package notificationtemplate

import (
	"context"
	"strings"

	constant "github.com/NpoolPlatform/cloud-hashing-apis/pkg/const"
	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	notificationpbpb "github.com/NpoolPlatform/message/npool/notification"
	"golang.org/x/xerrors"
)

func GetTemplateByAppLangUsedFor(ctx context.Context, in *notificationpbpb.GetTemplateByAppLangUsedForRequest, message, userName string) (*notificationpbpb.Template, error) {
	template, err := grpc2.GetTemplateByAppLangUsedFor(ctx, in)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, xerrors.Errorf("fail get template")
	}
	if message != "" {
		template.Content = strings.ReplaceAll(template.Content, constant.MessageTemplate, message)
	}
	if userName != "" {
		template.Content = strings.ReplaceAll(template.Content, constant.NameTemplate, userName)
	}
	return template, err
}
