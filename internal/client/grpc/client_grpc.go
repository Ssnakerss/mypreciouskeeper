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
		//TODO: Asset create error: rpc error: code = ResourceExhausted desc = grpc: received message larger than max (40123519 vs. 4194304)
		// Why MaxCall ne rabotaet???
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(1024*1024*50),
			grpc.MaxCallSendMsgSize(1024*1024*50),
		), //50mb message size
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
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

		Conn: Conn,
	}
}

func (c *GRPCClient) Close() {
	if c.Conn != nil {
		c.Conn.Close()
	}
}
