package report_service

import (
	"backend/internal/models"
	"backend/internal/models/dtos"
	"backend/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

const allowAdminCreateAnyChecklist = true

type ReportService struct {
	db             *repository.DB
	reportRepo     repository.ReportRepository
	questionRepo   repository.QuestionRepository
	checklistRepo  repository.ChecklistRepository
	userRepo       repository.UserRepository
	varietyRepo    repository.VarietyRepository
	phenophaseRepo repository.PhenophaseRepository
	formulaRepo    repository.QuestionFormulaRepository
}

func NewReportService(
	db *repository.DB,
	reportRepo repository.ReportRepository,
	questionRepo repository.QuestionRepository,
	checklistRepo repository.ChecklistRepository,
	userRepo repository.UserRepository,
	varietyRepo repository.VarietyRepository,
	phenophaseRepo repository.PhenophaseRepository,
	formulaRepo repository.QuestionFormulaRepository,
) *ReportService {
	return &ReportService{
		db:             db,
		reportRepo:     reportRepo,
		questionRepo:   questionRepo,
		checklistRepo:  checklistRepo,
		userRepo:       userRepo,
		varietyRepo:    varietyRepo,
		phenophaseRepo: phenophaseRepo,
		formulaRepo:    formulaRepo,
	}
}

func (s *ReportService) CreateReport(ctx context.Context, userID string, req dtos.CreateReportRequest) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user id")
	}

	checklistID, err := uuid.Parse(req.ChecklistID)
	if err != nil {
		return errors.New("invalid checklist id")
	}

	user, roleName, err := s.userRepo.GetByID(ctx, uid)
	if err != nil {
		return errors.New("user not found")
	}

	checklist, err := s.checklistRepo.GetByID(ctx, checklistID)
	if err != nil {
		return errors.New("checklist not found")
	}

	if !canCreateReportForChecklist(user.RoleID, roleName, checklist.AllowedRoleID) {
		return errors.New("access denied for this checklist")
	}

	if req.ReportDate == "" || req.ResponsibleName == "" {
		return errors.New("invalid report data")
	}

	var varietyID *uuid.UUID
	if req.VarietyID != nil && strings.TrimSpace(*req.VarietyID) != "" {
		parsedVarietyID, err := uuid.Parse(*req.VarietyID)
		if err != nil {
			return errors.New("invalid variety id")
		}

		if _, err := s.varietyRepo.GetByID(ctx, parsedVarietyID); err != nil {
			return errors.New("variety not found")
		}

		varietyID = &parsedVarietyID
	}

	var selectedPhenophase *models.Phenophase

	var phenophaseID *uuid.UUID
	if req.PhenophaseID != nil && strings.TrimSpace(*req.PhenophaseID) != "" {
		parsedPhenophaseID, err := uuid.Parse(*req.PhenophaseID)
		if err != nil {
			return errors.New("invalid phenophase id")
		}

		phenophase, err := s.phenophaseRepo.GetByID(ctx, parsedPhenophaseID)
		if err != nil {
			return errors.New("phenophase not found")
		}

		selectedPhenophase = phenophase
		phenophaseID = &parsedPhenophaseID
	}

	if checklist.Code == "sort_control" {
		if varietyID == nil {
			return errors.New("variety is required for phenophase checklist")
		}

		if phenophaseID == nil {
			return errors.New("phenophase is required for phenophase checklist")
		}
	}

	metadataBytes, err := json.Marshal(req.Metadata)
	if err != nil {
		return errors.New("invalid metadata")
	}

	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	reportID, err := s.reportRepo.CreateReport(
		ctx,
		tx,
		uid,
		checklistID,
		varietyID,
		phenophaseID,
		req.ReportDate,
		req.ResponsibleName,
		metadataBytes,
	)
	if err != nil {
		return err
	}

	for _, ans := range req.Answers {
		qID, err := uuid.Parse(ans.QuestionID)
		if err != nil {
			return errors.New("invalid question id")
		}

		question, err := s.questionRepo.GetByID(ctx, qID)
		if err != nil {
			return errors.New("question not found")
		}

		if question == nil {
			return errors.New("question not found")
		}

		if question.ChecklistID != checklistID {
			return errors.New("question does not belong to checklist")
		}

		var imageURL *string
		if strings.TrimSpace(ans.ImageURL) != "" {
			v := strings.TrimSpace(ans.ImageURL)
			imageURL = &v
		}

		formula := question.Formula

		if phenophaseID != nil {
			phenophaseFormula, err := s.formulaRepo.GetFormula(ctx, qID, *phenophaseID)
			if err != nil {
				return errors.New("failed to get question formula")
			}

			if phenophaseFormula != nil {
				formula = phenophaseFormula
			}
		}

		result := evaluateReportAnswer(question, formula, ans.AnswerText, selectedPhenophase)

		err = s.reportRepo.CreateAnswer(
			ctx,
			tx,
			reportID,
			qID,
			ans.AnswerText,
			imageURL,
			result,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (s *ReportService) GetReports(ctx context.Context, filters repository.ReportFilters) ([]dtos.ReportResponse, error) {
	return s.reportRepo.GetReports(ctx, filters)
}

func (s *ReportService) GetReportByID(ctx context.Context, id string) (*dtos.ReportDetailResponse, error) {
	return s.reportRepo.GetReportByID(ctx, id)
}

func (s *ReportService) ExportReports(ctx context.Context, filters repository.ReportFilters) ([]dtos.ReportDetailResponse, error) {
	return s.reportRepo.GetReportsDetailed(ctx, filters)
}

func (s *ReportService) GetPhenophaseMatrixReport(
	ctx context.Context,
	varietyID uuid.UUID,
) (*dtos.PhenophaseMatrixReportResponse, error) {
	if varietyID == uuid.Nil {
		return nil, errors.New("variety_id is required")
	}

	return s.reportRepo.GetPhenophaseMatrixReport(ctx, varietyID)
}

func (s *ReportService) DeleteReport(ctx context.Context, id string) error {
	reportID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid report id")
	}

	return s.reportRepo.DeleteReport(ctx, reportID)
}

func evaluateAnswer(formula *string, answerText string) *string {
	if formula == nil || strings.TrimSpace(*formula) == "" {
		neutral := "neutral"
		return &neutral
	}

	f := strings.TrimSpace(*formula)
	answer := strings.TrimSpace(answerText)

	if strings.HasPrefix(f, "=") {
		expected := strings.TrimSpace(strings.TrimPrefix(f, "="))

		if strings.EqualFold(answer, expected) {
			good := "good"
			return &good
		}

		bad := "bad"
		return &bad
	}

	answerValue, err := strconv.ParseFloat(strings.ReplaceAll(answer, ",", "."), 64)
	if err != nil {
		neutral := "neutral"
		return &neutral
	}

	operators := []string{">=", "<=", ">", "<"}

	for _, op := range operators {
		if strings.HasPrefix(f, op) {
			rawValue := strings.TrimSpace(strings.TrimPrefix(f, op))

			threshold, err := strconv.ParseFloat(strings.ReplaceAll(rawValue, ",", "."), 64)
			if err != nil {
				neutral := "neutral"
				return &neutral
			}

			ok := false

			switch op {
			case ">":
				ok = answerValue > threshold
			case ">=":
				ok = answerValue >= threshold
			case "<":
				ok = answerValue < threshold
			case "<=":
				ok = answerValue <= threshold
			}

			if ok {
				good := "good"
				return &good
			}

			bad := "bad"
			return &bad
		}
	}

	neutral := "neutral"
	return &neutral
}

func canCreateReportForChecklist(userRoleID uuid.UUID, roleName string, checklistRoleID uuid.UUID) bool {
	if allowAdminCreateAnyChecklist && roleName == "admin" {
		return true
	}

	return userRoleID == checklistRoleID
}

func evaluateReportAnswer(
	question *models.Question,
	formula *string,
	answerText string,
	phenophase *models.Phenophase,
) *string {
	if question == nil {
		return evaluateAnswer(formula, answerText)
	}

	if question.TechnicalCode != nil {
		switch strings.TrimSpace(*question.TechnicalCode) {
		case "actual_temperature":
			return evaluateActualTemperature(answerText, phenophase)

		case "min_critical_temperature", "critical_temperature":
			neutral := "neutral"
			return &neutral
		}
	}

	return evaluateAnswer(formula, answerText)
}

//

func evaluateActualTemperature(answerText string, phenophase *models.Phenophase) *string {
	if phenophase == nil || phenophase.MinCriticalTemperature == nil {
		neutral := "neutral"
		return &neutral
	}

	actual, err := strconv.ParseFloat(
		strings.ReplaceAll(strings.TrimSpace(answerText), ",", "."),
		64,
	)
	if err != nil {
		neutral := "neutral"
		return &neutral
	}

	minCritical := *phenophase.MinCriticalTemperature

	if actual > minCritical {
		good := "good"
		return &good
	}

	bad := "bad"
	return &bad
}
