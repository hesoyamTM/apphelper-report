package psql

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hesoyamTM/apphelper-report/internal/models"
	"github.com/hesoyamTM/apphelper-sso/pkg/logger"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(host, user, password, db string, port int) *Storage {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, db)

	conn, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		panic(err)
	}

	return &Storage{
		db: conn,
	}
}

func (s *Storage) CreateReport(ctx context.Context, groupId, studentId, trainerId uuid.UUID, desc string) error {
	const op = "psql.CreateReport"

	query := `INSERT INTO reports (student_id, trainer_id, group_id, description, date) VALUES ($1, $2, $3, $4, now())`

	if _, err := s.db.Exec(ctx, query, studentId, trainerId, groupId, desc); err != nil {
		// TODO: refactor error

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) ProvideReport(ctx context.Context, groupId, studentId, trainerId uuid.UUID) ([]models.Report, error) {
	const op = "psql.ProvideReport"

	if groupId == uuid.Nil && studentId == uuid.Nil && trainerId == uuid.Nil {
		return nil, nil
	}

	q, i1, i2, i3 := parseProvideReportsCondition(groupId, studentId, trainerId)

	query := fmt.Sprintf(`SELECT group_id, student_id, trainer_id, description, date FROM reports WHERE %s`, q)
	log := logger.GetLoggerFromCtx(ctx)
	log.Info(ctx, "query", zap.String("query", query))

	args := make([]interface{}, 0)
	if i1 {
		args = append(args, groupId)
	}
	if i2 {
		args = append(args, studentId)
	}
	if i3 {
		args = append(args, trainerId)
	}

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	reports := make([]models.Report, 0)

	for rows.Next() {
		var rep models.Report
		if err := rows.Scan(&rep.GroupId, &rep.StudentId, &rep.TrainerId, &rep.Description, &rep.Date); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		reports = append(reports, rep)
	}
	return reports, nil
}

func parseProvideReportsCondition(groupId, studentId, trainerId uuid.UUID) (query string, i1 bool, i2 bool, i3 bool) {

	query = ""
	i := 1
	i1, i2, i3 = false, false, false

	if groupId != uuid.Nil {
		if query != "" {
			query += " AND"
		}
		query += fmt.Sprintf(" group_id = $%d", i)
		i++
		i1 = true
	}
	if studentId != uuid.Nil {
		if query != "" {
			query += " AND"
		}
		query += fmt.Sprintf(" student_id = $%d", i)
		i++
		i2 = true
	}
	if trainerId != uuid.Nil {
		if query != "" {
			query += " AND"
		}
		query += fmt.Sprintf(" trainer_id = $%d", i)
		i++
		i3 = true
	}

	return
}
