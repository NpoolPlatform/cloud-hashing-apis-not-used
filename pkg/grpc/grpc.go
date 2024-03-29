package grpc

import (
	"context"
	"fmt"
	"time"

	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"

	reviewpb "github.com/NpoolPlatform/message/npool/review-service"
	reviewconst "github.com/NpoolPlatform/review-service/pkg/message/const" //nolint

	coininfopb "github.com/NpoolPlatform/message/npool/coininfo"
	coininfoconst "github.com/NpoolPlatform/sphinx-coininfo/pkg/message/const" //nolint

	sphinxproxypb "github.com/NpoolPlatform/message/npool/sphinxproxy"
	sphinxproxyconst "github.com/NpoolPlatform/sphinx-proxy/pkg/message/const" //nolint

	billingconst "github.com/NpoolPlatform/cloud-hashing-billing/pkg/message/const" //nolint
	billingpb "github.com/NpoolPlatform/message/npool/cloud-hashing-billing"

	inspireconst "github.com/NpoolPlatform/cloud-hashing-inspire/pkg/message/const" //nolint
	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"

	notificationpb "github.com/NpoolPlatform/message/npool/notification"
	notificationconst "github.com/NpoolPlatform/notification/pkg/message/const"
)

const (
	grpcTimeout = 60 * time.Second
)

func GetCoinInfo(ctx context.Context, in *coininfopb.GetCoinInfoRequest) (*coininfopb.CoinInfo, error) {
	conn, err := grpc2.GetGRPCConn(coininfoconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get coininfo connection: %v", err)
	}
	defer conn.Close()

	cli := coininfopb.NewSphinxCoinInfoClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetCoinInfo(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get coin info: %v", err)
	}

	return resp.Info, nil
}

//---------------------------------------------------------------------------------------------------------------------------

//---------------------------------------------------------------------------------------------------------------------------

func GetCouponAllocated(ctx context.Context, in *inspirepb.GetCouponAllocatedDetailRequest) (*inspirepb.CouponAllocatedDetail, error) {
	conn, err := grpc2.GetGRPCConn(inspireconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get inspire connection: %v", err)
	}
	defer conn.Close()

	cli := inspirepb.NewCloudHashingInspireClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetCouponAllocatedDetail(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get coupon allocated: %v", err)
	}

	return resp.Info, nil
}

func GetCouponsAllocatedByAppUser(ctx context.Context, in *inspirepb.GetCouponsAllocatedDetailByAppUserRequest) ([]*inspirepb.CouponAllocatedDetail, error) {
	conn, err := grpc2.GetGRPCConn(inspireconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get inspire connection: %v", err)
	}
	defer conn.Close()

	cli := inspirepb.NewCloudHashingInspireClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetCouponsAllocatedDetailByAppUser(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get coupons allocated: %v", err)
	}

	return resp.Infos, nil
}

func GetUserSpecialReduction(ctx context.Context, in *inspirepb.GetUserSpecialReductionRequest) (*inspirepb.UserSpecialReduction, error) {
	conn, err := grpc2.GetGRPCConn(inspireconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get inspire connection: %v", err)
	}
	defer conn.Close()

	cli := inspirepb.NewCloudHashingInspireClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetUserSpecialReduction(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get user special reduction: %v", err)
	}

	return resp.Info, nil
}

func GetUserSpecialReductionsByAppUser(ctx context.Context, in *inspirepb.GetUserSpecialReductionsByAppUserRequest) ([]*inspirepb.UserSpecialReduction, error) {
	conn, err := grpc2.GetGRPCConn(inspireconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get inspire connection: %v", err)
	}
	defer conn.Close()

	cli := inspirepb.NewCloudHashingInspireClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetUserSpecialReductionsByAppUser(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get user special reductions: %v", err)
	}

	return resp.Infos, nil
}

