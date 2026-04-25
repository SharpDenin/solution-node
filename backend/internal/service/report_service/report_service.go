package report_service

import (
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
	db            *repository.DB
	reportRepo    repository.ReportRepository
	questionRepo  repository.QuestionRepository
	checklistRepo repository.ChecklistRepository
	userRepo      repository.UserRepository
}

func NewReportService(
	db *repository.DB,
	reportRepo repository.ReportRepository,
	questionRepo repository.QuestionRepository,
	checklistRepo repository.ChecklistRepository,
	userRepo repository.UserRepository,
) *ReportService {
	return &ReportService{
		db:            db,
		reportRepo:    reportRepo,
		questionRepo:  questionRepo,
		checklistRepo: checklistRepo,
		userRepo:      userRepo,
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

		if question.ChecklistID != checklistID {
			return errors.New("question does not belong to checklist")
		}

		var imageURL *string
		if ans.ImageURL != "" {
			v := ans.ImageURL
			imageURL = &v
		}

		result := evaluateAnswer(question.Formula, ans.AnswerText)

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

func evaluateAnswer(formula *string, answerText string) *string {
	if formula == nil || strings.TrimSpace(*formula) == "" {
		neutral := "neutral"
		return &neutral
	}

	answerValue, err := strconv.ParseFloat(strings.ReplaceAll(answerText, ",", "."), 64)
	if err != nil {
		neutral := "neutral"
		return &neutral
	}

	f := strings.TrimSpace(*formula)

	operators := []string{">=", "<=", ">", "<", "="}

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
			case "=":
				ok = answerValue == threshold
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
