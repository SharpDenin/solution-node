package question_service

import (
	"backend/internal/handler/dtos/requests"
	"backend/internal/handler/dtos/responses"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type QuestionService struct {
	questionRepo   repository.QuestionRepository
	formulaRepo    repository.QuestionFormulaRepository
	phenophaseRepo repository.PhenophaseRepository
}

func NewQuestionService(
	questionRepo repository.QuestionRepository,
	formulaRepo repository.QuestionFormulaRepository,
	phenophaseRepo repository.PhenophaseRepository,
) *QuestionService {
	return &QuestionService{
		questionRepo:   questionRepo,
		formulaRepo:    formulaRepo,
		phenophaseRepo: phenophaseRepo,
	}
}

func (s *QuestionService) Create(ctx context.Context, req requests.CreateQuestionRequest) (*responses.QuestionResponse, error) {
	if strings.TrimSpace(req.Text) == "" {
		return nil, errors.New("question text is required")
	}

	if req.ChecklistID == uuid.Nil {
		return nil, errors.New("checklist_id is required")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	question := &models.Question{
		Text:        strings.TrimSpace(req.Text),
		OrderIndex:  req.OrderIndex,
		IsActive:    isActive,
		ChecklistID: req.ChecklistID,
		Formula:     normalizeOptionalString(req.Formula),
		ImageURL:    normalizeOptionalString(req.ImageURL),
	}

	if err := s.questionRepo.Create(ctx, question); err != nil {
		return nil, err
	}

	if err := s.syncPhenophaseFormulas(ctx, question.ID, req.Formulas); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, question.ID)
}

func (s *QuestionService) GetByID(ctx context.Context, id uuid.UUID) (*responses.QuestionResponse, error) {
	if id == uuid.Nil {
		return nil, errors.New("question id is required")
	}

	question, err := s.questionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if question == nil {
		return nil, errors.New("question not found")
	}

	formulas, err := s.formulaRepo.GetByQuestion(ctx, id)
	if err != nil {
		return nil, err
	}

	return mapQuestionToResponse(question, formulas), nil
}

func (s *QuestionService) GetAll(ctx context.Context) ([]responses.QuestionResponse, error) {
	questions, err := s.questionRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return s.attachFormulas(ctx, questions)
}

func (s *QuestionService) GetByChecklist(ctx context.Context, checklistID uuid.UUID) ([]responses.QuestionResponse, error) {
	if checklistID == uuid.Nil {
		return nil, errors.New("checklist_id is required")
	}

	questions, err := s.questionRepo.GetByChecklist(ctx, checklistID)
	if err != nil {
		return nil, err
	}

	return s.attachFormulas(ctx, questions)
}

func (s *QuestionService) Update(ctx context.Context, id uuid.UUID, req requests.UpdateQuestionRequest) (*responses.QuestionResponse, error) {
	if id == uuid.Nil {
		return nil, errors.New("question id is required")
	}

	if strings.TrimSpace(req.Text) == "" {
		return nil, errors.New("question text is required")
	}

	if req.ChecklistID == uuid.Nil {
		return nil, errors.New("checklist_id is required")
	}

	existing, err := s.questionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, errors.New("question not found")
	}

	existing.Text = strings.TrimSpace(req.Text)
	existing.OrderIndex = req.OrderIndex
	existing.IsActive = req.IsActive
	existing.ChecklistID = req.ChecklistID
	existing.Formula = normalizeOptionalString(req.Formula)
	existing.ImageURL = normalizeOptionalString(req.ImageURL)

	if err := s.questionRepo.Update(ctx, existing); err != nil {
		return nil, err
	}

	if err := s.syncPhenophaseFormulas(ctx, existing.ID, req.Formulas); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, existing.ID)
}

func (s *QuestionService) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("question id is required")
	}

	existing, err := s.questionRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existing == nil {
		return errors.New("question not found")
	}

	if existing.TechnicalCode != nil && strings.TrimSpace(*existing.TechnicalCode) != "" {
		return errors.New("system question cannot be deleted")
	}

	return s.questionRepo.Delete(ctx, id)
}

