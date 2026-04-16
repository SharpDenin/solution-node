package handler

import (
	"backend/internal/middleware"
	"backend/internal/models/dtos"
	"backend/internal/service/report_service"
	"encoding/json"
	"net/http"
)

type ReportHandler struct {
	reportService *report_service.ReportService
}

func NewReportHandler(reportService *report_service.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

func (h *ReportHandler) CreateReport(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreateReportRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)

	err := h.reportService.CreateReport(r.Context(), userID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
