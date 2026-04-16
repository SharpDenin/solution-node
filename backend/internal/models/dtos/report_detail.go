package dtos

import "time"

type ReportDetailResponse struct {
	ID              string           `json:"id"`
	UserID          string           `json:"user_id"`
	Place           string           `json:"place"`
	ReportDate      time.Time        `json:"report_date"`
	ResponsibleName string           `json:"responsible_name"`
	CreatedAt       time.Time        `json:"created_at"`
	Answers         []AnswerResponse `json:"answers"`
}

type AnswerResponse struct {
	QuestionID   string  `json:"question_id"`
	QuestionText string  `json:"question_text"`
	AnswerText   string  `json:"answer_text"`
	ImageURL     *string `json:"image_url"`
}
