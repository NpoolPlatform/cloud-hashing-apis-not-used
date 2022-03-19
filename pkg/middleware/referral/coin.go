package referral

import (
	"context"

	cache "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/cache"
	orderconst "github.com/NpoolPlatform/cloud-hashing-order/pkg/const"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"

	"golang.org/x/xerrors"
)

const (
	cacheCoinSummaries = "referral:coin:summaries"
)

func getCoinSummaries(ctx context.Context, appID, userID string) ([]*npool.CoinSummary, error) {
	mySummaries := cache.GetEntry(CacheKey(appID, userID, cacheCoinSummaries))
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
				Unit:       order.Good.Good.Good.Unit,
			}
			summaries = append(summaries, summary)
		}

		summary.Units += order.Order.Order.Units
		summary.Amount += order.Order.Payment.Amount * order.Order.Payment.CoinUSDCurrency
	}

	cache.AddEntry(CacheKey(appID, userID, cacheCoinSummaries), summaries)

	return summaries, nil
}
