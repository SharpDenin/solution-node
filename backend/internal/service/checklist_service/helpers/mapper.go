package helpers

import (
	"backend/internal/handler/dtos/responses"
	"backend/internal/models"
)

func MapChecklists(checklists []models.Checklist) []responses.ChecklistResponse {
	res := make([]responses.ChecklistResponse, 0, len(checklists))

	for _, c := range checklists {
		res = append(res, responses.ChecklistResponse{
			ID:   c.ID.String(),
			Name: c.Name,
			Code: c.Code,
		})
	}

	return res
}
