package repository

import (
	"backend/internal/models/dtos"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ReportFilters struct {
	DateFrom *string
	DateTo   *string
	Place    *string
	UserID   *string
	UserName *string
	Limit    int
	Offset   int
}

type ReportRepository interface {
	CreateReport(ctx context.Context, tx Tx, userID uuid.UUID, place string, reportDate string, responsibleName string) (uuid.UUID, error)
	CreateAnswer(ctx context.Context, tx Tx, reportID uuid.UUID, questionID uuid.UUID, answerText string, imageURL *string) error
	GetReports(ctx context.Context, filters ReportFilters) ([]dtos.ReportResponse, error)
	GetReportByID(ctx context.Context, id string) (*dtos.ReportDetailResponse, error)

	GetReportsDetailed(ctx context.Context, filters ReportFilters) ([]dtos.ReportDetailResponse, error)
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
        SELECT r.id, r.user_id, r.place, r.report_date, r.responsible_name, r.created_at
        FROM reports r
        LEFT JOIN users u ON u.id = r.user_id
        WHERE 1=1
    `

	args := []interface{}{}
	argID := 1

	if f.DateFrom != nil {
		query += " AND r.report_date >= $" + fmt.Sprint(argID)
		args = append(args, *f.DateFrom)
		argID++
	}

	if f.DateTo != nil {
		query += " AND r.report_date <= $" + fmt.Sprint(argID)
		args = append(args, *f.DateTo)
		argID++
	}

	if f.Place != nil {
		query += " AND r.place ILIKE $" + fmt.Sprint(argID)
		args = append(args, "%"+*f.Place+"%")
		argID++
	}

	if f.UserID != nil {
		query += " AND r.user_id = $" + fmt.Sprint(argID)
		args = append(args, *f.UserID)
		argID++
	}

	// НОВЫЙ ФИЛЬТР ПО ИМЕНИ ПОЛЬЗОВАТЕЛЯ
	if f.UserName != nil {
		query += " AND u.full_name ILIKE $" + fmt.Sprint(argID)
		args = append(args, "%"+*f.UserName+"%")
		argID++
	}

	query += " ORDER BY r.report_date DESC"
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
		err := rows.Scan(&r.ID, &r.UserID, &r.Place, &r.ReportDate, &r.ResponsibleName, &r.CreatedAt)
		if err != nil {
			return nil, err
		}
		reports = append(reports, r)
	}
	return reports, nil
}

func (r *reportRepository) GetReportByID(ctx context.Context, id string) (*dtos.ReportDetailResponse, error) {

	query := `
		SELECT 
			r.id,
			r.user_id,
			r.place,
			r.report_date,
			r.responsible_name,
			r.created_at,
			q.id,
			q.text,
			a.answer_text,
			a.image_url
		FROM reports r
		LEFT JOIN answers a ON a.report_id = r.id
		LEFT JOIN questions q ON q.id = a.question_id
		WHERE r.id = $1
		ORDER BY q.order_index
	`

	rows, err := r.db.Pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var report *dtos.ReportDetailResponse

	for rows.Next() {
		var (
			rID, userID, place, responsible string
			reportDate, createdAt           time.Time

			qID, qText *string
			answerText *string
			imageURL   *string
		)

		err := rows.Scan(
			&rID,
			&userID,
			&place,
			&reportDate,
			&responsible,
			&createdAt,
			&qID,
			&qText,
			&answerText,
			&imageURL,
		)
		if err != nil {
			return nil, err
		}

		// создаём report один раз
		if report == nil {
			report = &dtos.ReportDetailResponse{
				ID:              rID,
				UserID:          userID,
				Place:           place,
				ReportDate:      reportDate,
				ResponsibleName: responsible,
				CreatedAt:       createdAt,
				Answers:         []dtos.AnswerResponse{},
			}
		}

		// если есть ответы
		if qID != nil {
			report.Answers = append(report.Answers, dtos.AnswerResponse{
				QuestionID:   *qID,
				QuestionText: *qText,
				AnswerText:   derefString(answerText),
				ImageURL:     imageURL,
			})
		}
	}

	if report == nil {
		return nil, errors.New("report not found")
	}

	return report, nil
}

func (r *reportRepository) GetReportsDetailed(ctx context.Context, f ReportFilters) ([]dtos.ReportDetailResponse, error) {

	query := `
		SELECT 
			r.id,
			r.user_id,
			r.place,
			r.report_date,
			r.responsible_name,
			r.created_at,
			q.id,
			q.text,
			a.answer_text,
			a.image_url
		FROM reports r
		LEFT JOIN answers a ON a.report_id = r.id
		LEFT JOIN questions q ON q.id = a.question_id
		LEFT JOIN users u ON u.id = r.user_id
		WHERE 1=1
	`

	args := []interface{}{}
	argID := 1

	if f.DateFrom != nil {
		query += " AND r.report_date >= $" + fmt.Sprint(argID)
		args = append(args, *f.DateFrom)
		argID++
	}

	if f.DateTo != nil {
		query += " AND r.report_date <= $" + fmt.Sprint(argID)
		args = append(args, *f.DateTo)
		argID++
	}

	if f.Place != nil {
		query += " AND r.place ILIKE $" + fmt.Sprint(argID)
		args = append(args, "%"+*f.Place+"%")
		argID++
	}

	if f.UserID != nil {
		query += " AND r.user_id = $" + fmt.Sprint(argID)
		args = append(args, *f.UserID)
		argID++
	}

	if f.UserName != nil {
		query += " AND u.full_name ILIKE $" + fmt.Sprint(argID)
		args = append(args, "%"+*f.UserName+"%")
		argID++
	}

	query += " ORDER BY r.created_at DESC"

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reportMap := make(map[string]*dtos.ReportDetailResponse)
	var result []dtos.ReportDetailResponse

	for rows.Next() {

		var (
			rID, userID, place, responsible string
			reportDate, createdAt           time.Time

			qID, qText *string
			answerText *string
			imageURL   *string
		)

		err := rows.Scan(
			&rID,
			&userID,
			&place,
			&reportDate,
			&responsible,
			&createdAt,
			&qID,
			&qText,
			&answerText,
			&imageURL,
		)
		if err != nil {
			return nil, err
		}

		// создаём report если нет
		if _, exists := reportMap[rID]; !exists {
			reportMap[rID] = &dtos.ReportDetailResponse{
				ID:              rID,
				UserID:          userID,
				Place:           place,
				ReportDate:      reportDate,
				ResponsibleName: responsible,
				CreatedAt:       createdAt,
				Answers:         []dtos.AnswerResponse{},
			}
		}

		// добавляем ответы
		if qID != nil {
			reportMap[rID].Answers = append(reportMap[rID].Answers, dtos.AnswerResponse{
				QuestionID:   *qID,
				QuestionText: *qText,
				AnswerText:   derefString(answerText),
				ImageURL:     imageURL,
			})
		}
	}

	for _, rep := range reportMap {
		result = append(result, *rep)
	}

	return result, nil
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
