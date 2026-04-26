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
	DateFrom        *string
	DateTo          *string
	ChecklistID     *string
	VarietyID       *string
	PhenophaseID    *string
	UserID          *string
	UserName        *string
	MetadataFilters map[string]string
	Limit           int
	Offset          int
}

type ReportRepository interface {
	CreateReport(
		ctx context.Context,
		tx Tx,
		userID uuid.UUID,
		checklistID uuid.UUID,
		varietyID *uuid.UUID,
		phenophaseID *uuid.UUID,
		reportDate string,
		responsibleName string,
		metadata []byte,
	) (uuid.UUID, error)

	CreateAnswer(
		ctx context.Context,
		tx Tx,
		reportID uuid.UUID,
		questionID uuid.UUID,
		answerText string,
		imageURL *string,
		result *string,
	) error

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
	checklistID uuid.UUID,
	varietyID *uuid.UUID,
	phenophaseID *uuid.UUID,
	reportDate string,
	responsibleName string,
	metadata []byte,
) (uuid.UUID, error) {
	query := `
		INSERT INTO reports (
			user_id,
			checklist_id,
			variety_id,
			phenophase_id,
			report_date,
			responsible_name,
			metadata
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	var reportID uuid.UUID

	err := tx.QueryRow(ctx, query,
		userID,
		checklistID,
		varietyID,
		phenophaseID,
		reportDate,
		responsibleName,
		metadata,
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
	result *string,
) error {
	query := `
		INSERT INTO answers (
			report_id,
			question_id,
			answer_text,
			image_url,
			result
		)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := tx.Exec(ctx,
		query,
		reportID,
		questionID,
		answerText,
		imageURL,
		result,
	)

	return err
}

