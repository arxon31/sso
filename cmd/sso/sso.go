package main //nolint:all

import (
	"context"
	ssogrpc "github.com/arxon31/sso/internal/controller/grpc"
	"github.com/arxon31/sso/internal/repo/postgres"
	"github.com/arxon31/sso/internal/service/auth"
	"github.com/arxon31/sso/internal/service/register"
	"github.com/arxon31/sso/pkg/pgconn"
	"github.com/arxon31/yapr-proto/pkg/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/arxon31/sso/pkg/logger"
)

const port = ":8081"

var Build string

func main() {
	logger.Logger.Info("starting app", slog.String("build", Build))
	err := run()
	if err != nil {
		logger.Logger.Error("app exited with error", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Logger.Info("app exited without errors")
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	postgresConn, err := pgconn.New()
	if err != nil {
		return err
	}
	defer postgresConn.Close()

	repo, err := postgres.NewPostgres(postgresConn)
	if err != nil {
		return err
	}

	authService := auth.NewService(repo)

	registerService := register.NewService(repo)

	ssoController := ssogrpc.NewController(authService, registerService)

	grpcServer := grpc.NewServer()
	defer grpcServer.GracefulStop()

	reflection.Register(grpcServer)

	sso.RegisterSSOServer(grpcServer, ssoController)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	errChan := make(chan error, 1)

	go func() {
		errChan <- grpcServer.Serve(listener)
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return nil
	}

}
