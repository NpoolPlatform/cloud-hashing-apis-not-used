package currency

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"

	"golang.org/x/xerrors"
)

const (
	coinGeckoAPI = "https://api.coingecko.com/api/v3"
)

func USDPrice(ctx context.Context, coinName string) (float64, error) {
	coin := strings.ToLower(coinName)

	url := fmt.Sprintf("%v%v?ids=%v&vs_currencies=usd", coinGeckoAPI, "/simple/price", coin)
	resp, err := resty.New().R().Get(url)
	if err != nil {
		return 0, xerrors.Errorf("fail get currency: %v", err)
	}
	respMap := map[string]map[string]float64{}
	err = json.Unmarshal(resp.Body(), &respMap)
	if err != nil {
		return 0, xerrors.Errorf("fail parse currency: %v", err)
	}

	priceMap, ok := respMap[coin]
	if !ok {
		return 0, xerrors.Errorf("fail get currency")
	}

	myPrice, ok := priceMap["usd"]
	if !ok {
		return 0, xerrors.Errorf("fail get usd currency")
	}

	return myPrice, nil
}
