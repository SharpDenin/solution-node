package handler

import (
	"backend/internal/middleware"
	"backend/internal/models/dtos"
	"backend/internal/repository"
	"backend/internal/service/report_service"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xuri/excelize/v2"
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
	vars := mux.Vars(r)
	id := vars["id"]

	report, err := h.reportService.GetReportByID(r.Context(), id)
	if err != nil {
		http.Error(w, "report not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(report)
}

func (h *ReportHandler) ExportExcel(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	var filters repository.ReportFilters

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

	reports, err := h.reportService.ExportReports(r.Context(), filters)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	file := excelize.NewFile()
	sheet := "Reports"
	index, _ := file.NewSheet(sheet)
	file.SetActiveSheet(index)

	// HEADER
	file.SetCellValue(sheet, "A1", "Report ID")
	file.SetCellValue(sheet, "B1", "Place")
	file.SetCellValue(sheet, "C1", "Date")
	file.SetCellValue(sheet, "D1", "Responsible")
	file.SetCellValue(sheet, "E1", "Question")
	file.SetCellValue(sheet, "F1", "Answer")
	file.SetCellValue(sheet, "G1", "Image")

	row := 2

	for _, report := range reports {
		for _, ans := range report.Answers {

			file.SetCellValue(sheet, fmt.Sprintf("A%d", row), report.ID)
			file.SetCellValue(sheet, fmt.Sprintf("B%d", row), report.Place)
			file.SetCellValue(sheet, fmt.Sprintf("C%d", row), report.ReportDate.Format("2006-01-02"))
			file.SetCellValue(sheet, fmt.Sprintf("D%d", row), report.ResponsibleName)

			file.SetCellValue(sheet, fmt.Sprintf("E%d", row), ans.QuestionText)
			file.SetCellValue(sheet, fmt.Sprintf("F%d", row), ans.AnswerText)

			if ans.ImageURL != nil {
				file.SetCellValue(sheet, fmt.Sprintf("G%d", row), *ans.ImageURL)
			} else {
				file.SetCellValue(sheet, fmt.Sprintf("G%d", row), "")
			}

			row++
		}
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=reports.xlsx")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Cache-Control", "no-cache")

	if err := file.Write(w); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