func (s *QuestionService) GetByChecklistWithDefaults(
	ctx context.Context,
	checklistID uuid.UUID,
	phenophaseID uuid.UUID,
) ([]responses.QuestionResponse, error) {
	if checklistID == uuid.Nil {
		return nil, errors.New("checklist_id is required")
	}

	if phenophaseID == uuid.Nil {
		return nil, errors.New("phenophase_id is required")
	}

	phenophase, err := s.phenophaseRepo.GetByID(ctx, phenophaseID)
	if err != nil {
		return nil, err
	}

	if phenophase == nil {
		return nil, errors.New("phenophase not found")
	}

	questions, err := s.questionRepo.GetByChecklist(ctx, checklistID)
	if err != nil {
		return nil, err
	}

	result := make([]responses.QuestionResponse, 0, len(questions))

	for _, question := range questions {
		formulas, err := s.formulaRepo.GetByQuestion(ctx, question.ID)
		if err != nil {
			return nil, err
		}

		response := mapQuestionToResponse(&question, formulas)
		response.DefaultAnswer = resolveDefaultAnswer(&question, phenophase)

		result = append(result, *response)
	}

	return result, nil
}

func (s *QuestionService) syncPhenophaseFormulas(
	ctx context.Context,
	questionID uuid.UUID,
	formulas []requests.QuestionPhenophaseFormulaRequest,
) error {
	if questionID == uuid.Nil {
		return errors.New("question id is required")
	}

	if err := s.formulaRepo.DeleteByQuestion(ctx, questionID); err != nil {
		return err
	}

	for _, item := range formulas {
		if item.PhenophaseID == uuid.Nil {
			return errors.New("phenophase_id is required")
		}

		if strings.TrimSpace(item.Formula) == "" {
			continue
		}

		formula := &models.QuestionPhenophaseFormula{
			QuestionID:   questionID,
			PhenophaseID: item.PhenophaseID,
			Formula:      strings.TrimSpace(item.Formula),
		}

		if err := s.formulaRepo.CreateOrUpdate(ctx, formula); err != nil {
			return err
		}
	}

	return nil
}

func (s *QuestionService) attachFormulas(
	ctx context.Context,
	questions []models.Question,
) ([]responses.QuestionResponse, error) {
	result := make([]responses.QuestionResponse, 0, len(questions))

	for _, question := range questions {
		formulas, err := s.formulaRepo.GetByQuestion(ctx, question.ID)
		if err != nil {
			return nil, err
		}

		response := mapQuestionToResponse(&question, formulas)
		result = append(result, *response)
	}

	return result, nil
}

func mapQuestionToResponse(
	question *models.Question,
	formulas []models.QuestionPhenophaseFormula,
) *responses.QuestionResponse {
	response := &responses.QuestionResponse{
		ID:            question.ID,
		Text:          question.Text,
		OrderIndex:    question.OrderIndex,
		IsActive:      question.IsActive,
		ChecklistID:   question.ChecklistID,
		Formula:       question.Formula,
		ImageURL:      question.ImageURL,
		TechnicalCode: question.TechnicalCode,
		DefaultAnswer: nil,
		Formulas:      make([]responses.QuestionPhenophaseFormulaResponse, 0, len(formulas)),
	}

	for _, formula := range formulas {
		response.Formulas = append(response.Formulas, responses.QuestionPhenophaseFormulaResponse{
			ID:           formula.ID,
			QuestionID:   formula.QuestionID,
			PhenophaseID: formula.PhenophaseID,
			Formula:      formula.Formula,
		})
	}

	return response
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}

func resolveDefaultAnswer(
	question *models.Question,
	phenophase *models.Phenophase,
) *string {
	if question == nil || question.TechnicalCode == nil || phenophase == nil {
		return nil
	}

	switch strings.TrimSpace(*question.TechnicalCode) {
	case "min_critical_temperature":
		return formatFloatPtr(phenophase.MinCriticalTemperature)

	case "critical_temperature":
		return formatFloatPtr(phenophase.CriticalTemperature)

	default:
		return nil
	}
}

func formatFloatPtr(value *float64) *string {
	if value == nil {
		return nil
	}

	formatted := strings.TrimRight(strings.TrimRight(
		strconv.FormatFloat(*value, 'f', 2, 64),
		"0",
	), ".")

	return &formatted
}
