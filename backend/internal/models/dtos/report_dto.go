package dtos

import (
	"encoding/json"
	"time"
)

type CreateReportRequest struct {
	ChecklistID     string                 `json:"checklist_id"`
	VarietyID       *string                `json:"variety_id"`
	PhenophaseID    *string                `json:"phenophase_id"`
	ReportDate      string                 `json:"report_date"`
	ResponsibleName string                 `json:"responsible_name"`
	Metadata        map[string]interface{} `json:"metadata"`
	Answers         []AnswerRequest        `json:"answers"`
}

type AnswerRequest struct {
	QuestionID string `json:"question_id"`
	AnswerText string `json:"answer_text"`
	ImageURL   string `json:"image_url"`
}

type ReportResponse struct {
	ID              string          `json:"id"`
	UserID          string          `json:"user_id"`
	ChecklistID     string          `json:"checklist_id"`
	VarietyID       *string         `json:"variety_id"`
	PhenophaseID    *string         `json:"phenophase_id"`
	Metadata        json.RawMessage `json:"metadata"`
	Place           string          `json:"place"`
	ReportDate      time.Time       `json:"report_date"`
	ResponsibleName string          `json:"responsible_name"`
	CreatedAt       time.Time       `json:"created_at"`
}
