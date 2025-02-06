package server

import (
	"context"
	"time"

	grpcserver "github.com/Ssnakerss/mypreciouskeeper/proto/gen"
	"google.golang.org/grpc"
)

type serverPingAPI struct {
	grpcserver.UnimplementedPingServer
}

func NewServerPingAPI() *serverPingAPI {
	return &serverPingAPI{}
}

func (p *serverPingAPI) RegisterGRPC(gRPCServer *grpc.Server) {
	grpcserver.RegisterPingServer(gRPCServer, p)
}

func (p *serverPingAPI) Ping(ctx context.Context, in *grpcserver.PingRequest) (*grpcserver.PingResponse, error) {
	return &grpcserver.PingResponse{
		Resp: time.Now().Unix(),
	}, nil
}