func (r *reportRepository) GetReports(ctx context.Context, f ReportFilters) ([]dtos.ReportResponse, error) {
	query := `
		SELECT 
			r.id,
			r.user_id,
			r.checklist_id,
			r.variety_id,
			r.phenophase_id,
			COALESCE(r.place, '') AS place,
			r.report_date,
			r.responsible_name,
			r.metadata,
			r.created_at
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

	if f.ChecklistID != nil {
		query += " AND r.checklist_id = $" + fmt.Sprint(argID)
		args = append(args, *f.ChecklistID)
		argID++
	}

	if f.VarietyID != nil {
		query += " AND r.variety_id = $" + fmt.Sprint(argID)
		args = append(args, *f.VarietyID)
		argID++
	}

	if f.PhenophaseID != nil {
		query += " AND r.phenophase_id = $" + fmt.Sprint(argID)
		args = append(args, *f.PhenophaseID)
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

	for key, value := range f.MetadataFilters {
		query += fmt.Sprintf(" AND r.metadata->>$%d ILIKE $%d", argID, argID+1)
		args = append(args, key, "%"+value+"%")
		argID += 2
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

	reports := make([]dtos.ReportResponse, 0)

	for rows.Next() {
		var report dtos.ReportResponse

		err := rows.Scan(
			&report.ID,
			&report.UserID,
			&report.ChecklistID,
			&report.VarietyID,
			&report.PhenophaseID,
			&report.Place,
			&report.ReportDate,
			&report.ResponsibleName,
			&report.Metadata,
			&report.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reports, nil
}

func (r *reportRepository) GetReportByID(ctx context.Context, id string) (*dtos.ReportDetailResponse, error) {
	query := `
		SELECT 
			r.id,
			r.user_id,
			r.checklist_id,
			r.variety_id,
			r.phenophase_id,
			COALESCE(r.place, '') AS place,
			r.report_date,
			r.responsible_name,
			r.metadata,
			r.created_at,
			q.id,
			q.text,
			a.answer_text,
			a.image_url,
			a.result
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
			reportID        string
			userID          string
			checklistID     string
			varietyID       *string
			phenophaseID    *string
			place           string
			reportDate      time.Time
			responsibleName string
			metadata        []byte
			createdAt       time.Time

			questionID   *string
			questionText *string
			answerText   *string
			imageURL     *string
			answerResult *string
		)

		err := rows.Scan(
			&reportID,
			&userID,
			&checklistID,
			&varietyID,
			&phenophaseID,
			&place,
			&reportDate,
			&responsibleName,
			&metadata,
			&createdAt,
			&questionID,
			&questionText,
			&answerText,
			&imageURL,
			&answerResult,
		)
		if err != nil {
			return nil, err
		}

		if report == nil {
			report = &dtos.ReportDetailResponse{
				ID:              reportID,
				UserID:          userID,
				ChecklistID:     checklistID,
				VarietyID:       varietyID,
				PhenophaseID:    phenophaseID,
				Place:           place,
				ReportDate:      reportDate,
				ResponsibleName: responsibleName,
				Metadata:        metadata,
				CreatedAt:       createdAt,
				Answers:         []dtos.AnswerResponse{},
			}
		}

		if questionID != nil {
			report.Answers = append(report.Answers, dtos.AnswerResponse{
				QuestionID:   *questionID,
				QuestionText: derefString(questionText),
				AnswerText:   derefString(answerText),
				ImageURL:     imageURL,
				Result:       answerResult,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
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
			r.checklist_id,
			r.variety_id,
			r.phenophase_id,
			COALESCE(r.place, '') AS place,
			r.report_date,
			r.responsible_name,
			r.metadata,
			r.created_at,
			q.id,
			q.text,
			a.answer_text,
			a.image_url,
			a.result
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

	if f.ChecklistID != nil {
		query += " AND r.checklist_id = $" + fmt.Sprint(argID)
		args = append(args, *f.ChecklistID)
		argID++
	}

	if f.VarietyID != nil {
		query += " AND r.variety_id = $" + fmt.Sprint(argID)
		args = append(args, *f.VarietyID)
		argID++
	}

	if f.PhenophaseID != nil {
		query += " AND r.phenophase_id = $" + fmt.Sprint(argID)
		args = append(args, *f.PhenophaseID)
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

	for key, value := range f.MetadataFilters {
		query += fmt.Sprintf(" AND r.metadata->>$%d ILIKE $%d", argID, argID+1)
		args = append(args, key, "%"+value+"%")
		argID += 2
	}

	query += " ORDER BY r.created_at DESC, q.order_index ASC"

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reportMap := make(map[string]*dtos.ReportDetailResponse)

	for rows.Next() {
		var (
			reportID        string
			userID          string
			checklistID     string
			varietyID       *string
			phenophaseID    *string
			place           string
			reportDate      time.Time
			responsibleName string
			metadata        []byte
			createdAt       time.Time

			questionID   *string
			questionText *string
			answerText   *string
			imageURL     *string
			answerResult *string
		)

		err := rows.Scan(
			&reportID,
			&userID,
			&checklistID,
			&varietyID,
			&phenophaseID,
			&place,
			&reportDate,
			&responsibleName,
			&metadata,
			&createdAt,
			&questionID,
			&questionText,
			&answerText,
			&imageURL,
			&answerResult,
		)
		if err != nil {
			return nil, err
		}

		report, exists := reportMap[reportID]
		if !exists {
			report = &dtos.ReportDetailResponse{
				ID:              reportID,
				UserID:          userID,
				ChecklistID:     checklistID,
				VarietyID:       varietyID,
				PhenophaseID:    phenophaseID,
				Place:           place,
				ReportDate:      reportDate,
				ResponsibleName: responsibleName,
				Metadata:        metadata,
				CreatedAt:       createdAt,
				Answers:         []dtos.AnswerResponse{},
			}

			reportMap[reportID] = report
		}

		if questionID != nil {
			report.Answers = append(report.Answers, dtos.AnswerResponse{
				QuestionID:   *questionID,
				QuestionText: derefString(questionText),
				AnswerText:   derefString(answerText),
				ImageURL:     imageURL,
				Result:       answerResult,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := make([]dtos.ReportDetailResponse, 0, len(reportMap))
	for _, report := range reportMap {
		result = append(result, *report)
	}

	return result, nil
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
