package psql

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/hesoyamTM/apphelper-sso/pkg/logger"
	"go.uber.org/zap"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrate(
	ctx context.Context,
	host, user, password, db string,
) error {
	const op = "psql.RunMigrate"
	log := logger.GetLoggerFromCtx(ctx)

	m, err := migrate.New(
		"file://migrations",
		fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, host, db),
	)
	if err != nil {
		log.Error(ctx, "failed to create migrate", zap.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info(ctx, "running migrations")

	err = m.Up()
	if err != nil {
		if err != migrate.ErrNoChange {
			log.Error(ctx, "failed to migrate", zap.Error(err))
			return fmt.Errorf("%s: %w", op, err)
		} else {
			log.Info(ctx, "no change")
		}
	}

	log.Info(ctx, "migrations done")

	return nil
}
