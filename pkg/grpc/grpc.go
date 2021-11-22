package grpc

import (
	"context"

	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"

	goodspb "github.com/NpoolPlatform/cloud-hashing-goods/message/npool"
	goodsconst "github.com/NpoolPlatform/cloud-hashing-goods/pkg/message/const" //nolint

	coininfopb "github.com/NpoolPlatform/message/npool/coininfo"
	coininfoconst "github.com/NpoolPlatform/sphinx-coininfo/pkg/message/const" //nolint

	"golang.org/x/xerrors"
	"google.golang.org/protobuf/types/known/emptypb"
)

func GetGoodsDetail(ctx context.Context, in *goodspb.GetGoodsDetailRequest) (*goodspb.GetGoodsDetailResponse, error) {
	conn, err := grpc2.GetGRPCConn(goodsconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get goods connection: %v", err)
	}

	cli := goodspb.NewCloudHashingGoodsClient(conn)
	return cli.GetGoodsDetail(ctx, in)
}

func GetCoinInfos(ctx context.Context, in *emptypb.Empty) (*coininfopb.GetCoinInfosResponse, error) {
	conn, err := grpc2.GetGRPCConn(coininfoconst.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, xerrors.Errorf("fail get coininfo connection: %v", err)
	}

	cli := coininfopb.NewSphinxCoinInfoClient(conn)
	return cli.GetCoinInfos(ctx, in)
}
