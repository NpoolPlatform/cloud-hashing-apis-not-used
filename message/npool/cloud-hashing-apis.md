# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [npool/cloud-hashing-apis.proto](#npool/cloud-hashing-apis.proto)
    - [CoinInfo](#cloud.hashing.apis.v1.CoinInfo)
    - [DeviceInfo](#cloud.hashing.apis.v1.DeviceInfo)
    - [GetGoodsDetailRequest](#cloud.hashing.apis.v1.GetGoodsDetailRequest)
    - [GetGoodsDetailResponse](#cloud.hashing.apis.v1.GetGoodsDetailResponse)
    - [GoodDetail](#cloud.hashing.apis.v1.GoodDetail)
    - [GoodExtraInfo](#cloud.hashing.apis.v1.GoodExtraInfo)
    - [GoodInfo](#cloud.hashing.apis.v1.GoodInfo)
    - [PageInfo](#cloud.hashing.apis.v1.PageInfo)
    - [PriceCurrency](#cloud.hashing.apis.v1.PriceCurrency)
    - [TargetAreaInfo](#cloud.hashing.apis.v1.TargetAreaInfo)
    - [VendorLocationInfo](#cloud.hashing.apis.v1.VendorLocationInfo)
    - [VersionResponse](#cloud.hashing.apis.v1.VersionResponse)
  
    - [CloudHashingApis](#cloud.hashing.apis.v1.CloudHashingApis)
  
- [Scalar Value Types](#scalar-value-types)



<a name="npool/cloud-hashing-apis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## npool/cloud-hashing-apis.proto



<a name="cloud.hashing.apis.v1.CoinInfo"></a>

### CoinInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ID | [string](#string) |  |  |
| PreSale | [bool](#bool) |  |  |
| Name | [string](#string) |  |  |
| Unit | [string](#string) |  |  |
| Logo | [string](#string) |  |  |






<a name="cloud.hashing.apis.v1.DeviceInfo"></a>

### DeviceInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ID | [string](#string) |  |  |
| Type | [string](#string) |  |  |
| Manufacturer | [string](#string) |  |  |
| PowerComsuption | [int32](#int32) |  |  |
| ShipmentAt | [int32](#int32) |  |  |






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






<a name="cloud.hashing.apis.v1.PageInfo"></a>

### PageInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| PageIndex | [int32](#int32) |  |  |
| PageSize | [int32](#int32) |  |  |






<a name="cloud.hashing.apis.v1.PriceCurrency"></a>

### PriceCurrency



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ID | [string](#string) |  |  |
| Name | [string](#string) |  |  |
| Unit | [string](#string) |  |  |
| Symbol | [string](#string) |  |  |






<a name="cloud.hashing.apis.v1.TargetAreaInfo"></a>

### TargetAreaInfo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ID | [string](#string) |  |  |
| Continent | [string](#string) |  |  |
| Country | [string](#string) |  |  |






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

