package phenophase_service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
)

type PhenophaseService struct {
	repo repository.PhenophaseRepository
}

func NewPhenophaseService(repo repository.PhenophaseRepository) *PhenophaseService {
	return &PhenophaseService{repo: repo}
}

func (s *PhenophaseService) Create(
	ctx context.Context,
	name string,
	description string,
	imageURL string,
	orderIndex int,
	minCriticalTemperature *float64,
	criticalTemperature *float64,
) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("name is required")
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

	phenophase := &models.Phenophase{
		Name:                   strings.TrimSpace(name),
		Description:            descPtr,
		ImageURL:               imagePtr,
		OrderIndex:             orderIndex,
		MinCriticalTemperature: minCriticalTemperature,
		CriticalTemperature:    criticalTemperature,
	}

	return s.repo.Create(ctx, phenophase)
}

func (s *PhenophaseService) GetAll(ctx context.Context) ([]models.Phenophase, error) {
	return s.repo.GetAll(ctx)
}

func (s *PhenophaseService) GetByID(ctx context.Context, id string) (*models.Phenophase, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid phenophase id")
	}

	return s.repo.GetByID(ctx, parsedID)
}

func (s *PhenophaseService) Update(
	ctx context.Context,
	id string,
	name string,
	description string,
	imageURL string,
	orderIndex int,
	minCriticalTemperature *float64,
	criticalTemperature *float64,
) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid phenophase id")
	}

	if strings.TrimSpace(name) == "" {
		return errors.New("name is required")
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

	phenophase := &models.Phenophase{
		ID:                     parsedID,
		Name:                   strings.TrimSpace(name),
		Description:            descPtr,
		ImageURL:               imagePtr,
		OrderIndex:             orderIndex,
		MinCriticalTemperature: minCriticalTemperature,
		CriticalTemperature:    criticalTemperature,
	}

	return s.repo.Update(ctx, phenophase)
}

func (s *PhenophaseService) Delete(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid phenophase id")
	}

	return s.repo.Delete(ctx, parsedID)
}
