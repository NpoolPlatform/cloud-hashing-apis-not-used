package grpc

import (
	"context"
	"fmt"
	"time"

	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"

	goodsconst "github.com/NpoolPlatform/cloud-hashing-goods/pkg/message/const" //nolint
	goodspb "github.com/NpoolPlatform/message/npool/cloud-hashing-goods"

	reviewpb "github.com/NpoolPlatform/message/npool/review-service"
	reviewconst "github.com/NpoolPlatform/review-service/pkg/message/const" //nolint

	coininfopb "github.com/NpoolPlatform/message/npool/coininfo"
	coininfoconst "github.com/NpoolPlatform/sphinx-coininfo/pkg/message/const" //nolint

	sphinxproxypb "github.com/NpoolPlatform/message/npool/sphinxproxy"
	sphinxproxyconst "github.com/NpoolPlatform/sphinx-proxy/pkg/message/const" //nolint

	orderconst "github.com/NpoolPlatform/cloud-hashing-order/pkg/message/const" //nolint
	orderpb "github.com/NpoolPlatform/message/npool/cloud-hashing-order"

	billingconst "github.com/NpoolPlatform/cloud-hashing-billing/pkg/message/const" //nolint
	billingpb "github.com/NpoolPlatform/message/npool/cloud-hashing-billing"

	inspireconst "github.com/NpoolPlatform/cloud-hashing-inspire/pkg/message/const" //nolint
	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"

	appusermgrconst "github.com/NpoolPlatform/appuser-manager/pkg/message/const" //nolint
	appusermgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v1"

	kycmgrconst "github.com/NpoolPlatform/kyc-management/pkg/message/const" //nolint
	kycmgrpb "github.com/NpoolPlatform/message/npool/kyc"

	thirdgwpb "github.com/NpoolPlatform/message/npool/thirdgateway"
	thirdgwconst "github.com/NpoolPlatform/third-gateway/pkg/message/const"

	logingwconst "github.com/NpoolPlatform/login-gateway/pkg/message/const"
	logingwpb "github.com/NpoolPlatform/message/npool/logingateway"

	notificationpb "github.com/NpoolPlatform/message/npool/notification"
	notificationconst "github.com/NpoolPlatform/notification/pkg/message/const"
)

const (
	grpcTimeout = 60 * time.Second
)

func CreateGood(ctx context.Context, in *goodspb.CreateGoodRequest) (*goodspb.GoodInfo, error) {
	conn, err := grpc2.GetGRPCConn(goodsconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get goods connection: %v", err)
	}
	defer conn.Close()

	cli := goodspb.NewCloudHashingGoodsClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.CreateGood(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail create good: %v", err)
	}

	return resp.Info, nil
}

func GetGood(ctx context.Context, in *goodspb.GetGoodRequest) (*goodspb.GoodInfo, error) {
	conn, err := grpc2.GetGRPCConn(goodsconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get goods connection: %v", err)
	}
	defer conn.Close()

	cli := goodspb.NewCloudHashingGoodsClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetGood(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail create good: %v", err)
	}

	return resp.Info, nil
}

func GetGoodsDetail(ctx context.Context, in *goodspb.GetGoodsDetailRequest) ([]*goodspb.GoodDetail, error) {
	conn, err := grpc2.GetGRPCConn(goodsconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get goods connection: %v", err)
	}
	defer conn.Close()

	cli := goodspb.NewCloudHashingGoodsClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetGoodsDetail(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get goods: %v", err)
	}

	return resp.Infos, nil
}

func GetGoodsDetailByApp(ctx context.Context, in *goodspb.GetGoodsDetailByAppRequest) ([]*goodspb.GoodDetail, error) {
	conn, err := grpc2.GetGRPCConn(goodsconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get goods connection: %v", err)
	}
	defer conn.Close()

	cli := goodspb.NewCloudHashingGoodsClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetGoodsDetailByApp(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get goods: %v", err)
	}

	return resp.Infos, nil
}

func GetGoodDetail(ctx context.Context, in *goodspb.GetGoodDetailRequest) (*goodspb.GoodDetail, error) {
	conn, err := grpc2.GetGRPCConn(goodsconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get goods connection: %v", err)
	}
	defer conn.Close()

	cli := goodspb.NewCloudHashingGoodsClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetGoodDetail(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get good: %v", err)
	}

	return resp.Info, nil
}

