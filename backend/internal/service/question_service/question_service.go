package question_service

import (
	"backend/internal/handler/dtos/requests"
	"backend/internal/handler/dtos/responses"
	"backend/internal/models"
	"backend/internal/repository"
	"context"

	"github.com/google/uuid"
)

type QuestionService struct {
	repo repository.QuestionRepository
}

func NewQuestionService(repo repository.QuestionRepository) *QuestionService {
	return &QuestionService{repo: repo}
}

func (s *QuestionService) Create(ctx context.Context, req requests.CreateQuestionRequest) error {

	q := &models.Question{
		Text:       req.Text,
		OrderIndex: req.OrderIndex,
		IsActive:   true,
	}

	return s.repo.Create(ctx, q)
}

func (s *QuestionService) GetAll(ctx context.Context) ([]responses.QuestionResponse, error) {

	questions, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]responses.QuestionResponse, 0, len(questions))

	for _, q := range questions {
		res = append(res, responses.QuestionResponse{
			ID:         q.ID.String(),
			Text:       q.Text,
			OrderIndex: q.OrderIndex,
			IsActive:   q.IsActive,
		})
	}

	return res, nil
}

func (s *QuestionService) Update(ctx context.Context, id uuid.UUID, req requests.UpdateQuestionRequest) error {
	q := &models.Question{
		ID:         id,
		Text:       req.Text,
		OrderIndex: req.OrderIndex,
		IsActive:   req.IsActive,
	}

	return s.repo.Update(ctx, q)
}

func (s *QuestionService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
