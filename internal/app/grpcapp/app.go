package grpcapp

import (
	"context"
	"fmt"
	"net"

	"github.com/hesoyamTM/apphelper-report/internal/grpc/report"

	"github.com/hesoyamTM/apphelper-sso/pkg/logger"
	"google.golang.org/grpc"
)

type App struct {
	grpcServer *grpc.Server
	host       string
	port       int
}

func New(ctx context.Context,
	host string,
	port int,
	repServ report.Report,
) *App {
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(logger.LoggingInterceptor(ctx)))

	report.RegisterServer(grpcServer, repServ)

	return &App{
		host:       host,
		port:       port,
		grpcServer: grpcServer,
	}
}

func (a *App) MustRun(ctx context.Context) {
	if err := a.run(ctx); err != nil {
		panic(err)
	}
}

func (a *App) run(ctx context.Context) error {
	log := logger.GetLoggerFromCtx(ctx)

	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.host, a.port))
	if err != nil {
		return err
	}

	log.Info(ctx, "grpc server is running")

	if err = a.grpcServer.Serve(l); err != nil {
		return err
	}

	return nil
}

func (a *App) Stop(ctx context.Context) {
	log := logger.GetLoggerFromCtx(ctx)

	log.Info(ctx, "grpc server is stopping")

	a.grpcServer.GracefulStop()
}