func GetUserInvitationCodeByCode(ctx context.Context, in *inspirepb.GetUserInvitationCodeByCodeRequest) (*inspirepb.UserInvitationCode, error) {
	conn, err := grpc2.GetGRPCConn(inspireconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get inspire connection: %v", err)
	}
	defer conn.Close()

	cli := inspirepb.NewCloudHashingInspireClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetUserInvitationCodeByCode(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get user invitation code: %v", err)
	}

	return resp.Info, nil
}

func GetUserInvitationCodeByAppUser(ctx context.Context, in *inspirepb.GetUserInvitationCodeByAppUserRequest) (*inspirepb.UserInvitationCode, error) {
	conn, err := grpc2.GetGRPCConn(inspireconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get inspire connection: %v", err)
	}
	defer conn.Close()

	cli := inspirepb.NewCloudHashingInspireClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetUserInvitationCodeByAppUser(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get user invitation code: %v", err)
	}

	return resp.Info, nil
}

func CreateRegistrationInvitation(ctx context.Context, in *inspirepb.CreateRegistrationInvitationRequest) (*inspirepb.RegistrationInvitation, error) {
	conn, err := grpc2.GetGRPCConn(inspireconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get inspire connection: %v", err)
	}
	defer conn.Close()

	cli := inspirepb.NewCloudHashingInspireClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.CreateRegistrationInvitation(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail create registration invitation: %v", err)
	}

	return resp.Info, nil
}

func GetRegistrationInvitationsByAppInviter(ctx context.Context, in *inspirepb.GetRegistrationInvitationsByAppInviterRequest) ([]*inspirepb.RegistrationInvitation, error) {
	conn, err := grpc2.GetGRPCConn(inspireconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get inspire connection: %v", err)
	}
	defer conn.Close()

	cli := inspirepb.NewCloudHashingInspireClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetRegistrationInvitationsByAppInviter(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get registration invitations: %v", err)
	}

	return resp.Infos, err
}

func GetRegistrationInvitationByAppInvitee(ctx context.Context, in *inspirepb.GetRegistrationInvitationByAppInviteeRequest) (*inspirepb.RegistrationInvitation, error) {
	conn, err := grpc2.GetGRPCConn(inspireconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get inspire connection: %v", err)
	}
	defer conn.Close()

	cli := inspirepb.NewCloudHashingInspireClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetRegistrationInvitationByAppInvitee(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get registration invitation: %v", err)
	}

	return resp.Info, err
}

func GetCommissionCoinSettings(ctx context.Context, in *inspirepb.GetCommissionCoinSettingsRequest) ([]*inspirepb.CommissionCoinSetting, error) {
	conn, err := grpc2.GetGRPCConn(inspireconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get inspire connection: %v", err)
	}
	defer conn.Close()

	cli := inspirepb.NewCloudHashingInspireClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetCommissionCoinSettings(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get commission coin settings: %v", err)
	}

	return resp.Infos, nil
}

func GetAppCommissionSetting(ctx context.Context, in *inspirepb.GetAppCommissionSettingRequest) (*inspirepb.AppCommissionSetting, error) {
	conn, err := grpc2.GetGRPCConn(inspireconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get inspire connection: %v", err)
	}
	defer conn.Close()

	cli := inspirepb.NewCloudHashingInspireClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetAppCommissionSetting(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get app commission setting: %v", err)
	}

	return resp.Info, nil
}

func GetAppCommissionSettingByApp(ctx context.Context, in *inspirepb.GetAppCommissionSettingByAppRequest) (*inspirepb.AppCommissionSetting, error) {
	conn, err := grpc2.GetGRPCConn(inspireconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get inspire connection: %v", err)
	}
	defer conn.Close()

	cli := inspirepb.NewCloudHashingInspireClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetAppCommissionSettingByApp(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get app commission setting: %v", err)
	}

	return resp.Info, nil
}

