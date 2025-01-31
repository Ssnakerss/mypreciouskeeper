package client

import (
	grpcClient "github.com/Ssnakerss/mypreciouskeeper/internal/client/grpc"
)

var App *ClientApp

type ClientApp struct {
	GRPC      *grpcClient.GRPCClient
	AuthToken string
	UserID    int64
	UserName  string
}

func NewClientApp(
	gRPCAddress string,
) *ClientApp {
	return &ClientApp{
		GRPC: grpcClient.NewGRPCClient(gRPCAddress),
	}
}
