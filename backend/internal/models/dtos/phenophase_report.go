package dtos

import "github.com/google/uuid"

type PhenophaseMatrixReportResponse struct {
	VarietyID uuid.UUID                `json:"variety_id"`
	Rows      []PhenophaseMatrixRow    `json:"rows"`
	Columns   []PhenophaseMatrixColumn `json:"columns"`
}

type PhenophaseMatrixColumn struct {
	PhenophaseID uuid.UUID `json:"phenophase_id"`
	Name         string    `json:"name"`
	OrderIndex   int       `json:"order_index"`
}

type PhenophaseMatrixRow struct {
	QuestionID uuid.UUID              `json:"question_id"`
	Text       string                 `json:"text"`
	OrderIndex int                    `json:"order_index"`
	Cells      []PhenophaseMatrixCell `json:"cells"`
}

type PhenophaseMatrixCell struct {
	PhenophaseID uuid.UUID  `json:"phenophase_id"`
	AnswerText   *string    `json:"answer_text"`
	Result       *string    `json:"result"`
	ImageURL     *string    `json:"image_url"`
	ReportID     *uuid.UUID `json:"report_id"`
}
