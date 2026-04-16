package responses

type QuestionResponse struct {
	ID         string `json:"id"`
	Text       string `json:"text"`
	OrderIndex int    `json:"order_index"`
	IsActive   bool   `json:"is_active"`
}
