package dtos

import "time"

type CreateReportRequest struct {
	Place           string            `json:"place"`
	ReportDate      string            `json:"report_date"`
	ResponsibleName string            `json:"responsible_name"`
	Answers         []CreateAnswerDTO `json:"answers"`
}

type CreateAnswerDTO struct {
	QuestionID string `json:"question_id"`
	AnswerText string `json:"answer_text"`
	ImageURL   string `json:"image_url"`
}

type ReportResponse struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	Place           string    `json:"place"`
	ReportDate      time.Time `json:"report_date"`
	ResponsibleName string    `json:"responsible_name"`
	CreatedAt       time.Time `json:"created_at"`
}