func GetRecommendGoodsByApp(ctx context.Context, in *goodspb.GetRecommendGoodsByAppRequest) ([]*goodspb.RecommendGood, error) {
	conn, err := grpc2.GetGRPCConn(goodsconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get goods connection: %v", err)
	}
	defer conn.Close()

	cli := goodspb.NewCloudHashingGoodsClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetRecommendGoodsByApp(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get goods: %v", err)
	}

	return resp.Infos, nil
}

func GetAppGoodByAppGood(ctx context.Context, in *goodspb.GetAppGoodByAppGoodRequest) (*goodspb.AppGoodInfo, error) {
	conn, err := grpc2.GetGRPCConn(goodsconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get goods connection: %v", err)
	}
	defer conn.Close()

	cli := goodspb.NewCloudHashingGoodsClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetAppGoodByAppGood(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get app good: %v", err)
	}

	return resp.Info, nil
}

func GetAppGoodPromotionByAppGoodTimestamp(ctx context.Context, in *goodspb.GetAppGoodPromotionByAppGoodTimestampRequest) (*goodspb.AppGoodPromotion, error) {
	conn, err := grpc2.GetGRPCConn(goodsconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get goods connection: %v", err)
	}
	defer conn.Close()

	cli := goodspb.NewCloudHashingGoodsClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetAppGoodPromotionByAppGoodTimestamp(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get app good: %v", err)
	}

	return resp.Info, nil
}

func GetAppGoodPromotion(ctx context.Context, in *goodspb.GetAppGoodPromotionRequest) (*goodspb.AppGoodPromotion, error) {
	conn, err := grpc2.GetGRPCConn(goodsconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get goods connection: %v", err)
	}
	defer conn.Close()

	cli := goodspb.NewCloudHashingGoodsClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetAppGoodPromotion(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get app good: %v", err)
	}

	return resp.Info, nil
}

//---------------------------------------------------------------------------------------------------------------------------

func GetCoinInfos(ctx context.Context, in *coininfopb.GetCoinInfosRequest) ([]*coininfopb.CoinInfo, error) {
	conn, err := grpc2.GetGRPCConn(coininfoconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get coininfo connection: %v", err)
	}
	defer conn.Close()

	cli := coininfopb.NewSphinxCoinInfoClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetCoinInfos(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get coin infos: %v", err)
	}

	return resp.Infos, nil
}

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

func GetOrder(ctx context.Context, in *orderpb.GetOrderRequest) (*orderpb.Order, error) {
	conn, err := grpc2.GetGRPCConn(orderconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get order connection: %v", err)
	}
	defer conn.Close()

	cli := orderpb.NewCloudHashingOrderClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetOrder(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get order: %v", err)
	}

	return resp.Info, nil
}

func GetOrderByAppUserCouponTypeID(ctx context.Context, in *orderpb.GetOrderByAppUserCouponTypeIDRequest) (*orderpb.Order, error) {
	conn, err := grpc2.GetGRPCConn(orderconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get order connection: %v", err)
	}
	defer conn.Close()

	cli := orderpb.NewCloudHashingOrderClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetOrderByAppUserCouponTypeID(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get order: %v", err)
	}

	return resp.Info, nil
}

func GetOrderDetail(ctx context.Context, in *orderpb.GetOrderDetailRequest) (*orderpb.OrderDetail, error) {
	conn, err := grpc2.GetGRPCConn(orderconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get order connection: %v", err)
	}
	defer conn.Close()

	cli := orderpb.NewCloudHashingOrderClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetOrderDetail(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get order: %v", err)
	}

	return resp.Info, nil
}

func GetOrdersDetailByAppUser(ctx context.Context, in *orderpb.GetOrdersDetailByAppUserRequest) ([]*orderpb.OrderDetail, error) {
	conn, err := grpc2.GetGRPCConn(orderconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get order connection: %v", err)
	}
	defer conn.Close()

	cli := orderpb.NewCloudHashingOrderClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetOrdersDetailByAppUser(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get orders: %v", err)
	}

	return resp.Infos, nil
}

func GetOrdersShortDetailByAppUser(ctx context.Context, in *orderpb.GetOrdersShortDetailByAppUserRequest) ([]*orderpb.OrderDetail, error) {
	conn, err := grpc2.GetGRPCConn(orderconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get order connection: %v", err)
	}
	defer conn.Close()

	cli := orderpb.NewCloudHashingOrderClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetOrdersShortDetailByAppUser(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get order short detail")
	}

	return resp.Infos, nil
}

