package server

import (
	"context"

	"github.com/Ssnakerss/mypreciouskeeper/proto/gen/grpcauth"
	"google.golang.org/grpc"
)

type serverAPI struct {
	grpcauth.UnimplementedAuthServer
	a Auth
}

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		pass string,
	) (token string, err error)
	Register(
		ctx context.Context,
		email string,
		pass string,
	) (userID int64, err error)
}

func Register(gRPCServer *grpc.Server, a Auth) {
	grpcauth.RegisterAuthServer(gRPCServer, &serverAPI{a: a})
}
