package requests

type CreateQuestionRequest struct {
	Text        string  `json:"text"`
	OrderIndex  int     `json:"order_index"`
	ChecklistID string  `json:"checklist_id"`
	Formula     *string `json:"formula"`
}

type UpdateQuestionRequest struct {
	Text        string  `json:"text"`
	OrderIndex  int     `json:"order_index"`
	IsActive    bool    `json:"is_active"`
	ChecklistID string  `json:"checklist_id"`
	Formula     *string `json:"formula"`
}
