package report_service

import (
	"backend/internal/models/dtos"
	"backend/internal/repository"
	"context"
	"errors"

	"github.com/google/uuid"
)

type ReportService struct {
	db         *repository.DB
	reportRepo repository.ReportRepository
}

func NewReportService(db *repository.DB, reportRepo repository.ReportRepository) *ReportService {
	return &ReportService{
		db:         db,
		reportRepo: reportRepo,
	}
}

func (s *ReportService) CreateReport(ctx context.Context, userID string, req dtos.CreateReportRequest) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user id")
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
		req.Place,
		req.ReportDate,
		req.ResponsibleName,
	)
	if err != nil {
		return err
	}

	for _, ans := range req.Answers {
		qID, err := uuid.Parse(ans.QuestionID)
		if err != nil {
			return err
		}

		var imageURL *string
		if ans.ImageURL != "" {
			v := ans.ImageURL
			imageURL = &v
		}

		err = s.reportRepo.CreateAnswer(ctx, tx, reportID, qID, ans.AnswerText, imageURL)
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
