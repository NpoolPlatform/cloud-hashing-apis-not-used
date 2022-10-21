package referral

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	cache "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/cache"
	cachekey "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/referral/cachekey"
	setting "github.com/NpoolPlatform/cloud-hashing-apis/pkg/middleware/referral/setting"
	orderconst "github.com/NpoolPlatform/cloud-hashing-order/pkg/const"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"
)

const (
	cacheGoodSummaries = "referral:good:summaries"
)

func getGoodSummaries(ctx context.Context, appID, userID string) ([]*npool.GoodSummary, error) { //nolint
	mySummaries := cache.GetEntry(cachekey.CacheKey(appID, userID, cacheGoodSummaries), func(data []byte) (interface{}, error) {
		return cache.UnmarshalGoodSummaries(data)
	})
	if mySummaries != nil {
		return mySummaries.([]*npool.GoodSummary), nil
	}

	orders, err := GetOrders(ctx, appID, userID)
	if err != nil {
		return nil, fmt.Errorf("fail get orders: %v", err)
	}

	settings, err := setting.GetAmountSettingsByAppUser(ctx, appID, userID)
	if err != nil {
		return nil, fmt.Errorf("fail get amount settings: %v", err)
	}

	summaries := []*npool.GoodSummary{}
	for _, order := range orders {
		if order.Order.Payment == nil || order.Order.Payment.State != orderconst.PaymentStateDone {
			continue
		}

		var summary *npool.GoodSummary
		for _, sum := range summaries {
			if sum.GoodID == order.Good.Good.ID {
				summary = sum
				break
			}
		}

		percent := uint32(0)
		s := setting.GetGoodAmountSetting(settings, order.Good.Good.ID)
		if s != nil {
			percent = s.Percent
		}

		if summary == nil {
			summary = &npool.GoodSummary{
				GoodID:     order.Good.Good.ID,
				CoinTypeID: order.Good.Main.ID,
				CoinName:   order.Good.Main.Unit,
				Unit:       order.Good.Good.Unit,
			}
			summaries = append(summaries, summary)
		}

		amount := order.Order.Payment.Amount * order.Order.Payment.CoinUSDCurrency
		logger.Sugar().Infof("order %v good %v coin %v | %v units %v amount %v user %v",
			order.Order.Order.ID, order.Good.Good.ID,
			order.Good.Main.ID, order.Good.Main.Unit,
			order.Order.Order.Units, amount, userID)

		summary.Units += order.Order.Order.Units
		summary.Amount += amount
		summary.Percent = percent
	}

	if len(summaries) > 0 {
		cache.AddEntry(cachekey.CacheKey(appID, userID, cacheGoodSummaries), summaries)
	}

	return summaries, nil
}

func GetGoodSummaries(ctx context.Context, appID, userID string) ([]*npool.GoodSummary, error) {
	comms, err := getGoodSummaries(ctx, appID, userID)
	if err != nil {
		return nil, fmt.Errorf("fail get user good summaries: %v", err)
	}

	invitees, err := GetLayeredInvitees(ctx, appID, userID)
	if err != nil {
		return nil, fmt.Errorf("fail get user invitees: %v", err)
	}

	for _, iv := range invitees {
		commissions, err := getGoodSummaries(ctx, appID, iv.InviteeID)
		if err != nil {
			return nil, fmt.Errorf("fail get user good summaries: %v", err)
		}
		for _, commission := range commissions {
			found := false
			for _, comm := range comms {
				if comm.GoodID == commission.GoodID {
					comm.Units += commission.Units
					comm.Amount += commission.Amount
					found = true
					break
				}
			}
			if !found {
				comms = append(comms, commission)
			}
		}
	}

	return comms, err
}

func UpdateGoodSummariesCache(ctx context.Context, appID, userID string, summaries []*npool.GoodSummary) {
	if len(summaries) > 0 {
		cache.AddEntry(cachekey.CacheKey(appID, userID, cacheGoodSummaries), summaries)
	}
}
