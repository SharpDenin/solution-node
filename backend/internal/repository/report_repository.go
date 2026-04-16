package repository

import (
	"context"
	"github.com/google/uuid"
)

type ReportRepository interface {
	CreateReport(ctx context.Context, tx Tx, userID uuid.UUID, place string, reportDate string, responsibleName string) (uuid.UUID, error)
	CreateAnswer(ctx context.Context, tx Tx, reportID uuid.UUID, questionID uuid.UUID, answer_text string, imageURL *string) error
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
	answer_text string,
	imageURL *string,
) error {

	query := `
		INSERT INTO answers (report_id, question_id, answer_text, image_url)
		VALUES ($1, $2, $3, $4)
	`

	_, err := tx.Exec(ctx, query, reportID, questionID, answer_text, imageURL)
	return err
}
