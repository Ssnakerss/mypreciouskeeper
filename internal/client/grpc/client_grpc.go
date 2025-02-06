package grpcClient

import (
	"log"

	grpcserver "github.com/Ssnakerss/mypreciouskeeper/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	AuthClient  grpcserver.AuthClient
	AssetClient grpcserver.AssetClient
	PingClient  grpcserver.PingClient
	token       string
	Conn        *grpc.ClientConn
}

// NewGRPCClient create client with Auth and Asset Endpoints from gRPC server
func NewGRPCClient(grpcAddress string) *GRPCClient {
	Conn, err := grpc.NewClient(
		grpcAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("grpc server connection failed: %v", err)
	}
	authClient := grpcserver.NewAuthClient(Conn)
	assetClient := grpcserver.NewAssetClient(Conn)
	pingClient := grpcserver.NewPingClient(Conn)

	return &GRPCClient{
		AuthClient:  authClient,
		AssetClient: assetClient,
		PingClient:  pingClient,
	}
}
