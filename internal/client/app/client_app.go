package client

import (
	grpcClient "github.com/Ssnakerss/mypreciouskeeper/internal/client/grpc"
)

var App *ClietApp

type ClietApp struct {
	GRPC      *grpcClient.GRPCClient
	AuthToken string
	UserID    int64
	UserName  string
}

func NewClientApp(
	gRPCAddress string,
) *ClietApp {
	return &ClietApp{
		GRPC: grpcClient.NewGRPCClient(gRPCAddress),
	}
}
