package cache

import (
	"encoding/json"
	"fmt"

	appusermgrpb "github.com/NpoolPlatform/message/npool/appuser/mgr/v1"
	npool "github.com/NpoolPlatform/message/npool/cloud-hashing-apis"
	goodspb "github.com/NpoolPlatform/message/npool/cloud-hashing-goods"
	inspirepb "github.com/NpoolPlatform/message/npool/cloud-hashing-inspire"
	coininfopb "github.com/NpoolPlatform/message/npool/coininfo"
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

func UnmarshalCouponAllocated(data []byte) (*inspirepb.CouponAllocatedDetail, error) {
	s := inspirepb.CouponAllocatedDetail{}
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("fail unmarshal coupon allocated: %v", err)
	}
	return &s, nil
}

func UnmarshalOrders(data []byte) ([]*npool.Order, error) {
	s := []*npool.Order{}
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("fail unmarshal orders: %v", err)
	}
	return s, nil
}

func UnmarshalAppUsers(data []byte) ([]*appusermgrpb.AppUser, error) {
	s := []*appusermgrpb.AppUser{}
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("fail unmarshal app users: %v", err)
	}
	return s, nil
}

func UnmarshalAppUser(data []byte) (*appusermgrpb.AppUser, error) {
	s := appusermgrpb.AppUser{}
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("fail unmarshal app users: %v", err)
	}
	return &s, nil
}

func UnmarshalAppUserExtra(data []byte) (*appusermgrpb.AppUserExtra, error) {
	s := appusermgrpb.AppUserExtra{}
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("fail unmarshal app user extra: %v", err)
	}
	return &s, nil
}

func UnmarshalAppGoodInfo(data []byte) (*goodspb.AppGoodInfo, error) {
	s := goodspb.AppGoodInfo{}
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("fail unmarshal app good info: %v", err)
	}
	return &s, nil
}

func UnmarshalAppGoodPromotion(data []byte) (*goodspb.AppGoodPromotion, error) {
	s := goodspb.AppGoodPromotion{}
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("fail unmarshal app good promotion: %v", err)
	}
	return &s, nil
}

func UnmarshalCoinInfo(data []byte) (*coininfopb.CoinInfo, error) {
	s := coininfopb.CoinInfo{}
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("fail unmarshal coin info: %v", err)
	}
	return &s, nil
}

func UnmarshalGood(data []byte) (*npool.Good, error) {
	s := npool.Good{}
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("fail unmarshal good: %v", err)
	}
	return &s, nil
}

func UnmarshalCoinSummaries(data []byte) ([]*npool.CoinSummary, error) {
	s := []*npool.CoinSummary{}
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("fail unmarshal coin summary: %v", err)
	}
	return s, nil
}

func UnmarshalGoodSummaries(data []byte) ([]*npool.GoodSummary, error) {
	s := []*npool.GoodSummary{}
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("fail unmarshal good summary: %v", err)
	}
	return s, nil
}

func UnmarshalSpecialOffer(data []byte) (*inspirepb.UserSpecialReduction, error) {
	s := inspirepb.UserSpecialReduction{}
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("fail unmarshal user special reduction: %v", err)
	}
	return &s, nil
}
