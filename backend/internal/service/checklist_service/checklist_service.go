package checklist_service

import (
	"backend/internal/handler/dtos/responses"
	"backend/internal/repository"
	"backend/internal/service/checklist_service/helpers"
	"context"

	"github.com/google/uuid"
)

type ChecklistService struct {
	checklistRepo repository.ChecklistRepository
	userRepo      repository.UserRepository
}

func NewChecklistService(
	checklistRepo repository.ChecklistRepository,
	userRepo repository.UserRepository,
) *ChecklistService {
	return &ChecklistService{
		checklistRepo: checklistRepo,
		userRepo:      userRepo,
	}
}

func (s *ChecklistService) GetAll(ctx context.Context) ([]responses.ChecklistResponse, error) {
	checklists, err := s.checklistRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]responses.ChecklistResponse, 0, len(checklists))

	for _, c := range checklists {
		res = append(res, responses.ChecklistResponse{
			ID:   c.ID.String(),
			Name: c.Name,
			Code: c.Code,
		})
	}

	return res, nil
}

func (s *ChecklistService) GetAvailableForUser(ctx context.Context, userID string) ([]responses.ChecklistResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	user, roleName, err := s.userRepo.GetByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	if roleName == "admin" {
		checklists, err := s.checklistRepo.GetAll(ctx)
		if err != nil {
			return nil, err
		}

		return helpers.MapChecklists(checklists), nil
	}

	checklists, err := s.checklistRepo.GetByRoleID(ctx, user.RoleID)
	if err != nil {
		return nil, err
	}

	return helpers.MapChecklists(checklists), nil
}
