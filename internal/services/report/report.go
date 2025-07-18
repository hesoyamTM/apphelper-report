package report

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	"github.com/google/uuid"
	"github.com/hesoyamTM/apphelper-report/internal/models"

	"github.com/hesoyamTM/apphelper-sso/pkg/logger"
	"go.uber.org/zap"
)

type ReportStorage interface {
	CreateReport(ctx context.Context, groupId, studentId, trainerId uuid.UUID, desc string) error
	ProvideReport(ctx context.Context, groupId, studentId, trainerId uuid.UUID) ([]models.Report, error)
}

type Report struct {
	db ReportStorage

	keyCh chan *ecdsa.PublicKey
}

func New(ctx context.Context, db ReportStorage) *Report {
	keyCh := make(chan *ecdsa.PublicKey, 1)

	return &Report{
		db:    db,
		keyCh: keyCh,
	}
}

func (r *Report) CreateReport(ctx context.Context, groupId, studentId, trainerId uuid.UUID, desc string) error {
	const op = "report.CreateReport"
	log := logger.GetLoggerFromCtx(ctx)

	if err := r.db.CreateReport(ctx, groupId, studentId, trainerId, desc); err != nil {
		log.Error(ctx, "failed to create report", zap.Error(err))

		//TODO: refactor error

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Report) GetReports(ctx context.Context, groupId, student_id, trainer_id uuid.UUID) ([]models.Report, error) {
	const op = "report.GetReport"
	log := logger.GetLoggerFromCtx(ctx)

	rep, err := r.db.ProvideReport(ctx, groupId, student_id, trainer_id)
	if err != nil {
		log.Error(ctx, "failed to get report", zap.Error(err))

		//TODO: refactor error

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return rep, nil
}
