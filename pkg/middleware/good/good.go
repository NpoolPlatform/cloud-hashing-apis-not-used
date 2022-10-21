package good

import (
	"context"
	"fmt"

	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	goodsconst "github.com/NpoolPlatform/cloud-hashing-goods/pkg/message/const"

	goodspb "github.com/NpoolPlatform/message/npool/good/mw/v1/good"

	coininfopb "github.com/NpoolPlatform/message/npool/coininfo"
	reviewpb "github.com/NpoolPlatform/message/npool/review-service"

	constant "github.com/NpoolPlatform/cloud-hashing-apis/pkg/const"
	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc"

	goodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/good"
	"github.com/google/uuid"
)

func constructGood(
	info *goodspb.Good,
	coinInfos []*coininfopb.CoinInfo,
	reviews []*reviewpb.Review) (*npool.Good, error) {
	var myCoinInfo *coininfopb.CoinInfo
	var supportedCoinInfos []*coininfopb.CoinInfo

	for _, coinInfo := range coinInfos {
		if coinInfo.ID == info.CoinTypeID {
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
		return nil, fmt.Errorf("not found coin info %v", info.CoinTypeID)
	}

	return &npool.Good{
		Good:         info,
		Main:         myCoinInfo,
		SupportCoins: supportedCoinInfos,
		Reviews:      reviews,
	}, nil
}

func Get(ctx context.Context, in *npool.GetGoodRequest) (*npool.GetGoodResponse, error) {
	good, err := goodmwcli.GetGood(ctx, in.GetID())
	if err != nil {
		return nil, fmt.Errorf("fail get good detail: %v", err)
	}

	coinInfos, err := grpc2.GetCoinInfos(ctx, &coininfopb.GetCoinInfosRequest{
		Offset: 0,
		Limit:  100,
	})
	if err != nil {
		return nil, fmt.Errorf("fail get coin infos: %v", err)
	}

	reviews, err := grpc2.GetReviewsByAppDomainObjectTypeID(ctx, &reviewpb.GetReviewsByAppDomainObjectTypeIDRequest{
		AppID:      uuid.UUID{}.String(),
		Domain:     goodsconst.ServiceName,
		ObjectType: constant.ReviewObjectGood,
		ObjectID:   good.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("fail get reviews by app domain object type id: %v", err)
	}

	detail, err := constructGood(good, coinInfos, reviews)
	if err != nil {
		return nil, fmt.Errorf("fail construct good detail: %v", err)
	}

	return &npool.GetGoodResponse{
		Info: detail,
	}, nil
}
