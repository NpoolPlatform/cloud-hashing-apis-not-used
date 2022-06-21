package cache

import (
	"encoding/json"
	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"
)

type AppPurchaseAmountSettings []*inspirepb.AppPurchaseAmountSetting

func (s AppPurchaseAmountSettings) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s AppPurchaseAmountSettings) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}
