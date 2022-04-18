package dtm

import (
	appusermgrconst "github.com/NpoolPlatform/appuser-manager/pkg/message/const"
	inspireconst "github.com/NpoolPlatform/cloud-hashing-inspire/pkg/message/const"
	"github.com/NpoolPlatform/go-service-framework/pkg/dtm"
)

func GetGrpcURL(serviceName, actionGrpcName, compensateGrpcName string, service ...string) (action, compensate string, err error) {
	grpcURL, err := SetPackageAndService(serviceName, service...)
	if err != nil {
		return "", "", err
	}
	actionGrpcURL := grpcURL + "/" + actionGrpcName
	compensateGrpcURL := grpcURL + "/" + compensateGrpcName
	return actionGrpcURL, compensateGrpcURL, nil
}

func SetPackageAndService(serviceName string, service ...string) (string, error) {
	switch serviceName {
	case appusermgrconst.ServiceName:
		serviceName, err := dtm.GetService(serviceName)
		if err != nil {
			return "", err
		}
		if len(service) != 0 {
			return serviceName + "/app.user.manager.v1." + service[0], nil
		}
		return serviceName + "/app.user.manager.v1.AppUserManager", nil
	case inspireconst.ServiceName:
		serviceName, err := dtm.GetService(serviceName)
		if err != nil {
			return "", err
		}
		if len(service) != 0 {
			return serviceName + "/cloud.hashing.inspire.v1." + service[0], nil
		}
		return serviceName + "/cloud.hashing.inspire.v1.CloudHashingInspire", nil
	}
	return serviceName, nil
}