func GetOrdersDetailByApp(ctx context.Context, in *orderpb.GetOrdersDetailByAppRequest) ([]*orderpb.OrderDetail, error) {
	conn, err := grpc2.GetGRPCConn(orderconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get order connection: %v", err)
	}
	defer conn.Close()

	cli := orderpb.NewCloudHashingOrderClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetOrdersDetailByApp(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get orders detail: %v", err)
	}

	return resp.Infos, nil
}

func GetOrdersDetailByGood(ctx context.Context, in *orderpb.GetOrdersDetailByGoodRequest) ([]*orderpb.OrderDetail, error) {
	conn, err := grpc2.GetGRPCConn(orderconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get order connection: %v", err)
	}
	defer conn.Close()

	cli := orderpb.NewCloudHashingOrderClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetOrdersDetailByGood(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get orders detail: %v", err)
	}

	return resp.Infos, nil
}

func GetPaymentsByAppUserState(ctx context.Context, in *orderpb.GetPaymentsByAppUserStateRequest) ([]*orderpb.Payment, error) {
	conn, err := grpc2.GetGRPCConn(orderconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get order connection: %v", err)
	}
	defer conn.Close()

	cli := orderpb.NewCloudHashingOrderClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetPaymentsByAppUserState(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get payments detail: %v", err)
	}

	return resp.Infos, nil
}

func CreateOrder(ctx context.Context, in *orderpb.CreateOrderRequest) (*orderpb.Order, error) {
	conn, err := grpc2.GetGRPCConn(orderconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get order connection: %v", err)
	}
	defer conn.Close()

	cli := orderpb.NewCloudHashingOrderClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.CreateOrder(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail create order: %v", err)
	}

	return resp.Info, nil
}

func CreateGoodPaying(ctx context.Context, in *orderpb.CreateGoodPayingRequest) (*orderpb.GoodPaying, error) {
	conn, err := grpc2.GetGRPCConn(orderconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get order connection: %v", err)
	}
	defer conn.Close()

	cli := orderpb.NewCloudHashingOrderClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.CreateGoodPaying(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail create good paying: %v", err)
	}

	return resp.Info, nil
}

func CreateGasPaying(ctx context.Context, in *orderpb.CreateGasPayingRequest) (*orderpb.GasPaying, error) {
	conn, err := grpc2.GetGRPCConn(orderconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get order connection: %v", err)
	}
	defer conn.Close()

	cli := orderpb.NewCloudHashingOrderClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.CreateGasPaying(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail create gas paying: %v", err)
	}

	return resp.Info, nil
}

func CreatePayment(ctx context.Context, in *orderpb.CreatePaymentRequest) (*orderpb.Payment, error) {
	conn, err := grpc2.GetGRPCConn(orderconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get order connection: %v", err)
	}
	defer conn.Close()

	cli := orderpb.NewCloudHashingOrderClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.CreatePayment(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail create payment: %v", err)
	}

	return resp.Info, nil
}

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

func Signup(ctx context.Context, in *appusermgrpb.CreateAppUserWithSecretRequest) (*appusermgrpb.AppUser, error) {
	conn, err := grpc2.GetGRPCConn(appusermgrconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get app user manager connection: %v", err)
	}
	defer conn.Close()

	cli := appusermgrpb.NewAppUserManagerClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.CreateAppUserWithSecret(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail create app user: %v", err)
	}
	return resp.Info, nil
}

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

func VerifyAppUserByAppAccountPassword(ctx context.Context, in *appusermgrpb.VerifyAppUserByAppAccountPasswordRequest) (*appusermgrpb.AppUserInfo, error) {
	conn, err := grpc2.GetGRPCConn(appusermgrconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get app user manager connection: %v", err)
	}
	defer conn.Close()

	cli := appusermgrpb.NewAppUserManagerClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.VerifyAppUserByAppAccountPassword(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail verify app user: %v", err)
	}

	return resp.Info, nil
}

func GetAppUserByAppUser(ctx context.Context, in *appusermgrpb.GetAppUserByAppUserRequest) (*appusermgrpb.AppUser, error) {
	conn, err := grpc2.GetGRPCConn(appusermgrconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get app user manager connection: %v", err)
	}
	defer conn.Close()

	cli := appusermgrpb.NewAppUserManagerClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetAppUserByAppUser(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get app user: %v", err)
	}

	return resp.Info, nil
}

func GetAppUserExtraByAppUser(ctx context.Context, in *appusermgrpb.GetAppUserExtraByAppUserRequest) (*appusermgrpb.AppUserExtra, error) {
	conn, err := grpc2.GetGRPCConn(appusermgrconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get app user manager connection: %v", err)
	}
	defer conn.Close()

	cli := appusermgrpb.NewAppUserManagerClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetAppUserExtraByAppUser(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get app user extra: %v", err)
	}

	return resp.Info, nil
}

