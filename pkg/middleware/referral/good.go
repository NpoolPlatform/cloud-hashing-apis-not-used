package referral

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	cache "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/cache"
	orderconst "github.com/NpoolPlatform/cloud-hashing-order/pkg/const"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	"golang.org/x/xerrors"
)

const (
	cacheGoodSummaries = "referral:good:summaries"
)

func getGoodSummaries(ctx context.Context, appID, userID string) ([]*npool.GoodSummary, error) {
	mySummaries := cache.GetEntry(CacheKey(appID, userID, cacheGoodSummaries))
	if mySummaries != nil {
		return mySummaries.([]*npool.GoodSummary), nil
	}

	orders, err := GetOrders(ctx, appID, userID)
	if err != nil {
		return nil, xerrors.Errorf("fail get orders: %v", err)
	}

	summaries := []*npool.GoodSummary{}
	for _, order := range orders {
		if order.Order.Payment == nil || order.Order.Payment.State != orderconst.PaymentStateDone {
			continue
		}

		var summary *npool.GoodSummary
		for _, sum := range summaries {
			if sum.GoodID == order.Good.Good.Good.ID {
				summary = sum
				break
			}
		}

		if summary == nil {
			summary = &npool.GoodSummary{
				GoodID:     order.Good.Good.Good.ID,
				CoinTypeID: order.Good.Main.ID,
				CoinName:   order.Good.Main.Unit,
				Unit:       order.Good.Good.Good.Unit,
			}
			summaries = append(summaries, summary)
		}

		amount := order.Order.Payment.Amount * order.Order.Payment.CoinUSDCurrency
		logger.Sugar().Infof("order %v good %v coin %v | %v units %v amount %v user %v",
			order.Order.Order.ID, order.Good.Good.Good.ID,
			order.Good.Main.ID, order.Good.Main.Unit,
			order.Order.Order.Units, amount, userID)

		summary.Units += order.Order.Order.Units
		summary.Amount += amount
	}

	if len(summaries) > 0 {
		cache.AddEntry(CacheKey(appID, userID, cacheGoodSummaries), summaries)
	}

	return summaries, nil
}
