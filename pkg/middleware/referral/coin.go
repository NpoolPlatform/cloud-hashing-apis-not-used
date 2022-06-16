package referral

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	cache "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/cache"
	cachekey "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/referral/cachekey"
	orderconst "github.com/NpoolPlatform/cloud-hashing-order/pkg/const"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	"golang.org/x/xerrors"
)

const (
	cacheCoinSummaries = "referral:coin:summaries"
)

func getCoinSummaries(ctx context.Context, appID, userID string) ([]*npool.CoinSummary, error) {
	mySummaries := cache.GetEntry(cachekey.CacheKey(appID, userID, cacheCoinSummaries))
	if mySummaries != nil {
		return mySummaries.([]*npool.CoinSummary), nil
	}

	orders, err := GetOrders(ctx, appID, userID)
	if err != nil {
		return nil, xerrors.Errorf("fail get orders: %v", err)
	}

	summaries := []*npool.CoinSummary{}
	for _, order := range orders {
		if order.Order.Payment == nil || order.Order.Payment.State != orderconst.PaymentStateDone {
			continue
		}

		var summary *npool.CoinSummary
		for _, sum := range summaries {
			if sum.CoinTypeID == order.Good.Main.ID {
				summary = sum
				break
			}
		}

		if summary == nil {
			summary = &npool.CoinSummary{
				CoinTypeID: order.Good.Main.ID,
				CoinName:   order.Good.Main.Unit,
				Unit:       order.Good.Good.Good.Unit,
			}
			summaries = append(summaries, summary)
		}

		amount := order.Order.Payment.Amount * order.Order.Payment.CoinUSDCurrency
		logger.Sugar().Infof("order %v coin %v | %v units %v amount %v user %v",
			order.Order.Order.ID, order.Good.Main.ID, order.Good.Main.Unit,
			order.Order.Order.Units, amount, userID)

		summary.Units += order.Order.Order.Units
		summary.Amount += amount
	}

	if len(summaries) > 0 {
		cache.AddEntry(cachekey.CacheKey(appID, userID, cacheCoinSummaries), summaries)
	}

	return summaries, nil
}
