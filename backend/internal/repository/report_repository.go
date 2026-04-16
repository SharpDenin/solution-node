package repository

import (
	"backend/internal/models/dtos"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type ReportFilters struct {
	DateFrom *string
	DateTo   *string
	Place    *string
	UserID   *string
	Limit    int
	Offset   int
}

type ReportRepository interface {
	CreateReport(ctx context.Context, tx Tx, userID uuid.UUID, place string, reportDate string, responsibleName string) (uuid.UUID, error)
	CreateAnswer(ctx context.Context, tx Tx, reportID uuid.UUID, questionID uuid.UUID, answerText string, imageURL *string) error
	GetReports(ctx context.Context, filters ReportFilters) ([]dtos.ReportResponse, error)
}

type reportRepository struct {
	db *DB
}

func NewReportRepository(db *DB) ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) CreateReport(
	ctx context.Context,
	tx Tx,
	userID uuid.UUID,
	place string,
	reportDate string,
	responsibleName string,
) (uuid.UUID, error) {

	query := `
		INSERT INTO reports (user_id, place, report_date, responsible_name)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var reportID uuid.UUID

	err := tx.QueryRow(ctx, query,
		userID,
		place,
		reportDate,
		responsibleName,
	).Scan(&reportID)

	if err != nil {
		return uuid.Nil, err
	}

	return reportID, nil
}

func (r *reportRepository) CreateAnswer(
	ctx context.Context,
	tx Tx,
	reportID uuid.UUID,
	questionID uuid.UUID,
	answerText string,
	imageURL *string,
) error {

	query := `
		INSERT INTO answers (report_id, question_id, answer_text, image_url)
		VALUES ($1, $2, $3, $4)
	`

	_, err := tx.Exec(ctx, query, reportID, questionID, answerText, imageURL)
	return err
}

func (r *reportRepository) GetReports(ctx context.Context, f ReportFilters) ([]dtos.ReportResponse, error) {
	query := `
		SELECT id, user_id, place, report_date, responsible_name, created_at
		FROM reports
		WHERE 1=1
	`

	args := []interface{}{}
	argID := 1

	if f.DateFrom != nil {
		query += " AND report_date >= $" + fmt.Sprint(argID)
		args = append(args, *f.DateFrom)
		argID++
	}

	if f.DateTo != nil {
		query += " AND report_date <= $" + fmt.Sprint(argID)
		args = append(args, *f.DateTo)
		argID++
	}

	if f.Place != nil {
		query += " AND place ILIKE $" + fmt.Sprint(argID)
		args = append(args, "%"+*f.Place+"%")
		argID++
	}

	if f.UserID != nil {
		query += " AND user_id = $" + fmt.Sprint(argID)
		args = append(args, *f.UserID)
		argID++
	}

	query += " ORDER BY report_date DESC"

	query += " LIMIT $" + fmt.Sprint(argID)
	args = append(args, f.Limit)
	argID++

	query += " OFFSET $" + fmt.Sprint(argID)
	args = append(args, f.Offset)

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []dtos.ReportResponse

	for rows.Next() {
		var r dtos.ReportResponse

		err := rows.Scan(
			&r.ID,
			&r.UserID,
			&r.Place,
			&r.ReportDate,
			&r.ResponsibleName,
			&r.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		reports = append(reports, r)
	}

	return reports, nil
}
