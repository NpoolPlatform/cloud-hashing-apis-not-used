package grpc

import (
	"context"
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

	usermgrpb "github.com/NpoolPlatform/message/npool/user"
	usermgrconst "github.com/NpoolPlatform/user-management/pkg/message/const" //nolint

	appmgrconst "github.com/NpoolPlatform/application-management/pkg/message/const" //nolint
	appmgrpb "github.com/NpoolPlatform/message/npool/application"

	"golang.org/x/xerrors"
)

const (
	grpcTimeout = 5 * time.Second
)

func GetGood(ctx context.Context, in *goodspb.GetGoodRequest) (*goodspb.GetGoodResponse, error) {
	conn, err := grpc2.GetGRPCConn(goodsconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get goods connection: %v", err)
	}
	defer conn.Close()

	cli := goodspb.NewCloudHashingGoodsClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.GetGood(ctx, in)
}

func GetGoodsDetail(ctx context.Context, in *goodspb.GetGoodsDetailRequest) (*goodspb.GetGoodsDetailResponse, error) {
	conn, err := grpc2.GetGRPCConn(goodsconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get goods connection: %v", err)
	}
	defer conn.Close()

	cli := goodspb.NewCloudHashingGoodsClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.GetGoodsDetail(ctx, in)
}

func GetGoodDetail(ctx context.Context, in *goodspb.GetGoodDetailRequest) (*goodspb.GetGoodDetailResponse, error) {
	conn, err := grpc2.GetGRPCConn(goodsconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get goods connection: %v", err)
	}
	defer conn.Close()

	cli := goodspb.NewCloudHashingGoodsClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.GetGoodDetail(ctx, in)
}

func GetRecommendGoodsByApp(ctx context.Context, in *goodspb.GetRecommendGoodsByAppRequest) (*goodspb.GetRecommendGoodsByAppResponse, error) {
	conn, err := grpc2.GetGRPCConn(goodsconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get goods connection: %v", err)
	}
	defer conn.Close()

	cli := goodspb.NewCloudHashingGoodsClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.GetRecommendGoodsByApp(ctx, in)
}

//---------------------------------------------------------------------------------------------------------------------------

func GetCoinInfos(ctx context.Context, in *coininfopb.GetCoinInfosRequest) (*coininfopb.GetCoinInfosResponse, error) {
	conn, err := grpc2.GetGRPCConn(coininfoconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get coininfo connection: %v", err)
	}
	defer conn.Close()

	cli := coininfopb.NewSphinxCoinInfoClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.GetCoinInfos(ctx, in)
}

func GetCoinInfo(ctx context.Context, in *coininfopb.GetCoinInfoRequest) (*coininfopb.GetCoinInfoResponse, error) {
	conn, err := grpc2.GetGRPCConn(coininfoconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get coininfo connection: %v", err)
	}
	defer conn.Close()

	cli := coininfopb.NewSphinxCoinInfoClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.GetCoinInfo(ctx, in)
}

//---------------------------------------------------------------------------------------------------------------------------

func GetOrder(ctx context.Context, in *orderpb.GetOrderRequest) (*orderpb.GetOrderResponse, error) {
	conn, err := grpc2.GetGRPCConn(orderconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get order connection: %v", err)
	}
	defer conn.Close()

	cli := orderpb.NewCloudHashingOrderClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.GetOrder(ctx, in)
}

func GetOrderDetail(ctx context.Context, in *orderpb.GetOrderDetailRequest) (*orderpb.GetOrderDetailResponse, error) {
	conn, err := grpc2.GetGRPCConn(orderconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get order connection: %v", err)
	}
	defer conn.Close()

	cli := orderpb.NewCloudHashingOrderClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.GetOrderDetail(ctx, in)
}

func GetOrdersDetailByAppUser(ctx context.Context, in *orderpb.GetOrdersDetailByAppUserRequest) (*orderpb.GetOrdersDetailByAppUserResponse, error) {
	conn, err := grpc2.GetGRPCConn(orderconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get order connection: %v", err)
	}
	defer conn.Close()

	cli := orderpb.NewCloudHashingOrderClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.GetOrdersDetailByAppUser(ctx, in)
}

