package gooddetail

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	commonpb "github.com/NpoolPlatform/message/npool"
	goodspb "github.com/NpoolPlatform/message/npool/cloud-hashing-goods"
	coininfopb "github.com/NpoolPlatform/message/npool/coininfo"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc" //nolint

	"golang.org/x/xerrors"
)

// TODO: add review result
func constructGood(info *goodspb.GoodDetail, coinInfos []*coininfopb.CoinInfo) (*npool.Good, error) {
	var myCoinInfo *coininfopb.CoinInfo
	var supportedCoinInfos []*coininfopb.CoinInfo

	for _, coinInfo := range coinInfos {
		if coinInfo.ID == info.CoinInfoID {
			myCoinInfo = coinInfo
			break
		}
	}

	for _, coinInfo := range coinInfos {
		for _, coinInfoID := range info.SupportCoinTypeIDs {
			if coinInfoID == coinInfo.ID {
				supportedCoinInfos = append(supportedCoinInfos, coinInfo)
			}
		}
	}

	if myCoinInfo == nil {
		return nil, xerrors.Errorf("not found coin info %v", info.CoinInfoID)
	}

	return &npool.Good{
		Good:         info,
		Main:         myCoinInfo,
		SupportCoins: supportedCoinInfos,
	}, nil
}

func GetAll(ctx context.Context, in *npool.GetGoodsRequest) (*npool.GetGoodsResponse, error) {
	logger.Sugar().Infof("get all %v", in)
	goodsResp, err := grpc2.GetGoodsDetail(ctx, &goodspb.GetGoodsDetailRequest{
		PageInfo: &commonpb.PageInfo{
			PageIndex: in.GetPageInfo().GetPageIndex(),
			PageSize:  in.GetPageInfo().GetPageSize(),
		},
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get goods info: %v", err)
	}

	logger.Sugar().Infof("get all %v: %v", in, len(goodsResp.Details))

	coininfoResp, err := grpc2.GetCoinInfos(ctx, &coininfopb.GetCoinInfosRequest{})
	if err != nil {
		return nil, xerrors.Errorf("fail get coin infos: %v", err)
	}

	details := []*npool.Good{}
	for _, info := range goodsResp.Details {
		detail, err := constructGood(info, coininfoResp.Infos)
		if err != nil {
			logger.Sugar().Errorf("fail to get coin info %v: %v", info.CoinInfoID, err)
			continue
		}

		details = append(details, detail)
	}

	return &npool.GetGoodsResponse{
		Infos: details,
	}, nil
}

func Get(ctx context.Context, in *npool.GetGoodRequest) (*npool.GetGoodResponse, error) {
	goodResp, err := grpc2.GetGoodDetail(ctx, &goodspb.GetGoodDetailRequest{
		ID: in.GetID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get good detail: %v", err)
	}

	coininfoResp, err := grpc2.GetCoinInfos(ctx, &coininfopb.GetCoinInfosRequest{})
	if err != nil {
		return nil, xerrors.Errorf("fail get coin infos: %v", err)
	}

	detail, err := constructGood(goodResp.Detail, coininfoResp.Infos)
	if err != nil {
		return nil, xerrors.Errorf("fail construct good detail: %v", err)
	}

	return &npool.GetGoodResponse{
		Info: detail,
	}, nil
}

func GetRecommendsByApp(ctx context.Context, in *npool.GetRecommendGoodsByAppRequest) (*npool.GetRecommendGoodsByAppResponse, error) {
	goodsResp, err := grpc2.GetRecommendGoodsByApp(ctx, &goodspb.GetRecommendGoodsByAppRequest{
		AppID: in.GetAppID(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get goods info: %v", err)
	}

	coininfoResp, err := grpc2.GetCoinInfos(ctx, &coininfopb.GetCoinInfosRequest{})
	if err != nil {
		return nil, xerrors.Errorf("fail get coin infos: %v", err)
	}

	details := []*npool.RecommendGood{}

	for _, info := range goodsResp.Infos {
		detail, err := constructGood(info.Good, coininfoResp.Infos)
		if err != nil {
			logger.Sugar().Errorf("fail to get coin info %v: %v", info.Good.CoinInfoID, err)
			continue
		}

		details = append(details, &npool.RecommendGood{
			Recommend: info.Recommend,
			Good:      detail,
		})
	}

	return &npool.GetRecommendGoodsByAppResponse{
		Infos: details,
	}, nil
}
