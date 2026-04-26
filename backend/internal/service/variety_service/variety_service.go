package variety_service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
)

type VarietyService struct {
	repo repository.VarietyRepository
}

func NewVarietyService(repo repository.VarietyRepository) *VarietyService {
	return &VarietyService{repo: repo}
}

func (s *VarietyService) Create(ctx context.Context, name, description, priority, imageURL string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("name is required")
	}

	if priority == "" {
		priority = "medium"
	}

	var descPtr *string
	if strings.TrimSpace(description) != "" {
		v := strings.TrimSpace(description)
		descPtr = &v
	}

	var imagePtr *string
	if strings.TrimSpace(imageURL) != "" {
		v := strings.TrimSpace(imageURL)
		imagePtr = &v
	}

	variety := &models.Variety{
		Name:        strings.TrimSpace(name),
		Description: descPtr,
		Priority:    priority,
		ImageURL:    imagePtr,
	}

	return s.repo.Create(ctx, variety)
}

func (s *VarietyService) GetAll(ctx context.Context) ([]models.Variety, error) {
	return s.repo.GetAll(ctx)
}

func (s *VarietyService) GetByID(ctx context.Context, id string) (*models.Variety, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid variety id")
	}

	return s.repo.GetByID(ctx, parsedID)
}

func (s *VarietyService) Update(ctx context.Context, id, name, description, priority, imageURL string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid variety id")
	}

	if strings.TrimSpace(name) == "" {
		return errors.New("name is required")
	}

	if priority == "" {
		priority = "medium"
	}

	var descPtr *string
	if strings.TrimSpace(description) != "" {
		v := strings.TrimSpace(description)
		descPtr = &v
	}

	var imagePtr *string
	if strings.TrimSpace(imageURL) != "" {
		v := strings.TrimSpace(imageURL)
		imagePtr = &v
	}

	variety := &models.Variety{
		ID:          parsedID,
		Name:        strings.TrimSpace(name),
		Description: descPtr,
		Priority:    priority,
		ImageURL:    imagePtr,
	}

	return s.repo.Update(ctx, variety)
}

func (s *VarietyService) Delete(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid variety id")
	}

	return s.repo.Delete(ctx, parsedID)
}
