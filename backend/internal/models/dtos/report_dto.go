package dtos

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
