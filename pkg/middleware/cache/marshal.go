package cache

import (
	"encoding/json"
	"fmt"

	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"
)

type Unmarshal func([]byte) (interface{}, error)

func UnmarshalInvitees(data []byte) ([]*inspirepb.RegistrationInvitation, error) {
	s := []*inspirepb.RegistrationInvitation{}
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("fail unmarshal registration invitation: %v", err)
	}
	return s, nil
}

func UnmarshalInviter(data []byte) (*inspirepb.RegistrationInvitation, error) {
	s := inspirepb.RegistrationInvitation{}
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("fail unmarshal registration invitation: %v", err)
	}
	return &s, nil
}

func UnmarshalAmountSettings(data []byte) ([]*inspirepb.AppPurchaseAmountSetting, error) {
	s := []*inspirepb.AppPurchaseAmountSetting{}
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("fail unmarshal app purchase amount setting: %v", err)
	}
	return s, nil
}
