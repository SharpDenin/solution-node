package dtos

import (
	"encoding/json"
	"time"
)

type ReportDetailResponse struct {
	ID              string           `json:"id"`
	UserID          string           `json:"user_id"`
	Place           string           `json:"place"`
	ReportDate      time.Time        `json:"report_date"`
	ResponsibleName string           `json:"responsible_name"`
	ChecklistID     string           `json:"checklist_id"`
	Metadata        json.RawMessage  `json:"metadata"`
	CreatedAt       time.Time        `json:"created_at"`
	Answers         []AnswerResponse `json:"answers"`
}

type AnswerResponse struct {
	QuestionID   string  `json:"question_id"`
	QuestionText string  `json:"question_text"`
	AnswerText   string  `json:"answer_text"`
	Result       *string `json:"result"`
	ImageURL     *string `json:"image_url"`
}
