package server

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"time"

	"github.com/Ssnakerss/mypreciouskeeper/internal/logger"
	"github.com/Ssnakerss/mypreciouskeeper/internal/services/auth"
	"github.com/Ssnakerss/mypreciouskeeper/internal/storage"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpclogging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	l    *slog.Logger
	gRPC *grpc.Server
	port int
}

// New creates an instance of server with logger and gRPC server
func New(l *slog.Logger, port int) *Server {

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			l.Error("panic", slog.Any("panic", p))
			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	loggingOpts := []logging.Option{
		grpclogging.WithLogOnEvents(
			grpclogging.PayloadReceived, grpclogging.PayloadSent,
		),
	}

	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		grpclogging.UnaryServerInterceptor(logger.InterceptorLogger(l), loggingOpts...),
	))

	dsn := "postgres://orchestra:orchestra12qwaszx@pg-ext.os.serk.lan:5103/orchestra?sslmode=disable"
	db, err := storage.New(context.Background(), dsn, time.Second*3)
	if err != nil {
		log.Fatal("db connection failed: ", err)
	}

	a := auth.New(l, db, time.Second*3)

	RegisterGRPC(gRPCServer, *a)

	return &Server{
		l:    l,
		gRPC: gRPCServer,
		port: port,
	}
}

// Run start listener on a specified port
func (s *Server) Run() error {

	ls, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	s.l.Info("server started", slog.String("addr", ls.Addr().String()))

	if err := s.gRPC.Serve(ls); err != nil {
		return fmt.Errorf("%s : %w", ls.Addr().String(), err)
	}
	return nil
}

// MustRun start listener on a specified port and panic on error
func (s *Server) MustRun() {
	if err := s.Run(); err != nil {
		log.Fatal("server start failed", err)
	}
}

func (s *Server) Shutdown() {
	s.l.Info("server shuting down ....")
	s.gRPC.GracefulStop()
	s.l.Info("server shutdown complete, exit")
}
