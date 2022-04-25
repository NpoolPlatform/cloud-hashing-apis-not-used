package dtm

import (
	"context"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"
	"github.com/NpoolPlatform/go-service-framework/pkg/dtm"
	apimgrpb "github.com/NpoolPlatform/message/npool/apimgr"
)

func GetGrpcURL(ctx context.Context, serviceName, actionGrpcName, compensateGrpcName string) (action, compensate string, err error) {
	methodName := []string{actionGrpcName, compensateGrpcName}
	grpcURL, err := GetPath(ctx, serviceName, methodName)
	if err != nil {
		return "", "", err
	}

	serviceName, err = dtm.GetService(serviceName)
	if err != nil {
		return "", "", err
	}

	actionPath := ""
	compensatePath := ""
	for _, val := range grpcURL {
		if val.MethodName == actionGrpcName {
			actionPath = val.Path
		}
		if val.MethodName == compensateGrpcName {
			compensatePath = val.Path
		}
	}

	actionGrpcURL := serviceName + "/" + actionPath
	compensateGrpcURL := serviceName + "/" + compensatePath
	return actionGrpcURL, compensateGrpcURL, nil
}

func GetPath(ctx context.Context, serviceName string, methodName []string) ([]*apimgrpb.ServicePath, error) {
	grpcApis, err := grpc2.GetAPIByServiceMethod(ctx, &apimgrpb.GetApisByServiceMethodRequest{
		ServiceName: serviceName,
		MethodName:  methodName,
	})
	if err != nil {
		return nil, err
	}
	return grpcApis.Infos, nil
}
