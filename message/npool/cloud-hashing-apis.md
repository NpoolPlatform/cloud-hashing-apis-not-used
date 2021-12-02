# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [npool/cloud-hashing-apis.proto](#npool/cloud-hashing-apis.proto)
    - [Account](#cloud.hashing.apis.v1.Account)
    - [CoinInfo](#cloud.hashing.apis.v1.CoinInfo)
    - [Compensate](#cloud.hashing.apis.v1.Compensate)
    - [Coupon](#cloud.hashing.apis.v1.Coupon)
    - [CouponPool](#cloud.hashing.apis.v1.CouponPool)
    - [CreateOrderPaymentRequest](#cloud.hashing.apis.v1.CreateOrderPaymentRequest)
    - [CreateOrderPaymentResponse](#cloud.hashing.apis.v1.CreateOrderPaymentResponse)
    - [DeviceInfo](#cloud.hashing.apis.v1.DeviceInfo)
    - [Discount](#cloud.hashing.apis.v1.Discount)
    - [DiscountPool](#cloud.hashing.apis.v1.DiscountPool)
    - [Fee](#cloud.hashing.apis.v1.Fee)
    - [GasPaying](#cloud.hashing.apis.v1.GasPaying)
    - [GetGoodDetailRequest](#cloud.hashing.apis.v1.GetGoodDetailRequest)
    - [GetGoodDetailResponse](#cloud.hashing.apis.v1.GetGoodDetailResponse)
    - [GetGoodsDetailRequest](#cloud.hashing.apis.v1.GetGoodsDetailRequest)
    - [GetGoodsDetailResponse](#cloud.hashing.apis.v1.GetGoodsDetailResponse)
    - [GetOrderDetailRequest](#cloud.hashing.apis.v1.GetOrderDetailRequest)
    - [GetOrderDetailResponse](#cloud.hashing.apis.v1.GetOrderDetailResponse)
    - [GetOrdersDetailByAppRequest](#cloud.hashing.apis.v1.GetOrdersDetailByAppRequest)
    - [GetOrdersDetailByAppResponse](#cloud.hashing.apis.v1.GetOrdersDetailByAppResponse)
    - [GetOrdersDetailByAppUserRequest](#cloud.hashing.apis.v1.GetOrdersDetailByAppUserRequest)
    - [GetOrdersDetailByAppUserResponse](#cloud.hashing.apis.v1.GetOrdersDetailByAppUserResponse)
    - [GetOrdersDetailByGoodRequest](#cloud.hashing.apis.v1.GetOrdersDetailByGoodRequest)
    - [GetOrdersDetailByGoodResponse](#cloud.hashing.apis.v1.GetOrdersDetailByGoodResponse)
    - [GoodDetail](#cloud.hashing.apis.v1.GoodDetail)
    - [GoodExtraInfo](#cloud.hashing.apis.v1.GoodExtraInfo)
    - [GoodFee](#cloud.hashing.apis.v1.GoodFee)
    - [GoodInfo](#cloud.hashing.apis.v1.GoodInfo)
    - [GoodPaying](#cloud.hashing.apis.v1.GoodPaying)
    - [Order](#cloud.hashing.apis.v1.Order)
    - [OrderDetail](#cloud.hashing.apis.v1.OrderDetail)
    - [OrderFee](#cloud.hashing.apis.v1.OrderFee)
    - [OutOfGas](#cloud.hashing.apis.v1.OutOfGas)
    - [PageInfo](#cloud.hashing.apis.v1.PageInfo)
    - [Payment](#cloud.hashing.apis.v1.Payment)
    - [PriceCurrency](#cloud.hashing.apis.v1.PriceCurrency)
    - [SubmitOrderRequest](#cloud.hashing.apis.v1.SubmitOrderRequest)
    - [SubmitOrderResponse](#cloud.hashing.apis.v1.SubmitOrderResponse)
    - [TargetAreaInfo](#cloud.hashing.apis.v1.TargetAreaInfo)
    - [UserSpecialReduction](#cloud.hashing.apis.v1.UserSpecialReduction)
    - [VendorLocationInfo](#cloud.hashing.apis.v1.VendorLocationInfo)
    - [VersionResponse](#cloud.hashing.apis.v1.VersionResponse)
  
    - [CloudHashingApis](#cloud.hashing.apis.v1.CloudHashingApis)
  
- [Scalar Value Types](#scalar-value-types)



<a name="npool/cloud-hashing-apis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## npool/cloud-hashing-apis.proto



<a name="cloud.hashing.apis.v1.Account"></a>

### Account



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ID | [string](#string) |  |  |
| CoinTypeID | [string](#string) |  |  |
| Address | [string](#string) |  |  |
| AppID | [string](#string) |  |  |
| UserID | [string](#string) |  |  |






<a name="cloud.hashing.apis.v1.CoinInfo"></a>

### CoinInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ID | [string](#string) |  |  |
| PreSale | [bool](#bool) |  |  |
| Name | [string](#string) |  |  |
| Unit | [string](#string) |  |  |
| Logo | [string](#string) |  |  |






<a name="cloud.hashing.apis.v1.Compensate"></a>

### Compensate



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ID | [string](#string) |  |  |
| OrderID | [string](#string) |  |  |
| Start | [uint32](#uint32) |  |  |
| End | [uint32](#uint32) |  |  |
| Message | [string](#string) |  |  |






<a name="cloud.hashing.apis.v1.Coupon"></a>

### Coupon



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ID | [string](#string) |  |  |
| UserID | [string](#string) |  |  |
| AppID | [string](#string) |  |  |
| Pool | [CouponPool](#cloud.hashing.apis.v1.CouponPool) |  |  |






<a name="cloud.hashing.apis.v1.CouponPool"></a>

### CouponPool



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ID | [string](#string) |  |  |
| AppID | [string](#string) |  |  |
| Denomination | [double](#double) |  |  |
| Start | [uint32](#uint32) |  |  |
| DurationDays | [int32](#int32) |  |  |
| Message | [string](#string) |  |  |
| Name | [string](#string) |  |  |






<a name="cloud.hashing.apis.v1.CreateOrderPaymentRequest"></a>

### CreateOrderPaymentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| OrderID | [string](#string) |  |  |
| PaymentCoinTypeID | [string](#string) |  |  |
| Fees | [OrderFee](#cloud.hashing.apis.v1.OrderFee) | repeated |  |






<a name="cloud.hashing.apis.v1.CreateOrderPaymentResponse"></a>

### CreateOrderPaymentResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Info | [OrderDetail](#cloud.hashing.apis.v1.OrderDetail) |  |  |






<a name="cloud.hashing.apis.v1.DeviceInfo"></a>

### DeviceInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ID | [string](#string) |  |  |
| Type | [string](#string) |  |  |
| Manufacturer | [string](#string) |  |  |
| PowerComsuption | [int32](#int32) |  |  |
| ShipmentAt | [int32](#int32) |  |  |






<a name="cloud.hashing.apis.v1.Discount"></a>

### Discount



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ID | [string](#string) |  |  |
| UserID | [string](#string) |  |  |
| AppID | [string](#string) |  |  |
| Pool | [DiscountPool](#cloud.hashing.apis.v1.DiscountPool) |  |  |






<a name="cloud.hashing.apis.v1.DiscountPool"></a>

### DiscountPool



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ID | [string](#string) |  |  |
| AppID | [string](#string) |  |  |
| ReleaseByUserID | [string](#string) |  |  |
| Discount | [uint32](#uint32) |  |  |
| Start | [uint32](#uint32) |  |  |
| DurationDays | [int32](#int32) |  |  |
| Message | [string](#string) |  |  |
| Name | [string](#string) |  |  |






<a name="cloud.hashing.apis.v1.Fee"></a>

### Fee



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Fee | [GoodFee](#cloud.hashing.apis.v1.GoodFee) |  |  |
| Value | [double](#double) |  |  |






<a name="cloud.hashing.apis.v1.GasPaying"></a>

### GasPaying



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ID | [string](#string) |  |  |
| OrderID | [string](#string) |  |  |
| PaymentID | [string](#string) |  |  |
| DurationMinutes | [uint32](#uint32) |  |  |






<a name="cloud.hashing.apis.v1.GetGoodDetailRequest"></a>

### GetGoodDetailRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ID | [string](#string) |  |  |






<a name="cloud.hashing.apis.v1.GetGoodDetailResponse"></a>

### GetGoodDetailResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Detail | [GoodDetail](#cloud.hashing.apis.v1.GoodDetail) |  |  |






<a name="cloud.hashing.apis.v1.GetGoodsDetailRequest"></a>

### GetGoodsDetailRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| AppID | [string](#string) |  |  |
| PageInfo | [PageInfo](#cloud.hashing.apis.v1.PageInfo) |  |  |






<a name="cloud.hashing.apis.v1.GetGoodsDetailResponse"></a>

### GetGoodsDetailResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Details | [GoodDetail](#cloud.hashing.apis.v1.GoodDetail) | repeated |  |
| Total | [int32](#int32) |  |  |






<a name="cloud.hashing.apis.v1.GetOrderDetailRequest"></a>

### GetOrderDetailRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ID | [string](#string) |  |  |






<a name="cloud.hashing.apis.v1.GetOrderDetailResponse"></a>

### GetOrderDetailResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Detail | [OrderDetail](#cloud.hashing.apis.v1.OrderDetail) |  |  |






<a name="cloud.hashing.apis.v1.GetOrdersDetailByAppRequest"></a>

### GetOrdersDetailByAppRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| AppID | [string](#string) |  |  |






<a name="cloud.hashing.apis.v1.GetOrdersDetailByAppResponse"></a>

### GetOrdersDetailByAppResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Details | [OrderDetail](#cloud.hashing.apis.v1.OrderDetail) | repeated |  |






<a name="cloud.hashing.apis.v1.GetOrdersDetailByAppUserRequest"></a>

### GetOrdersDetailByAppUserRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| AppID | [string](#string) |  |  |
| UserID | [string](#string) |  |  |






<a name="cloud.hashing.apis.v1.GetOrdersDetailByAppUserResponse"></a>

### GetOrdersDetailByAppUserResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Details | [OrderDetail](#cloud.hashing.apis.v1.OrderDetail) | repeated |  |






<a name="cloud.hashing.apis.v1.GetOrdersDetailByGoodRequest"></a>

### GetOrdersDetailByGoodRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| GoodID | [string](#string) |  |  |






<a name="cloud.hashing.apis.v1.GetOrdersDetailByGoodResponse"></a>

### GetOrdersDetailByGoodResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Details | [OrderDetail](#cloud.hashing.apis.v1.OrderDetail) | repeated |  |






<a name="cloud.hashing.apis.v1.GoodDetail"></a>

### GoodDetail



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ID | [string](#string) |  |  |
| DeviceInfo | [DeviceInfo](#cloud.hashing.apis.v1.DeviceInfo) |  |  |
| SeparateFee | [bool](#bool) |  |  |
| UnitPower | [int32](#int32) |  |  |
| DurationDays | [int32](#int32) |  |  |
| CoinInfo | [CoinInfo](#cloud.hashing.apis.v1.CoinInfo) |  |  |
| Actuals | [bool](#bool) |  |  |
| DeliveryAt | [int32](#int32) |  |  |
| InheritFromGood | [GoodInfo](#cloud.hashing.apis.v1.GoodInfo) |  |  |
| VendorLocation | [VendorLocationInfo](#cloud.hashing.apis.v1.VendorLocationInfo) |  |  |
| Price | [double](#double) |  |  |
| BenefitType | [string](#string) |  |  |
| Classic | [bool](#bool) |  |  |
| SupportCoins | [CoinInfo](#cloud.hashing.apis.v1.CoinInfo) | repeated |  |
| Total | [int32](#int32) |  |  |
| Extra | [GoodExtraInfo](#cloud.hashing.apis.v1.GoodExtraInfo) |  |  |
| Title | [string](#string) |  |  |
| Unit | [string](#string) |  |  |
| Start | [uint32](#uint32) |  |  |
| PriceCurrency | [PriceCurrency](#cloud.hashing.apis.v1.PriceCurrency) |  |  |
| Fees | [Fee](#cloud.hashing.apis.v1.Fee) | repeated |  |






<a name="cloud.hashing.apis.v1.GoodExtraInfo"></a>

### GoodExtraInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ID | [string](#string) |  |  |
| GoodID | [string](#string) |  |  |
| Posters | [string](#string) | repeated |  |
| Labels | [string](#string) | repeated |  |
| OutSale | [bool](#bool) |  |  |
| PreSale | [bool](#bool) |  |  |
| VoteCount | [uint32](#uint32) |  |  |
| Rating | [float](#float) |  |  |






<a name="cloud.hashing.apis.v1.GoodFee"></a>

### GoodFee



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ID | [string](#string) |  |  |
| AppID | [string](#string) |  |  |
| FeeType | [string](#string) |  |  |
| FeeDescription | [string](#string) |  |  |
| PayType | [string](#string) |  |  |






<a name="cloud.hashing.apis.v1.GoodInfo"></a>

### GoodInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ID | [string](#string) |  |  |
| DeviceInfoID | [string](#string) |  |  |
| SeparateFee | [bool](#bool) |  |  |
| UnitPower | [int32](#int32) |  |  |
| DurationDays | [int32](#int32) |  |  |
| CoinInfoID | [string](#string) |  |  |
| Actuals | [bool](#bool) |  |  |
| DeliveryAt | [int32](#int32) |  |  |
| InheritFromGoodID | [string](#string) |  |  |
| VendorLocationID | [string](#string) |  |  |
| Price | [double](#double) |  |  |
| BenefitType | [string](#string) |  |  |
| Classic | [bool](#bool) |  |  |
| SupportCoinTypeIDs | [string](#string) | repeated |  |
| Total | [int32](#int32) |  |  |
| PriceCurrency | [string](#string) |  |  |
| Title | [string](#string) |  |  |
| Unit | [string](#string) |  |  |
| Start | [uint32](#uint32) |  |  |






<a name="cloud.hashing.apis.v1.GoodPaying"></a>

### GoodPaying



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ID | [string](#string) |  |  |
| OrderID | [string](#string) |  |  |
| PaymentID | [string](#string) |  |  |






<a name="cloud.hashing.apis.v1.Order"></a>

### Order



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ID | [string](#string) |  |  |
| GoodID | [string](#string) |  |  |
| AppID | [string](#string) |  |  |
| UserID | [string](#string) |  |  |
| Units | [uint32](#uint32) |  |  |
| DiscountCouponID | [string](#string) |  |  |
| UserSpecialReductionID | [string](#string) |  |  |
| Start | [uint32](#uint32) |  |  |
| End | [uint32](#uint32) |  |  |
| CouponID | [string](#string) |  |  |






<a name="cloud.hashing.apis.v1.OrderDetail"></a>

### OrderDetail



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ID | [string](#string) |  |  |
| GoodID | [string](#string) |  |  |
| AppID | [string](#string) |  |  |
| UserID | [string](#string) |  |  |
| Units | [uint32](#uint32) |  |  |
| DiscountCoupon | [Discount](#cloud.hashing.apis.v1.Discount) |  |  |
| UserSpecialReduction | [UserSpecialReduction](#cloud.hashing.apis.v1.UserSpecialReduction) |  |  |
| GoodPaying | [GoodPaying](#cloud.hashing.apis.v1.GoodPaying) |  |  |
| GasPayings | [GasPaying](#cloud.hashing.apis.v1.GasPaying) | repeated |  |
| Compensates | [Compensate](#cloud.hashing.apis.v1.Compensate) | repeated |  |
| OutOfGases | [OutOfGas](#cloud.hashing.apis.v1.OutOfGas) | repeated |  |
| Payment | [Payment](#cloud.hashing.apis.v1.Payment) |  |  |
| Start | [uint32](#uint32) |  |  |
| End | [uint32](#uint32) |  |  |
| Coupon | [Coupon](#cloud.hashing.apis.v1.Coupon) |  |  |
| Discount | [uint32](#uint32) |  |  |
| SpecialReductionAmount | [double](#double) |  |  |






<a name="cloud.hashing.apis.v1.OrderFee"></a>

### OrderFee



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ID | [string](#string) |  |  |
| DurationDays | [uint32](#uint32) |  |  |






<a name="cloud.hashing.apis.v1.OutOfGas"></a>

### OutOfGas



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ID | [string](#string) |  |  |
| OrderID | [string](#string) |  |  |
| Start | [uint32](#uint32) |  |  |
| End | [uint32](#uint32) |  |  |






<a name="cloud.hashing.apis.v1.PageInfo"></a>

### PageInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| PageIndex | [int32](#int32) |  |  |
| PageSize | [int32](#int32) |  |  |






<a name="cloud.hashing.apis.v1.Payment"></a>

### Payment



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ID | [string](#string) |  |  |
| OrderID | [string](#string) |  |  |
| Account | [Account](#cloud.hashing.apis.v1.Account) |  |  |
| Amount | [double](#double) |  |  |
| CoinInfo | [CoinInfo](#cloud.hashing.apis.v1.CoinInfo) |  |  |
| State | [string](#string) |  |  |
| ChainTransactionID | [string](#string) |  |  |
| PlatformTransactionID | [string](#string) |  |  |






<a name="cloud.hashing.apis.v1.PriceCurrency"></a>

### PriceCurrency



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ID | [string](#string) |  |  |
| Name | [string](#string) |  |  |
| Unit | [string](#string) |  |  |
| Symbol | [string](#string) |  |  |






<a name="cloud.hashing.apis.v1.SubmitOrderRequest"></a>

### SubmitOrderRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| GoodID | [string](#string) |  |  |
| Units | [uint32](#uint32) |  |  |
| UserID | [string](#string) |  |  |
| AppID | [string](#string) |  |  |
| CouponID | [string](#string) |  |  |






<a name="cloud.hashing.apis.v1.SubmitOrderResponse"></a>

### SubmitOrderResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Info | [Order](#cloud.hashing.apis.v1.Order) |  |  |






<a name="cloud.hashing.apis.v1.TargetAreaInfo"></a>

### TargetAreaInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ID | [string](#string) |  |  |
| Continent | [string](#string) |  |  |
| Country | [string](#string) |  |  |






<a name="cloud.hashing.apis.v1.UserSpecialReduction"></a>

### UserSpecialReduction



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ID | [string](#string) |  |  |
| AppID | [string](#string) |  |  |
| UserID | [string](#string) |  |  |
| Amount | [double](#double) |  |  |
| ReleaseByUserID | [string](#string) |  |  |
| Start | [uint32](#uint32) |  |  |
| DurationDays | [int32](#int32) |  |  |
| Message | [string](#string) |  |  |






<a name="cloud.hashing.apis.v1.VendorLocationInfo"></a>

### VendorLocationInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ID | [string](#string) |  |  |
| Country | [string](#string) |  |  |
| Province | [string](#string) |  |  |
| City | [string](#string) |  |  |
| Address | [string](#string) |  |  |






<a name="cloud.hashing.apis.v1.VersionResponse"></a>

### VersionResponse
Request body and response


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Info | [string](#string) |  |  |





 

 

 


<a name="cloud.hashing.apis.v1.CloudHashingApis"></a>

### CloudHashingApis
Cloud Hashing Goods

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Version | [.google.protobuf.Empty](#google.protobuf.Empty) | [VersionResponse](#cloud.hashing.apis.v1.VersionResponse) |  |
| GetGoodsDetail | [GetGoodsDetailRequest](#cloud.hashing.apis.v1.GetGoodsDetailRequest) | [GetGoodsDetailResponse](#cloud.hashing.apis.v1.GetGoodsDetailResponse) |  |
| GetGoodDetail | [GetGoodDetailRequest](#cloud.hashing.apis.v1.GetGoodDetailRequest) | [GetGoodDetailResponse](#cloud.hashing.apis.v1.GetGoodDetailResponse) |  |
| SubmitOrder | [SubmitOrderRequest](#cloud.hashing.apis.v1.SubmitOrderRequest) | [SubmitOrderResponse](#cloud.hashing.apis.v1.SubmitOrderResponse) |  |
| CreateOrderPayment | [CreateOrderPaymentRequest](#cloud.hashing.apis.v1.CreateOrderPaymentRequest) | [CreateOrderPaymentResponse](#cloud.hashing.apis.v1.CreateOrderPaymentResponse) |  |
| GetOrderDetail | [GetOrderDetailRequest](#cloud.hashing.apis.v1.GetOrderDetailRequest) | [GetOrderDetailResponse](#cloud.hashing.apis.v1.GetOrderDetailResponse) |  |
| GetOrdersDetailByAppUser | [GetOrdersDetailByAppUserRequest](#cloud.hashing.apis.v1.GetOrdersDetailByAppUserRequest) | [GetOrdersDetailByAppUserResponse](#cloud.hashing.apis.v1.GetOrdersDetailByAppUserResponse) |  |
| GetOrdersDetailByApp | [GetOrdersDetailByAppRequest](#cloud.hashing.apis.v1.GetOrdersDetailByAppRequest) | [GetOrdersDetailByAppResponse](#cloud.hashing.apis.v1.GetOrdersDetailByAppResponse) |  |
| GetOrdersDetailByGood | [GetOrdersDetailByGoodRequest](#cloud.hashing.apis.v1.GetOrdersDetailByGoodRequest) | [GetOrdersDetailByGoodResponse](#cloud.hashing.apis.v1.GetOrdersDetailByGoodResponse) |  |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

