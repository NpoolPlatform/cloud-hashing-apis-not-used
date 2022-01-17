package good

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	goodsconst "github.com/NpoolPlatform/cloud-hashing-goods/pkg/message/const"
	goodspb "github.com/NpoolPlatform/message/npool/cloud-hashing-goods"
	coininfopb "github.com/NpoolPlatform/message/npool/coininfo"
	reviewpb "github.com/NpoolPlatform/message/npool/review-service"

	constant "github.com/NpoolPlatform/cloud-hashing-apis/pkg/const"
	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"

	"github.com/google/uuid"

	"golang.org/x/xerrors"
)

func constructGood(
	info *goodspb.GoodDetail,
	coinInfos []*coininfopb.CoinInfo,
	reviews []*reviewpb.Review) (*npool.Good, error) {
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
		Reviews:      reviews,
	}, nil
}

func GetAll(ctx context.Context, in *npool.GetGoodsRequest) (*npool.GetGoodsResponse, error) {
	goodsResp, err := grpc2.GetGoodsDetail(ctx, &goodspb.GetGoodsDetailRequest{
		PageInfo: in.GetPageInfo(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get goods info: %v", err)
	}

	coininfoResp, err := grpc2.GetCoinInfos(ctx, &coininfopb.GetCoinInfosRequest{})
	if err != nil {
		return nil, xerrors.Errorf("fail get coin infos: %v", err)
	}

	details := []*npool.Good{}
	for _, info := range goodsResp.Details {
		review, err := grpc2.GetReviewsByAppDomainObjectTypeID(ctx, &reviewpb.GetReviewsByAppDomainObjectTypeIDRequest{
			AppID:      uuid.UUID{}.String(),
			Domain:     goodsconst.ServiceName,
			ObjectType: constant.ReviewObjectGood,
			ObjectID:   info.ID,
		})
		if err != nil {
			return nil, xerrors.Errorf("fail get reviews by app domain object type id: %v", err)
		}

		detail, err := constructGood(info, coininfoResp.Infos, review.Infos)
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

func Create(ctx context.Context, in *npool.CreateGoodRequest) (*npool.CreateGoodResponse, error) {
	goodResp, err := grpc2.CreateGood(ctx, &goodspb.CreateGoodRequest{
		Info: in.GetInfo(),
	})
	if err != nil {
		return nil, xerrors.Errorf("fail create good: %v", err)
	}

	_, err = grpc2.CreateReview(ctx, &reviewpb.CreateReviewRequest{
		Info: &reviewpb.Review{
			AppID:      uuid.UUID{}.String(),
			Domain:     goodsconst.ServiceName,
			ObjectType: constant.ReviewObjectGood,
			ObjectID:   goodResp.Info.ID,
		},
	})
	if err != nil {
		// TODO: rollback good database
		return nil, xerrors.Errorf("fail create good review: %v", err)
	}

	detail, err := Get(ctx, &npool.GetGoodRequest{
		ID: goodResp.Info.ID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get good detail: %v", err)
	}

	return &npool.CreateGoodResponse{
		Info: detail.Info,
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

	review, err := grpc2.GetReviewsByAppDomainObjectTypeID(ctx, &reviewpb.GetReviewsByAppDomainObjectTypeIDRequest{
		AppID:      uuid.UUID{}.String(),
		Domain:     goodsconst.ServiceName,
		ObjectType: constant.ReviewObjectGood,
		ObjectID:   goodResp.Detail.ID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get reviews by app domain object type id: %v", err)
	}

	detail, err := constructGood(goodResp.Detail, coininfoResp.Infos, review.Infos)
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
		review, err := grpc2.GetReviewsByAppDomainObjectTypeID(ctx, &reviewpb.GetReviewsByAppDomainObjectTypeIDRequest{
			AppID:      uuid.UUID{}.String(),
			Domain:     goodsconst.ServiceName,
			ObjectType: constant.ReviewObjectGood,
			ObjectID:   info.Good.ID,
		})
		if err != nil {
			return nil, xerrors.Errorf("fail get reviews by app domain object type id: %v", err)
		}

		detail, err := constructGood(info.Good, coininfoResp.Infos, review.Infos)
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
