package server

import (
	"context"
	"errors"

	"github.com/Ssnakerss/mypreciouskeeper/internal/apperrs"

	grpcserver "github.com/Ssnakerss/mypreciouskeeper/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RegisterGRPC(gRPCServer *grpc.Server, a Auth) {
	grpcserver.RegisterAuthServer(gRPCServer, &serverAPI{auth: a})
}

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		pass string,
	) (token string, err error)

	RegisterUser(
		ctx context.Context,
		email string,
		pass string,
	) (userID int64, err error)
}

type serverAPI struct {
	grpcserver.UnimplementedAuthServer
	auth Auth
}

// Login serve gRPC calls -  check email and password, call app Login func and return token
func (s *serverAPI) Login(
	ctx context.Context,
	in *grpcserver.LoginRequest,
) (*grpcserver.LoginResponse, error) {
	if in.Email == "" || in.Pass == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password required")
	}
	token, err := s.auth.Login(ctx, in.Email, in.Pass)
	if err != nil {
		if errors.Is(err, apperrs.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &grpcserver.LoginResponse{Token: token}, nil
}

// Register serve gRPC calls - check email and password, call app Register func and return userID
func (s *serverAPI) Register(
	ctx context.Context,
	in *grpcserver.RegisterRequest,
) (*grpcserver.RegisterResponse, error) {
	//TODO
	if in.Email == "" || in.Pass == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password required")
	}
	userID, err := s.auth.RegisterUser(ctx, in.Email, in.Pass)
	if err != nil {
		if errors.Is(err, apperrs.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &grpcserver.RegisterResponse{UserId: userID}, nil
}
