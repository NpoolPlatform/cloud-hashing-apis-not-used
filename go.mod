module github.com/NpoolPlatform/cloud-hashing-apis

go 1.16

require (
	entgo.io/ent v0.10.0
	github.com/NpoolPlatform/api-manager v0.0.0-20220205130236-69d286e72dba
	github.com/NpoolPlatform/appuser-manager v0.0.0-20220210093932-4b9db1361d89
	github.com/NpoolPlatform/cloud-hashing-billing v0.0.0-20220214123916-517341a90b77
	github.com/NpoolPlatform/cloud-hashing-goods v0.0.0-20220113121137-c2b65b514bad
	github.com/NpoolPlatform/cloud-hashing-inspire v0.0.0-20220113121537-8e8b6c8966da
	github.com/NpoolPlatform/cloud-hashing-order v0.0.0-20220211133153-5c577f9d3988
	github.com/NpoolPlatform/cloud-hashing-staker v0.0.0-20220214050218-567efaad1cb3
	github.com/NpoolPlatform/go-service-framework v0.0.0-20220211051615-b2300d03022a
	github.com/NpoolPlatform/kyc-management v0.0.0-20220113122339-4bef8bdcbc5c
	github.com/NpoolPlatform/message v0.0.0-20220215031642-df3fb9a755e7
	github.com/NpoolPlatform/review-service v0.0.0-20220214135408-eb1dbe09de65
	github.com/NpoolPlatform/sphinx-coininfo v0.0.0-20211208035009-5ad2768d2290
	github.com/NpoolPlatform/sphinx-proxy v0.0.0-20211210102925-d9b8abe11021
	github.com/NpoolPlatform/third-gateway v0.0.0-20220204065515-e8e92a3fce48
	github.com/go-resty/resty/v2 v2.7.0
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.7.2
	github.com/streadway/amqp v1.0.0
	github.com/stretchr/testify v1.7.1-0.20210427113832-6241f9ab9942
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
	google.golang.org/grpc v1.44.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.2.0
	google.golang.org/protobuf v1.27.1
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.41.0
