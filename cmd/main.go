package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	app "github.com/hesoyamTM/apphelper-report/internal/app"
	"github.com/hesoyamTM/apphelper-report/internal/config"

	"github.com/hesoyamTM/apphelper-sso/pkg/logger"
)

func main() {
	ctx := context.Background()
	cfg := config.MustLoad()

	ctx, err := logger.New(ctx, cfg.Env)
	if err != nil {
		panic(err)
	}

	log := logger.GetLoggerFromCtx(ctx)
	log.Debug(ctx, "logger is working")

	grpcOpts := app.GrpcOpts{
		Host: cfg.Grpc.Host,
		Port: cfg.Grpc.Port,
	}

	psqlOpts := app.PsqlOpts{
		Host: cfg.Psql.Host,
		User: cfg.Psql.User,
		Port: cfg.Psql.Port,
		Pass: cfg.Psql.Password,
		DB:   cfg.Psql.DB,
	}

	application := app.New(ctx, grpcOpts, psqlOpts)
	go application.GRPCApp.MustRun(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	application.GRPCApp.Stop(ctx)
	log.Info(ctx, "application stopped")
}