func UpdateAppUser(ctx context.Context, in *appusermgrpb.UpdateAppUserRequest) (*appusermgrpb.AppUser, error) {
	conn, err := grpc2.GetGRPCConn(appusermgrconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get app user manager connection: %v", err)
	}
	defer conn.Close()

	cli := appusermgrpb.NewAppUserManagerClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.UpdateAppUser(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail update app user: %v", err)
	}

	return resp.Info, nil
}

func UpdateAppUserExtra(ctx context.Context, in *appusermgrpb.UpdateAppUserExtraRequest) (*appusermgrpb.AppUserExtra, error) {
	conn, err := grpc2.GetGRPCConn(appusermgrconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get app user manager connection: %v", err)
	}
	defer conn.Close()

	cli := appusermgrpb.NewAppUserManagerClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.UpdateAppUserExtra(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail update app user extra: %v", err)
	}

	return resp.Info, nil
}

func CreateAppUserExtra(ctx context.Context, in *appusermgrpb.CreateAppUserExtraRequest) (*appusermgrpb.AppUserExtra, error) {
	conn, err := grpc2.GetGRPCConn(appusermgrconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get app user manager connection: %v", err)
	}
	defer conn.Close()

	cli := appusermgrpb.NewAppUserManagerClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.CreateAppUserExtra(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail create app user extra: %v", err)
	}

	return resp.Info, nil
}

func UpdateAppUserControl(ctx context.Context, in *appusermgrpb.UpdateAppUserControlRequest) (*appusermgrpb.AppUserControl, error) {
	conn, err := grpc2.GetGRPCConn(appusermgrconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get app user manager connection: %v", err)
	}
	defer conn.Close()

	cli := appusermgrpb.NewAppUserManagerClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.UpdateAppUserControl(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail update app user control: %v", err)
	}

	return resp.Info, nil
}

func CreateAppUserControl(ctx context.Context, in *appusermgrpb.CreateAppUserControlRequest) (*appusermgrpb.AppUserControl, error) {
	conn, err := grpc2.GetGRPCConn(appusermgrconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get app user manager connection: %v", err)
	}
	defer conn.Close()

	cli := appusermgrpb.NewAppUserManagerClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.CreateAppUserControl(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail create app user control: %v", err)
	}

	return resp.Info, nil
}

func GetAppUserByAppAccount(ctx context.Context, in *appusermgrpb.GetAppUserByAppAccountRequest) (*appusermgrpb.AppUser, error) {
	conn, err := grpc2.GetGRPCConn(appusermgrconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get app user manager connection: %v", err)
	}
	defer conn.Close()

	cli := appusermgrpb.NewAppUserManagerClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetAppUserByAppAccount(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get app user: %v", err)
	}

	return resp.Info, nil
}

func GetAppUserInfoByAppUser(ctx context.Context, in *appusermgrpb.GetAppUserInfoByAppUserRequest) (*appusermgrpb.AppUserInfo, error) {
	conn, err := grpc2.GetGRPCConn(appusermgrconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get app user manager connection: %v", err)
	}
	defer conn.Close()

	cli := appusermgrpb.NewAppUserManagerClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetAppUserInfoByAppUser(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get app user info: %v", err)
	}

	return resp.Info, nil
}

func GetAppUserSecretByAppUser(ctx context.Context, in *appusermgrpb.GetAppUserSecretByAppUserRequest) (*appusermgrpb.AppUserSecret, error) {
	conn, err := grpc2.GetGRPCConn(appusermgrconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get app user manager connection: %v", err)
	}
	defer conn.Close()

	cli := appusermgrpb.NewAppUserManagerClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetAppUserSecretByAppUser(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get app user secret: %v", err)
	}

	return resp.Info, nil
}

func UpdateAppUserSecret(ctx context.Context, in *appusermgrpb.UpdateAppUserSecretRequest) (*appusermgrpb.AppUserSecret, error) {
	conn, err := grpc2.GetGRPCConn(appusermgrconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get app user manager connection: %v", err)
	}
	defer conn.Close()

	cli := appusermgrpb.NewAppUserManagerClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.UpdateAppUserSecret(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail update app user secret: %v", err)
	}

	return resp.Info, nil
}

//---------------------------------------------------------------------------------------------------------------------------

