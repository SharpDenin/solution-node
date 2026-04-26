package requests

import "github.com/google/uuid"

type CreateQuestionRequest struct {
	Text        string    `json:"text"`
	OrderIndex  int       `json:"order_index"`
	IsActive    *bool     `json:"is_active"`
	ChecklistID uuid.UUID `json:"checklist_id"`
	Formula     *string   `json:"formula"`
	ImageURL    *string   `json:"image_url"`

	Formulas []QuestionPhenophaseFormulaRequest `json:"formulas"`
}

type UpdateQuestionRequest struct {
	Text        string    `json:"text"`
	OrderIndex  int       `json:"order_index"`
	IsActive    bool      `json:"is_active"`
	ChecklistID uuid.UUID `json:"checklist_id"`
	Formula     *string   `json:"formula"`
	ImageURL    *string   `json:"image_url"`

	Formulas []QuestionPhenophaseFormulaRequest `json:"formulas"`
}

type QuestionPhenophaseFormulaRequest struct {
	PhenophaseID uuid.UUID `json:"phenophase_id"`
	Formula      string    `json:"formula"`
}
