package responses

import "github.com/google/uuid"

type QuestionResponse struct {
	ID            uuid.UUID                           `json:"id"`
	Text          string                              `json:"text"`
	OrderIndex    int                                 `json:"order_index"`
	IsActive      bool                                `json:"is_active"`
	ChecklistID   uuid.UUID                           `json:"checklist_id"`
	Formula       *string                             `json:"formula"`
	ImageURL      *string                             `json:"image_url"`
	TechnicalCode *string                             `json:"technical_code"`
	DefaultAnswer *string                             `json:"default_answer"`
	Formulas      []QuestionPhenophaseFormulaResponse `json:"formulas,omitempty"`
}

type QuestionPhenophaseFormulaResponse struct {
	ID           uuid.UUID `json:"id"`
	QuestionID   uuid.UUID `json:"question_id"`
	PhenophaseID uuid.UUID `json:"phenophase_id"`
	Formula      string    `json:"formula"`
}