func GetAppPurchaseAmountSettingsByApp(ctx context.Context, in *inspirepb.GetAppPurchaseAmountSettingsByAppRequest) ([]*inspirepb.AppPurchaseAmountSetting, error) {
	conn, err := grpc2.GetGRPCConn(inspireconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get inspire connection: %v", err)
	}
	defer conn.Close()

	cli := inspirepb.NewCloudHashingInspireClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetAppPurchaseAmountSettingsByApp(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get app purchase amount setting: %v", err)
	}

	return resp.Infos, nil
}

func GetAppPurchaseAmountSettingsByAppUser(ctx context.Context, in *inspirepb.GetAppPurchaseAmountSettingsByAppUserRequest) ([]*inspirepb.AppPurchaseAmountSetting, error) {
	conn, err := grpc2.GetGRPCConn(inspireconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get inspire connection: %v", err)
	}
	defer conn.Close()

	cli := inspirepb.NewCloudHashingInspireClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetAppPurchaseAmountSettingsByAppUser(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get app purchase amount settings: %v", err)
	}

	return resp.Infos, nil
}

func GetAppInvitationSettingsByApp(ctx context.Context, in *inspirepb.GetAppInvitationSettingsByAppRequest) ([]*inspirepb.AppInvitationSetting, error) {
	conn, err := grpc2.GetGRPCConn(inspireconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get inspire connection: %v", err)
	}
	defer conn.Close()

	cli := inspirepb.NewCloudHashingInspireClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetAppInvitationSettingsByApp(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get app invitation setting: %v", err)
	}

	return resp.Infos, nil
}

//---------------------------------------------------------------------------------------------------------------------------

func CreateBillingAccount(ctx context.Context, in *billingpb.CreateCoinAccountRequest) (*billingpb.CoinAccountInfo, error) {
	conn, err := grpc2.GetGRPCConn(billingconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get billing connection: %v", err)
	}
	defer conn.Close()

	cli := billingpb.NewCloudHashingBillingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.CreateCoinAccount(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail create coin account: %v", err)
	}

	return resp.Info, nil
}

func GetBillingAccount(ctx context.Context, in *billingpb.GetCoinAccountRequest) (*billingpb.CoinAccountInfo, error) {
	conn, err := grpc2.GetGRPCConn(billingconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get billing connection: %v", err)
	}
	defer conn.Close()

	cli := billingpb.NewCloudHashingBillingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetCoinAccount(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get coin account: %v", err)
	}

	return resp.Info, nil
}

func GetIdleGoodPaymentsByGoodPaymentCoin(ctx context.Context, in *billingpb.GetIdleGoodPaymentsByGoodPaymentCoinRequest) ([]*billingpb.GoodPayment, error) {
	conn, err := grpc2.GetGRPCConn(billingconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get billing connection: %v", err)
	}
	defer conn.Close()

	cli := billingpb.NewCloudHashingBillingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetIdleGoodPaymentsByGoodPaymentCoin(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get good payments: %v", err)
	}

	return resp.Infos, nil
}

func CreateGoodPayment(ctx context.Context, in *billingpb.CreateGoodPaymentRequest) (*billingpb.GoodPayment, error) {
	conn, err := grpc2.GetGRPCConn(billingconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get billing connection: %v", err)
	}
	defer conn.Close()

	cli := billingpb.NewCloudHashingBillingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.CreateGoodPayment(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail create good payment: %v", err)
	}

	return resp.Info, nil
}

func UpdateGoodPayment(ctx context.Context, in *billingpb.UpdateGoodPaymentRequest) (*billingpb.GoodPayment, error) {
	conn, err := grpc2.GetGRPCConn(billingconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get billing connection: %v", err)
	}
	defer conn.Close()

	cli := billingpb.NewCloudHashingBillingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.UpdateGoodPayment(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail update good payment: %v", err)
	}

	return resp.Info, nil
}

func CreateUserWithdrawItem(ctx context.Context, in *billingpb.CreateUserWithdrawItemRequest) (*billingpb.UserWithdrawItem, error) {
	conn, err := grpc2.GetGRPCConn(billingconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get billing connection: %v", err)
	}
	defer conn.Close()

	cli := billingpb.NewCloudHashingBillingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.CreateUserWithdrawItem(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail create user withdraw item: %v", err)
	}

	return resp.Info, nil
}

