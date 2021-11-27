package api

import (
	"context"

	"github.com/NpoolPlatform/cloud-hashing-apis/message/npool"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	npool.UnimplementedCloudHashingApisServer
}

func Register(server grpc.ServiceRegistrar) {
	npool.RegisterCloudHashingApisServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return npool.RegisterCloudHashingApisHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
