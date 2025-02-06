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

// AUth interface for authorization business logic
type AuthService interface {
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

type serverAuthAPI struct {
	grpcserver.UnimplementedAuthServer
	authService AuthService
}

func NewServerAuthAPI(a AuthService) *serverAuthAPI {
	return &serverAuthAPI{
		authService: a,
	}
}

func (s *serverAuthAPI) RegisterGRPC(gRPCServer *grpc.Server) {
	grpcserver.RegisterAuthServer(gRPCServer, s)
}

// Login serve gRPC calls -  check email and password, call app Login func and return token
func (s *serverAuthAPI) Login(
	ctx context.Context,
	in *grpcserver.LoginRequest,
) (*grpcserver.LoginResponse, error) {
	if in.Email == "" || in.Pass == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password required")
	}
	token, err := s.authService.Login(ctx, in.Email, in.Pass)
	if err != nil {
		if errors.Is(err, apperrs.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &grpcserver.LoginResponse{Token: token}, nil
}

// Register serve gRPC calls - check email and password, call app Register func and return userID
func (s *serverAuthAPI) Register(
	ctx context.Context,
	in *grpcserver.RegisterRequest,
) (*grpcserver.RegisterResponse, error) {
	//TODO
	if in.Email == "" || in.Pass == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password required")
	}
	userID, err := s.authService.Register(ctx, in.Email, in.Pass)
	if err != nil {
		if errors.Is(err, apperrs.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &grpcserver.RegisterResponse{UserId: userID}, nil
}

// verifyJWTPayload verify JWT token and extract User info