func CreateUserWithdraw(ctx context.Context, in *billingpb.CreateUserWithdrawRequest) (*billingpb.UserWithdraw, error) {
	conn, err := grpc2.GetGRPCConn(billingconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get billing connection: %v", err)
	}
	defer conn.Close()

	cli := billingpb.NewCloudHashingBillingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.CreateUserWithdraw(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail create user withdraw: %v", err)
	}

	return resp.Info, nil
}

func UpdateUserWithdrawItem(ctx context.Context, in *billingpb.UpdateUserWithdrawItemRequest) (*billingpb.UserWithdrawItem, error) {
	conn, err := grpc2.GetGRPCConn(billingconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get billing connection: %v", err)
	}
	defer conn.Close()

	cli := billingpb.NewCloudHashingBillingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.UpdateUserWithdrawItem(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail update user withdraw item: %v", err)
	}

	return resp.Info, nil
}

func GetUserWithdraw(ctx context.Context, in *billingpb.GetUserWithdrawRequest) (*billingpb.UserWithdraw, error) {
	conn, err := grpc2.GetGRPCConn(billingconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get billing connection: %v", err)
	}
	defer conn.Close()

	cli := billingpb.NewCloudHashingBillingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetUserWithdraw(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get user withdraw: %v", err)
	}

	return resp.Info, nil
}

func DeleteUserWithdraw(ctx context.Context, in *billingpb.DeleteUserWithdrawRequest) (*billingpb.UserWithdraw, error) {
	conn, err := grpc2.GetGRPCConn(billingconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get billing connection: %v", err)
	}
	defer conn.Close()

	cli := billingpb.NewCloudHashingBillingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.DeleteUserWithdraw(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get user withdraw: %v", err)
	}

	return resp.Info, nil
}

func GetUserWithdrawItem(ctx context.Context, in *billingpb.GetUserWithdrawItemRequest) (*billingpb.UserWithdrawItem, error) {
	conn, err := grpc2.GetGRPCConn(billingconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get billing connection: %v", err)
	}
	defer conn.Close()

	cli := billingpb.NewCloudHashingBillingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetUserWithdrawItem(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get user withdraw item: %v", err)
	}

	return resp.Info, nil
}

func GetUserWithdrawItemsByAppUser(ctx context.Context, in *billingpb.GetUserWithdrawItemsByAppUserRequest) ([]*billingpb.UserWithdrawItem, error) {
	conn, err := grpc2.GetGRPCConn(billingconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get billing connection: %v", err)
	}
	defer conn.Close()

	cli := billingpb.NewCloudHashingBillingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetUserWithdrawItemsByAppUser(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get user withdraw items: %v", err)
	}

	return resp.Infos, nil
}

func GetUserWithdrawItemsByAppUserCoinWithdrawType(ctx context.Context, in *billingpb.GetUserWithdrawItemsByAppUserCoinWithdrawTypeRequest) ([]*billingpb.UserWithdrawItem, error) {
	conn, err := grpc2.GetGRPCConn(billingconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get billing connection: %v", err)
	}
	defer conn.Close()

	cli := billingpb.NewCloudHashingBillingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetUserWithdrawItemsByAppUserCoinWithdrawType(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get user withdraw items: %v", err)
	}

	return resp.Infos, nil
}

func GetUserWithdrawItemsByAppUserWithdrawType(ctx context.Context, in *billingpb.GetUserWithdrawItemsByAppUserWithdrawTypeRequest) ([]*billingpb.UserWithdrawItem, error) {
	conn, err := grpc2.GetGRPCConn(billingconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get billing connection: %v", err)
	}
	defer conn.Close()

	cli := billingpb.NewCloudHashingBillingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetUserWithdrawItemsByAppUserWithdrawType(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get user withdraw item: %v", err)
	}

	return resp.Infos, nil
}

