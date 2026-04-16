package requests

type CreateQuestionRequest struct {
	Text       string `json:"text"`
	OrderIndex int    `json:"order_index"`
}

type UpdateQuestionRequest struct {
	ID         string `json:"id"`
	Text       string `json:"text"`
	OrderIndex int    `json:"order_index"`
	IsActive   bool   `json:"is_active"`
}
