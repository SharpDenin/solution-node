package helpers

import (
	"backend/internal/handler/dtos/responses"
	"backend/internal/models"
)

func MapQuestionsToResponse(questions []models.Question) []responses.QuestionResponse {
	res := make([]responses.QuestionResponse, 0, len(questions))

	for _, q := range questions {
		res = append(res, responses.QuestionResponse{
			ID:          q.ID.String(),
			Text:        q.Text,
			OrderIndex:  q.OrderIndex,
			IsActive:    q.IsActive,
			ChecklistID: q.ChecklistID.String(),
			Formula:     q.Formula,
		})
	}

	return res
}