func GetUserBenefitsByAppUserCoin(ctx context.Context, in *billingpb.GetUserBenefitsByAppUserCoinRequest) ([]*billingpb.UserBenefit, error) {
	conn, err := grpc2.GetGRPCConn(billingconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get billing connection: %v", err)
	}
	defer conn.Close()

	cli := billingpb.NewCloudHashingBillingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetUserBenefitsByAppUserCoin(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get user benefit: %v", err)
	}

	return resp.Infos, nil
}

func GetAppWithdrawSettingByAppCoin(ctx context.Context, in *billingpb.GetAppWithdrawSettingByAppCoinRequest) (*billingpb.AppWithdrawSetting, error) {
	conn, err := grpc2.GetGRPCConn(billingconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get billing connection: %v", err)
	}
	defer conn.Close()

	cli := billingpb.NewCloudHashingBillingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetAppWithdrawSettingByAppCoin(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get app withdraw setting: %v", err)
	}

	return resp.Info, nil
}

func GetPlatformSetting(ctx context.Context, in *billingpb.GetPlatformSettingRequest) (*billingpb.PlatformSetting, error) {
	conn, err := grpc2.GetGRPCConn(billingconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get billing connection: %v", err)
	}
	defer conn.Close()

	cli := billingpb.NewCloudHashingBillingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetPlatformSetting(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get platform setting: %v", err)
	}

	return resp.Info, nil
}

func GetUserWithdrawByAccount(ctx context.Context, in *billingpb.GetUserWithdrawByAccountRequest) (*billingpb.UserWithdraw, error) {
	conn, err := grpc2.GetGRPCConn(billingconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get billing connection: %v", err)
	}
	defer conn.Close()

	cli := billingpb.NewCloudHashingBillingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetUserWithdrawByAccount(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get user withdraw: %v", err)
	}

	return resp.Info, nil
}

func GetUserWithdrawsByAppUser(ctx context.Context, in *billingpb.GetUserWithdrawsByAppUserRequest) ([]*billingpb.UserWithdraw, error) {
	conn, err := grpc2.GetGRPCConn(billingconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get billing connection: %v", err)
	}
	defer conn.Close()

	cli := billingpb.NewCloudHashingBillingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetUserWithdrawsByAppUser(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get user withdraw: %v", err)
	}

	return resp.Infos, nil
}

func GetUserWithdrawsByAppUserCoin(ctx context.Context, in *billingpb.GetUserWithdrawsByAppUserCoinRequest) ([]*billingpb.UserWithdraw, error) {
	conn, err := grpc2.GetGRPCConn(billingconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get billing connection: %v", err)
	}
	defer conn.Close()

	cli := billingpb.NewCloudHashingBillingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetUserWithdrawsByAppUserCoin(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get user withdraw: %v", err)
	}

	return resp.Infos, nil
}

func GetCoinSettingByCoin(ctx context.Context, in *billingpb.GetCoinSettingByCoinRequest) (*billingpb.CoinSetting, error) {
	conn, err := grpc2.GetGRPCConn(billingconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get billing connection: %v", err)
	}
	defer conn.Close()

	cli := billingpb.NewCloudHashingBillingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetCoinSettingByCoin(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get coin setting: %v", err)
	}

	return resp.Info, nil
}

func CreateCoinAccountTransaction(ctx context.Context, in *billingpb.CreateCoinAccountTransactionRequest) (*billingpb.CoinAccountTransaction, error) {
	conn, err := grpc2.GetGRPCConn(billingconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get billing connection: %v", err)
	}
	defer conn.Close()

	cli := billingpb.NewCloudHashingBillingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.CreateCoinAccountTransaction(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail create coin account transaction: %v", err)
	}

	return resp.Info, nil
}

