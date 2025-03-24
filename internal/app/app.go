package app

import (
	"context"

	"github.com/hesoyamTM/apphelper-report/internal/app/grpcapp"
	"github.com/hesoyamTM/apphelper-report/internal/services/report"
	"github.com/hesoyamTM/apphelper-report/internal/storage/psql"
)

type App struct {
	GRPCApp grpcapp.App
}

type GrpcOpts struct {
	Host string
	Port int
}

type PsqlOpts struct {
	Host string
	Port int
	User string
	Pass string
	DB   string
}

func New(ctx context.Context, grpcOpts GrpcOpts, psqlOpts PsqlOpts) *App {
	storage := psql.New(psqlOpts.Host, psqlOpts.User, psqlOpts.Pass, psqlOpts.DB, psqlOpts.Port)

	reportService := report.New(ctx, storage)

	grpcApp := grpcapp.New(
		ctx,
		grpcOpts.Host,
		grpcOpts.Port,
		reportService,
	)

	return &App{
		GRPCApp: *grpcApp,
	}
}
