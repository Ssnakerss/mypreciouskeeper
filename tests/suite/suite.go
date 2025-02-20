package suite

import (
	"context"
	"net"
	"testing"
	"time"

	grpcserver "github.com/Ssnakerss/mypreciouskeeper/proto/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T
	// Cfg     *config.Config
	AClient     grpcserver.AuthClient
	AssetClient grpcserver.AssetClient
	PingClient  grpcserver.PingClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	// cfg := config.Load()

	ctx, cancelCtx := context.WithTimeout(context.Background(), time.Second*5) // cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	// Адрес нашего gRPC-сервера
	grpcAddress := net.JoinHostPort("localhost", "44044") // strconv.Itoa(cfg.GRPC.Port))

	// Создаем клиент
	cc, err := grpc.NewClient(
		grpcAddress,
		// Используем insecure-коннект для тестов
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	authClient := grpcserver.NewAuthClient(cc)
	assetClient := grpcserver.NewAssetClient(cc)
	pingClient := grpcserver.NewPingClient(cc)

	return ctx, &Suite{
		T: t,
		// Cfg:     cfg,
		AClient:     authClient,
		AssetClient: assetClient,
		PingClient:  pingClient,
	}

}
