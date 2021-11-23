package gooddetail

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	"github.com/NpoolPlatform/cloud-hashing-apis/message/npool"

	goodspb "github.com/NpoolPlatform/cloud-hashing-goods/message/npool"
	coininfopb "github.com/NpoolPlatform/message/npool/coininfo"

	grpc2 "github.com/NpoolPlatform/cloud-hashing-apis/pkg/grpc" //nolint

	"golang.org/x/xerrors"
	"google.golang.org/protobuf/types/known/emptypb"
)

func constructGoodDetail(info *goodspb.GoodDetail, coinInfos []*coininfopb.CoinInfo) (*npool.GoodDetail, error) {
	var myCoinInfo *npool.CoinInfo
	var supportedCoinInfos []*npool.CoinInfo

	for _, coinInfo := range coinInfos {
		if coinInfo.ID == info.CoinInfoID {
			myCoinInfo = &npool.CoinInfo{
				ID:      coinInfo.ID,
				PreSale: coinInfo.PreSale,
				Name:    coinInfo.Name,
				Unit:    coinInfo.Unit,
				Logo:    "",
			}
			break
		}
	}

	for _, coinInfo := range coinInfos {
		for _, coinInfoID := range info.SupportCoinTypeIDs {
			if coinInfoID == coinInfo.ID {
				supportedCoinInfos = append(supportedCoinInfos, &npool.CoinInfo{
					ID:      coinInfo.ID,
					PreSale: coinInfo.PreSale,
					Name:    coinInfo.Name,
					Unit:    coinInfo.Unit,
					Logo:    "",
				})
			}
		}
	}

	if myCoinInfo == nil {
		return nil, xerrors.Errorf("not found coin info %v", info.CoinInfoID)
	}

	var inheritFrom *npool.GoodInfo
	if info.InheritFromGood != nil {
		inheritFrom = &npool.GoodInfo{
			ID:                 info.InheritFromGood.ID,
			Title:              info.InheritFromGood.Title,
			DeviceInfoID:       info.InheritFromGood.DeviceInfoID,
			SeparateFee:        info.InheritFromGood.SeparateFee,
			UnitPower:          info.InheritFromGood.UnitPower,
			DurationDays:       info.InheritFromGood.DurationDays,
			CoinInfoID:         info.InheritFromGood.CoinInfoID,
			Actuals:            info.InheritFromGood.Actuals,
			DeliveryAt:         info.InheritFromGood.DeliveryAt,
			InheritFromGoodID:  info.InheritFromGood.InheritFromGoodID,
			VendorLocationID:   info.InheritFromGood.VendorLocationID,
			Price:              info.InheritFromGood.Price,
			PriceCurrency:      info.InheritFromGood.PriceCurrency,
			BenefitType:        info.InheritFromGood.BenefitType,
			Classic:            info.InheritFromGood.Classic,
			SupportCoinTypeIDs: info.InheritFromGood.SupportCoinTypeIDs,
			Total:              info.InheritFromGood.Total,
			Unit:               info.InheritFromGood.Unit,
			Start:              info.InheritFromGood.Start,
		}
	}

	var fees []*npool.Fee //nolint
	for _, fee := range info.Fees {
		fees = append(fees, &npool.Fee{
			Fee: &npool.GoodFee{
				ID:             fee.Fee.ID,
				AppID:          fee.Fee.ID,
				FeeType:        fee.Fee.FeeType,
				FeeDescription: fee.Fee.FeeDescription,
				PayType:        fee.Fee.PayType,
			},
			Value: fee.Value,
		})
	}

	return &npool.GoodDetail{
		ID: info.ID,
		DeviceInfo: &npool.DeviceInfo{
			ID:              info.DeviceInfo.ID,
			Type:            info.DeviceInfo.Type,
			Manufacturer:    info.DeviceInfo.Manufacturer,
			PowerComsuption: info.DeviceInfo.PowerComsuption,
			ShipmentAt:      info.DeviceInfo.ShipmentAt,
		},
		SeparateFee:     info.SeparateFee,
		UnitPower:       info.UnitPower,
		DurationDays:    info.DurationDays,
		CoinInfo:        myCoinInfo,
		Actuals:         info.Actuals,
		DeliveryAt:      info.DeliveryAt,
		InheritFromGood: inheritFrom,
		VendorLocation: &npool.VendorLocationInfo{
			ID:       info.VendorLocation.ID,
			Country:  info.VendorLocation.Country,
			Province: info.VendorLocation.Province,
			City:     info.VendorLocation.City,
			Address:  info.VendorLocation.Address,
		},
		Price: info.Price,
		PriceCurrency: &npool.PriceCurrency{
			ID:     info.PriceCurrency.ID,
			Name:   info.PriceCurrency.Name,
			Unit:   info.PriceCurrency.Unit,
			Symbol: info.PriceCurrency.Symbol,
		},
		BenefitType:  info.BenefitType,
		Classic:      info.Classic,
		SupportCoins: supportedCoinInfos,
		Total:        info.Total,
		Extra: &npool.GoodExtraInfo{
			ID:        info.Extra.ID,
			GoodID:    info.Extra.GoodID,
			Posters:   info.Extra.Posters,
			Labels:    info.Extra.Labels,
			OutSale:   info.Extra.OutSale,
			PreSale:   info.Extra.PreSale,
			VoteCount: info.Extra.VoteCount,
			Rating:    info.Extra.Rating,
		},
		Start: info.Start,
		Unit:  info.Unit,
		Title: info.Title,
		Fees:  fees,
	}, nil
}

func GetAll(ctx context.Context, in *npool.GetGoodsDetailRequest) (*npool.GetGoodsDetailResponse, error) {
	goodsResp, err := grpc2.GetGoodsDetail(ctx, &goodspb.GetGoodsDetailRequest{
		PageInfo: &goodspb.PageInfo{
			PageIndex: in.GetPageInfo().GetPageIndex(),
			PageSize:  in.GetPageInfo().GetPageSize(),
		},
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get goods info: %v", err)
	}

	coininfoResp, err := grpc2.GetCoinInfos(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, xerrors.Errorf("fail get coin infos: %v", err)
	}

	details := []*npool.GoodDetail{}
	for _, info := range goodsResp.Details {
		detail, err := constructGoodDetail(info, coininfoResp.Infos)
		if err != nil {
			logger.Sugar().Errorf("fail to get coin info %v: %v", info.CoinInfoID, err)
			continue
		}

		details = append(details, detail)
	}

	return &npool.GetGoodsDetailResponse{
		Details: details,
	}, nil
}

func Get(ctx context.Context, goodID string) (*npool.GoodDetail, error) {
	goodResp, err := grpc2.GetGoodDetail(ctx, &goodspb.GetGoodDetailRequest{
		ID: goodID,
	})
	if err != nil {
		return nil, xerrors.Errorf("fail get good detail: %v", err)
	}

	coininfoResp, err := grpc2.GetCoinInfos(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, xerrors.Errorf("fail get coin infos: %v", err)
	}

	return constructGoodDetail(goodResp.Detail, coininfoResp.Infos)
}