func GetOrdersShortDetailByAppUser(ctx context.Context, in *orderpb.GetOrdersShortDetailByAppUserRequest) (*orderpb.GetOrdersShortDetailByAppUserResponse, error) {
	conn, err := grpc2.GetGRPCConn(orderconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get order connection: %v", err)
	}
	defer conn.Close()

	cli := orderpb.NewCloudHashingOrderClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.GetOrdersShortDetailByAppUser(ctx, in)
}

func GetOrdersDetailByApp(ctx context.Context, in *orderpb.GetOrdersDetailByAppRequest) (*orderpb.GetOrdersDetailByAppResponse, error) {
	conn, err := grpc2.GetGRPCConn(orderconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get order connection: %v", err)
	}
	defer conn.Close()

	cli := orderpb.NewCloudHashingOrderClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.GetOrdersDetailByApp(ctx, in)
}

func GetOrdersDetailByGood(ctx context.Context, in *orderpb.GetOrdersDetailByGoodRequest) (*orderpb.GetOrdersDetailByGoodResponse, error) {
	conn, err := grpc2.GetGRPCConn(orderconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get order connection: %v", err)
	}
	defer conn.Close()

	cli := orderpb.NewCloudHashingOrderClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.GetOrdersDetailByGood(ctx, in)
}

func CreateOrder(ctx context.Context, in *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	conn, err := grpc2.GetGRPCConn(orderconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get order connection: %v", err)
	}
	defer conn.Close()

	cli := orderpb.NewCloudHashingOrderClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.CreateOrder(ctx, in)
}

func CreateGoodPaying(ctx context.Context, in *orderpb.CreateGoodPayingRequest) (*orderpb.CreateGoodPayingResponse, error) {
	conn, err := grpc2.GetGRPCConn(orderconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get order connection: %v", err)
	}
	defer conn.Close()

	cli := orderpb.NewCloudHashingOrderClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.CreateGoodPaying(ctx, in)
}

func CreateGasPaying(ctx context.Context, in *orderpb.CreateGasPayingRequest) (*orderpb.CreateGasPayingResponse, error) {
	conn, err := grpc2.GetGRPCConn(orderconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get order connection: %v", err)
	}
	defer conn.Close()

	cli := orderpb.NewCloudHashingOrderClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.CreateGasPaying(ctx, in)
}

func CreatePayment(ctx context.Context, in *orderpb.CreatePaymentRequest) (*orderpb.CreatePaymentResponse, error) {
	conn, err := grpc2.GetGRPCConn(orderconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get order connection: %v", err)
	}
	defer conn.Close()

	cli := orderpb.NewCloudHashingOrderClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.CreatePayment(ctx, in)
}

//---------------------------------------------------------------------------------------------------------------------------

func GetCouponAllocated(ctx context.Context, in *inspirepb.GetCouponAllocatedDetailRequest) (*inspirepb.GetCouponAllocatedDetailResponse, error) {
	conn, err := grpc2.GetGRPCConn(inspireconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get inspire connection: %v", err)
	}
	defer conn.Close()

	cli := inspirepb.NewCloudHashingInspireClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.GetCouponAllocatedDetail(ctx, in)
}

func GetUserSpecialReduction(ctx context.Context, in *inspirepb.GetUserSpecialReductionRequest) (*inspirepb.GetUserSpecialReductionResponse, error) {
	conn, err := grpc2.GetGRPCConn(inspireconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get inspire connection: %v", err)
	}
	defer conn.Close()

	cli := inspirepb.NewCloudHashingInspireClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.GetUserSpecialReduction(ctx, in)
}

func GetUserInvitationCodeByCode(ctx context.Context, in *inspirepb.GetUserInvitationCodeByCodeRequest) (*inspirepb.GetUserInvitationCodeByCodeResponse, error) {
	conn, err := grpc2.GetGRPCConn(inspireconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get inspire connection: %v", err)
	}
	defer conn.Close()

	cli := inspirepb.NewCloudHashingInspireClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.GetUserInvitationCodeByCode(ctx, in)
}

func GetUserInvitationCodeByAppUser(ctx context.Context, in *inspirepb.GetUserInvitationCodeByAppUserRequest) (*inspirepb.GetUserInvitationCodeByAppUserResponse, error) {
	conn, err := grpc2.GetGRPCConn(inspireconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get inspire connection: %v", err)
	}
	defer conn.Close()

	cli := inspirepb.NewCloudHashingInspireClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.GetUserInvitationCodeByAppUser(ctx, in)
}

