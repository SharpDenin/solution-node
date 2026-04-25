package question_service

import (
	"backend/internal/handler/dtos/requests"
	"backend/internal/handler/dtos/responses"
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/service/question_service/helpers"
	"context"
	"errors"

	"github.com/google/uuid"
)

type QuestionService struct {
	repo repository.QuestionRepository
}

func NewQuestionService(repo repository.QuestionRepository) *QuestionService {
	return &QuestionService{repo: repo}
}

func (s *QuestionService) Create(ctx context.Context, req requests.CreateQuestionRequest) error {
	if req.Text == "" || req.ChecklistID == "" {
		return errors.New("invalid input")
	}

	checklistID, err := uuid.Parse(req.ChecklistID)
	if err != nil {
		return errors.New("invalid checklist id")
	}

	q := &models.Question{
		Text:        req.Text,
		OrderIndex:  req.OrderIndex,
		IsActive:    true,
		ChecklistID: checklistID,
		Formula:     req.Formula,
	}

	return s.repo.Create(ctx, q)
}

func (s *QuestionService) GetAll(ctx context.Context) ([]responses.QuestionResponse, error) {
	questions, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return helpers.MapQuestionsToResponse(questions), nil
}

func (s *QuestionService) GetByChecklist(ctx context.Context, checklistID string) ([]responses.QuestionResponse, error) {
	if checklistID == "" {
		return nil, errors.New("checklist id is required")
	}

	if _, err := uuid.Parse(checklistID); err != nil {
		return nil, errors.New("invalid checklist id")
	}

	questions, err := s.repo.GetByChecklist(ctx, checklistID)
	if err != nil {
		return nil, err
	}

	return helpers.MapQuestionsToResponse(questions), nil
}

func (s *QuestionService) Update(ctx context.Context, id uuid.UUID, req requests.UpdateQuestionRequest) error {
	if req.Text == "" || req.ChecklistID == "" {
		return errors.New("invalid input")
	}

	checklistID, err := uuid.Parse(req.ChecklistID)
	if err != nil {
		return errors.New("invalid checklist id")
	}

	q := &models.Question{
		ID:          id,
		Text:        req.Text,
		OrderIndex:  req.OrderIndex,
		IsActive:    req.IsActive,
		ChecklistID: checklistID,
		Formula:     req.Formula,
	}

	return s.repo.Update(ctx, q)
}

func (s *QuestionService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
