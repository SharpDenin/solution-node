package handler

import (
	"backend/internal/middleware"
	"backend/internal/models/dtos"
	"backend/internal/repository"
	"backend/internal/service/report_service"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

func (h *ReportHandler) GetReports(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	filters := repository.ReportFilters{
		Limit:  10,
		Offset: 0,
	}

	if v := query.Get("date_from"); v != "" {
		filters.DateFrom = &v
	}

	if v := query.Get("date_to"); v != "" {
		filters.DateTo = &v
	}

	if v := query.Get("place"); v != "" {
		filters.Place = &v
	}

	if v := query.Get("user_id"); v != "" {
		filters.UserID = &v
	}

	if v := query.Get("limit"); v != "" {
		fmt.Sscanf(v, "%d", &filters.Limit)
	}

	if v := query.Get("offset"); v != "" {
		fmt.Sscanf(v, "%d", &filters.Offset)
	}

	reports, err := h.reportService.GetReports(r.Context(), filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(reports)
}

func (h *ReportHandler) GetReportByID(w http.ResponseWriter, r *http.Request) {

	id := strings.TrimPrefix(r.URL.Path, "/reports/")

	report, err := h.reportService.GetReportByID(r.Context(), id)
	if err != nil {
		http.Error(w, "report not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(report)
}
