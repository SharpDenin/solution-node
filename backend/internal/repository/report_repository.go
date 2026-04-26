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
	GetPhenophaseMatrixReport(ctx context.Context, varietyID uuid.UUID) (*dtos.PhenophaseMatrixReportResponse, error)
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

func (r *reportRepository) GetPhenophaseMatrixReport(
	ctx context.Context,
	varietyID uuid.UUID,
) (*dtos.PhenophaseMatrixReportResponse, error) {
	columnsQuery := `
		SELECT id, name, order_index
		FROM phenophases
		ORDER BY order_index ASC
	`

	columnRows, err := r.db.Pool.Query(ctx, columnsQuery)
	if err != nil {
		return nil, err
	}
	defer columnRows.Close()

	var columns []dtos.PhenophaseMatrixColumn

	for columnRows.Next() {
		var column dtos.PhenophaseMatrixColumn

		if err := columnRows.Scan(
			&column.PhenophaseID,
			&column.Name,
			&column.OrderIndex,
		); err != nil {
			return nil, err
		}

		columns = append(columns, column)
	}

	if err := columnRows.Err(); err != nil {
		return nil, err
	}

	questionsQuery := `
		SELECT
			q.id,
			q.text,
			q.order_index
		FROM questions q
		INNER JOIN checklists c ON c.id = q.checklist_id
		WHERE c.code = 'sort_control'
		  AND q.is_active = true
		ORDER BY q.order_index ASC
	`

	questionRows, err := r.db.Pool.Query(ctx, questionsQuery)
	if err != nil {
		return nil, err
	}
	defer questionRows.Close()

	var rows []dtos.PhenophaseMatrixRow

	for questionRows.Next() {
		var row dtos.PhenophaseMatrixRow

		if err := questionRows.Scan(
			&row.QuestionID,
			&row.Text,
			&row.OrderIndex,
		); err != nil {
			return nil, err
		}

		row.Cells = make([]dtos.PhenophaseMatrixCell, 0, len(columns))

		for _, column := range columns {
			row.Cells = append(row.Cells, dtos.PhenophaseMatrixCell{
				PhenophaseID: column.PhenophaseID,
				ReportID:     nil,
				AnswerText:   nil,
				Result:       nil,
				ImageURL:     nil,
			})
		}

		rows = append(rows, row)
	}

	if err := questionRows.Err(); err != nil {
		return nil, err
	}

	answersQuery := `
		SELECT DISTINCT ON (a.question_id, r.phenophase_id)
			a.question_id,
			r.phenophase_id,
			r.id AS report_id,
			a.answer_text,
			a.result,
			a.image_url
		FROM answers a
		INNER JOIN reports r ON r.id = a.report_id
		INNER JOIN questions q ON q.id = a.question_id
		INNER JOIN checklists c ON c.id = q.checklist_id
		WHERE r.variety_id = $1
		  AND r.phenophase_id IS NOT NULL
		  AND c.code = 'sort_control'
		ORDER BY a.question_id, r.phenophase_id, r.report_date DESC, r.created_at DESC
	`

	answerRows, err := r.db.Pool.Query(ctx, answersQuery, varietyID)
	if err != nil {
		return nil, err
	}
	defer answerRows.Close()

	type matrixAnswer struct {
		QuestionID   uuid.UUID
		PhenophaseID uuid.UUID
		ReportID     uuid.UUID
		AnswerText   string
		Result       *string
		ImageURL     *string
	}

	answersMap := make(map[uuid.UUID]map[uuid.UUID]matrixAnswer)

	for answerRows.Next() {
		var answer matrixAnswer

		if err := answerRows.Scan(
			&answer.QuestionID,
			&answer.PhenophaseID,
			&answer.ReportID,
			&answer.AnswerText,
			&answer.Result,
			&answer.ImageURL,
		); err != nil {
			return nil, err
		}

		if _, ok := answersMap[answer.QuestionID]; !ok {
			answersMap[answer.QuestionID] = make(map[uuid.UUID]matrixAnswer)
		}

		answersMap[answer.QuestionID][answer.PhenophaseID] = answer
	}

	if err := answerRows.Err(); err != nil {
		return nil, err
	}

	for rowIndex := range rows {
		questionAnswers, ok := answersMap[rows[rowIndex].QuestionID]
		if !ok {
			continue
		}

		for cellIndex := range rows[rowIndex].Cells {
			phenophaseID := rows[rowIndex].Cells[cellIndex].PhenophaseID

			answer, ok := questionAnswers[phenophaseID]
			if !ok {
				continue
			}

			reportID := answer.ReportID
			answerText := answer.AnswerText

			rows[rowIndex].Cells[cellIndex].ReportID = &reportID
			rows[rowIndex].Cells[cellIndex].AnswerText = &answerText
			rows[rowIndex].Cells[cellIndex].Result = answer.Result
			rows[rowIndex].Cells[cellIndex].ImageURL = answer.ImageURL
		}
	}

	return &dtos.PhenophaseMatrixReportResponse{
		VarietyID: varietyID,
		Columns:   columns,
		Rows:      rows,
	}, nil
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