func GetApp(ctx context.Context, in *appusermgrpb.GetAppInfoRequest) (*appusermgrpb.AppInfo, error) {
	conn, err := grpc2.GetGRPCConn(appusermgrconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get app user manager connection: %v", err)
	}
	defer conn.Close()

	cli := appusermgrpb.NewAppUserManagerClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetAppInfo(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get app info: %v", err)
	}

	return resp.Info, nil
}

//---------------------------------------------------------------------------------------------------------------------------

func GetReviewsByDomain(ctx context.Context, in *reviewpb.GetReviewsByDomainRequest) ([]*reviewpb.Review, error) {
	conn, err := grpc2.GetGRPCConn(reviewconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get review connection: %v", err)
	}
	defer conn.Close()

	cli := reviewpb.NewReviewServiceClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetReviewsByDomain(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get reviews: %v", err)
	}

	return resp.Infos, nil
}

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

func GetKycByIDs(ctx context.Context, in *kycmgrpb.GetKycByKycIDsRequest) ([]*kycmgrpb.KycInfo, error) {
	conn, err := grpc2.GetGRPCConn(kycmgrconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get kyc connection: %v", err)
	}
	defer conn.Close()

	cli := kycmgrpb.NewKycManagementClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetKycByKycIDs(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get kycs: %v", err)
	}

	return resp.Infos, nil
}

func GetKycByUserID(ctx context.Context, in *kycmgrpb.GetKycByUserIDRequest) (*kycmgrpb.KycInfo, error) {
	conn, err := grpc2.GetGRPCConn(kycmgrconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get kyc connection: %v", err)
	}
	defer conn.Close()

	cli := kycmgrpb.NewKycManagementClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.GetKycByUserID(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail get kyc: %v", err)
	}

	return resp.Info, nil
}

func CreateKyc(ctx context.Context, in *kycmgrpb.CreateKycRequest) (*kycmgrpb.KycInfo, error) {
	conn, err := grpc2.GetGRPCConn(kycmgrconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get kyc connection: %v", err)
	}
	defer conn.Close()

	cli := kycmgrpb.NewKycManagementClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.CreateKyc(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail create kyc: %v", err)
	}

	return resp.Info, nil
}

func UpdateKyc(ctx context.Context, in *kycmgrpb.UpdateKycRequest) (*kycmgrpb.KycInfo, error) {
	conn, err := grpc2.GetGRPCConn(kycmgrconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get kyc connection: %v", err)
	}
	defer conn.Close()

	cli := kycmgrpb.NewKycManagementClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.UpdateKyc(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail create kyc: %v", err)
	}

	return resp.Info, nil
}

//---------------------------------------------------------------------------------------------------------------------------

func VerifyEmailCode(ctx context.Context, in *thirdgwpb.VerifyEmailCodeRequest) (*thirdgwpb.VerifyEmailCodeResponse, error) {
	conn, err := grpc2.GetGRPCConn(thirdgwconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get third gateway connection: %v", err)
	}
	defer conn.Close()

	cli := thirdgwpb.NewThirdGatewayClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.VerifyEmailCode(ctx, in)
}

func VerifySMSCode(ctx context.Context, in *thirdgwpb.VerifySMSCodeRequest) (*thirdgwpb.VerifySMSCodeResponse, error) {
	conn, err := grpc2.GetGRPCConn(thirdgwconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get third gateway connection: %v", err)
	}
	defer conn.Close()

	cli := thirdgwpb.NewThirdGatewayClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.VerifySMSCode(ctx, in)
}

func VerifyGoogleAuthentication(ctx context.Context, in *thirdgwpb.VerifyGoogleAuthenticationRequest) (*thirdgwpb.VerifyGoogleAuthenticationResponse, error) {
	conn, err := grpc2.GetGRPCConn(thirdgwconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get third gateway connection: %v", err)
	}
	defer conn.Close()

	cli := thirdgwpb.NewThirdGatewayClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.VerifyGoogleAuthentication(ctx, in)
}

//---------------------------------------------------------------------------------------------------------------------------

func UpdateCache(ctx context.Context, in *logingwpb.UpdateCacheRequest) (*appusermgrpb.AppUserInfo, error) {
	conn, err := grpc2.GetGRPCConn(logingwconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get login gateway connection: %v", err)
	}
	defer conn.Close()

	cli := logingwpb.NewLoginGatewayClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	resp, err := cli.UpdateCache(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("fail update cache: %v", err)
	}

	return resp.Info, nil
}