func CreateRegistrationInvitation(ctx context.Context, in *inspirepb.CreateRegistrationInvitationRequest) (*inspirepb.CreateRegistrationInvitationResponse, error) {
	conn, err := grpc2.GetGRPCConn(inspireconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get inspire connection: %v", err)
	}
	defer conn.Close()

	cli := inspirepb.NewCloudHashingInspireClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.CreateRegistrationInvitation(ctx, in)
}

func GetRegistrationInvitationsByAppInviter(ctx context.Context, in *inspirepb.GetRegistrationInvitationsByAppInviterRequest) (*inspirepb.GetRegistrationInvitationsByAppInviterResponse, error) {
	conn, err := grpc2.GetGRPCConn(inspireconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get inspire connection: %v", err)
	}
	defer conn.Close()

	cli := inspirepb.NewCloudHashingInspireClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.GetRegistrationInvitationsByAppInviter(ctx, in)
}

//---------------------------------------------------------------------------------------------------------------------------

func CreateBillingAccount(ctx context.Context, in *billingpb.CreateCoinAccountRequest) (*billingpb.CreateCoinAccountResponse, error) {
	conn, err := grpc2.GetGRPCConn(billingconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get billing connection: %v", err)
	}
	defer conn.Close()

	cli := billingpb.NewCloudHashingBillingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.CreateCoinAccount(ctx, in)
}

func GetBillingAccount(ctx context.Context, in *billingpb.GetCoinAccountRequest) (*billingpb.GetCoinAccountResponse, error) {
	conn, err := grpc2.GetGRPCConn(billingconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get billing connection: %v", err)
	}
	defer conn.Close()

	cli := billingpb.NewCloudHashingBillingClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.GetCoinAccount(ctx, in)
}

//---------------------------------------------------------------------------------------------------------------------------

func CreateCoinAddress(ctx context.Context, in *sphinxproxypb.CreateWalletRequest) (*sphinxproxypb.CreateWalletResponse, error) {
	conn, err := grpc2.GetGRPCConn(sphinxproxyconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get sphinxproxy connection: %v", err)
	}
	defer conn.Close()

	cli := sphinxproxypb.NewSphinxProxyClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.CreateWallet(ctx, in)
}

func GetBalance(ctx context.Context, in *sphinxproxypb.GetBalanceRequest) (*sphinxproxypb.GetBalanceResponse, error) {
	conn, err := grpc2.GetGRPCConn(sphinxproxyconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get sphinxproxy connection: %v", err)
	}
	defer conn.Close()

	cli := sphinxproxypb.NewSphinxProxyClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.GetBalance(ctx, in)
}

//---------------------------------------------------------------------------------------------------------------------------

func Signup(ctx context.Context, in *usermgrpb.SignupRequest) (*usermgrpb.SignupResponse, error) {
	conn, err := grpc2.GetGRPCConn(usermgrconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get usermgr connection: %v", err)
	}
	defer conn.Close()

	cli := usermgrpb.NewUserClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.SignUp(ctx, in)
}

func GetUser(ctx context.Context, in *usermgrpb.GetUserRequest) (*usermgrpb.GetUserResponse, error) {
	conn, err := grpc2.GetGRPCConn(usermgrconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get usermgr connection: %v", err)
	}
	defer conn.Close()

	cli := usermgrpb.NewUserClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.GetUser(ctx, in)
}

//---------------------------------------------------------------------------------------------------------------------------

func GetApp(ctx context.Context, in *appmgrpb.GetApplicationRequest) (*appmgrpb.GetApplicationResponse, error) {
	conn, err := grpc2.GetGRPCConn(appmgrconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get appmgr connection: %v", err)
	}
	defer conn.Close()

	cli := appmgrpb.NewApplicationManagementClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.GetApplication(ctx, in)
}

//---------------------------------------------------------------------------------------------------------------------------

func GetReviewsByDomain(ctx context.Context, in *reviewpb.GetReviewsByDomainRequest) (*reviewpb.GetReviewsByDomainResponse, error) {
	conn, err := grpc2.GetGRPCConn(reviewconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get review connection: %v", err)
	}
	defer conn.Close()

	cli := reviewpb.NewReviewServiceClient(conn)

	ctx, cancel := context.WithTimeout(ctx, grpcTimeout)
	defer cancel()

	return cli.GetReviewsByDomain(ctx, in)
}
