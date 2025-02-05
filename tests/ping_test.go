package tests

import (
	"net"
	"testing"

	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Test_Ping(t *testing.T) {
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
	for i := 0; i < 20; i++ {
		cc.Connect()
		t.Log(cc.GetState().String())
		// if cc.GetState().String() != "IDLE" && cc.GetState().String() != "READY" {
		// 	cc.Close()
		// 	t.Fatal()
		// }
		time.Sleep(time.Second * 2)
	}
}