func GetCoinAccountTransactionsByAppUserCoin(ctx context.Context, in *billingpb.GetCoinAccountTransactionsByAppUserCoinRequest) ([]*billingpb.CoinAccountTransaction, error) {
	conn, err := grpc2.GetGRPCConn(billingconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get billing connection: %v", err)
	}
	defer conn.Close()

	cli := billingpb.NewCloudHashingBillingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetCoinAccountTransactionsByAppUserCoin(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get coin account transactions: %v", err)
	}

	return resp.Infos, nil
}

//---------------------------------------------------------------------------------------------------------------------------

func CreateCoinAddress(ctx context.Context, in *sphinxproxypb.CreateWalletRequest) (*sphinxproxypb.WalletInfo, error) {
	conn, err := grpc2.GetGRPCConn(sphinxproxyconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get sphinxproxy connection: %v", err)
	}
	defer conn.Close()

	cli := sphinxproxypb.NewSphinxProxyClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.CreateWallet(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail create wallet: %v", err)
	}

	return resp.Info, nil
}

func GetBalance(ctx context.Context, in *sphinxproxypb.GetBalanceRequest) (*sphinxproxypb.BalanceInfo, error) {
	conn, err := grpc2.GetGRPCConn(sphinxproxyconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get sphinxproxy connection: %v", err)
	}
	defer conn.Close()

	cli := sphinxproxypb.NewSphinxProxyClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetBalance(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get balance info: %v", err)
	}

	return resp.Info, nil
}

//---------------------------------------------------------------------------------------------------------------------------

func CreateNotification(ctx context.Context, in *notificationpb.CreateNotificationRequest) (*notificationpb.UserNotification, error) {
	conn, err := grpc2.GetGRPCConn(notificationconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get notification connection: %v", err)
	}
	defer conn.Close()

	cli := notificationpb.NewNotificationClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.CreateNotification(ctx, in)
	if err != nil {
		return nil, err
	}
	return resp.Info, err
}

//---------------------------------------------------------------------------------------------------------------------------

func CreateReview(ctx context.Context, in *reviewpb.CreateReviewRequest) (*reviewpb.Review, error) {
	conn, err := grpc2.GetGRPCConn(reviewconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get review connection: %v", err)
	}
	defer conn.Close()

	cli := reviewpb.NewReviewServiceClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.CreateReview(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail create review: %v", err)
	}

	return resp.Info, nil
}

func UpdateReview(ctx context.Context, in *reviewpb.UpdateReviewRequest) (*reviewpb.Review, error) {
	conn, err := grpc2.GetGRPCConn(reviewconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get review connection: %v", err)
	}
	defer conn.Close()

	cli := reviewpb.NewReviewServiceClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.UpdateReview(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail update review: %v", err)
	}

	return resp.Info, nil
}

func GetReview(ctx context.Context, in *reviewpb.GetReviewRequest) (*reviewpb.Review, error) {
	conn, err := grpc2.GetGRPCConn(reviewconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get review connection: %v", err)
	}
	defer conn.Close()

	cli := reviewpb.NewReviewServiceClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetReview(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get review: %v", err)
	}

	return resp.Info, nil
}

func GetReviewsByAppDomainObjectTypeID(ctx context.Context, in *reviewpb.GetReviewsByAppDomainObjectTypeIDRequest) ([]*reviewpb.Review, error) {
	conn, err := grpc2.GetGRPCConn(reviewconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get review connection: %v", err)
	}
	defer conn.Close()

	cli := reviewpb.NewReviewServiceClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetReviewsByAppDomainObjectTypeID(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get reviews: %v", err)
	}

	return resp.Infos, nil
}

func GetReviewsByAppDomain(ctx context.Context, in *reviewpb.GetReviewsByAppDomainRequest) ([]*reviewpb.Review, error) {
	conn, err := grpc2.GetGRPCConn(reviewconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get review connection: %v", err)
	}
	defer conn.Close()

	cli := reviewpb.NewReviewServiceClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetReviewsByAppDomain(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get reviews: %v", err)
	}

	return resp.Infos, nil
}

//---------------------------------------------------------------------------------------------------------------------------
