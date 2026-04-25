package handler

import (
	"backend/internal/middleware"
	"backend/internal/models/dtos"
	"backend/internal/repository"
	"backend/internal/service/report_service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

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

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

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
		Limit:           10,
		Offset:          0,
		MetadataFilters: make(map[string]string),
	}

	if v := query.Get("date_from"); v != "" {
		filters.DateFrom = &v
	}

	if v := query.Get("date_to"); v != "" {
		filters.DateTo = &v
	}

	if v := query.Get("checklist_id"); v != "" {
		filters.ChecklistID = &v
	}

	if v := query.Get("user_id"); v != "" {
		filters.UserID = &v
	}

	if v := query.Get("user_name"); v != "" {
		filters.UserName = &v
	}

	if v := query.Get("limit"); v != "" {
		fmt.Sscanf(v, "%d", &filters.Limit)
	}

	if v := query.Get("offset"); v != "" {
		fmt.Sscanf(v, "%d", &filters.Offset)
	}

	for key, values := range query {
		if len(values) == 0 {
			continue
		}

		if strings.HasPrefix(key, "metadata_") {
			metaKey := strings.TrimPrefix(key, "metadata_")
			if metaKey != "" && values[0] != "" {
				filters.MetadataFilters[metaKey] = values[0]
			}
		}
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

	filters := repository.ReportFilters{
		MetadataFilters: make(map[string]string),
	}

	if v := query.Get("date_from"); v != "" {
		filters.DateFrom = &v
	}

	if v := query.Get("date_to"); v != "" {
		filters.DateTo = &v
	}

	if v := query.Get("checklist_id"); v != "" {
		filters.ChecklistID = &v
	}

	if v := query.Get("user_id"); v != "" {
		filters.UserID = &v
	}

	if v := query.Get("user_name"); v != "" {
		filters.UserName = &v
	}

	for key, values := range query {
		if len(values) == 0 {
			continue
		}

		if strings.HasPrefix(key, "metadata_") {
			metaKey := strings.TrimPrefix(key, "metadata_")
			if metaKey != "" && values[0] != "" {
				filters.MetadataFilters[metaKey] = values[0]
			}
		}
	}

	reports, err := h.reportService.ExportReports(r.Context(), filters)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка получения данных: %v", err), http.StatusInternalServerError)
		return
	}

	if len(reports) == 0 {
		http.Error(w, "Нет отчётов по заданным фильтрам", http.StatusNotFound)
		return
	}

	file := excelize.NewFile()
	sheetName := "Отчёты"

	if err := file.SetSheetName("Sheet1", sheetName); err != nil {
		http.Error(w, "Ошибка переименования листа", http.StatusInternalServerError)
		return
	}

	headerStyle, err := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 12,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#E0E0E0"},
			Pattern: 1,
		},
	})
	if err != nil {
		http.Error(w, "Ошибка создания стиля", http.StatusInternalServerError)
		return
	}

	row := 1
	reportNum := 1

	for _, report := range reports {
		if reportNum > 1 {
			row++
		}

		headerText := fmt.Sprintf(
			"Отчёт №%d   Дата: %s   Ответственный: %s",
			reportNum,
			report.ReportDate.Format("02.01.2006"),
			report.ResponsibleName,
		)

		startCell := fmt.Sprintf("A%d", row)
		endCell := fmt.Sprintf("D%d", row)

		if err := file.MergeCell(sheetName, startCell, endCell); err != nil {
			http.Error(w, "Ошибка объединения ячеек", http.StatusInternalServerError)
			return
		}

		file.SetCellValue(sheetName, startCell, headerText)
		file.SetCellStyle(sheetName, startCell, endCell, headerStyle)
		row++

		metaCell := fmt.Sprintf("A%d", row)
		metaEndCell := fmt.Sprintf("D%d", row)

		if err := file.MergeCell(sheetName, metaCell, metaEndCell); err != nil {
			http.Error(w, "Ошибка объединения ячеек", http.StatusInternalServerError)
			return
		}

		file.SetCellValue(sheetName, metaCell, fmt.Sprintf("Метаданные: %s", string(report.Metadata)))
		row++

		headers := []string{"Вопрос", "Ответ", "Результат", "Изображение"}
		for i, title := range headers {
			col := string(rune('A' + i))
			cell := fmt.Sprintf("%s%d", col, row)

			file.SetCellValue(sheetName, cell, title)
			file.SetCellStyle(sheetName, cell, cell, headerStyle)
		}
		row++

		if len(report.Answers) == 0 {
			emptyCell := fmt.Sprintf("A%d", row)
			emptyEnd := fmt.Sprintf("D%d", row)

			file.MergeCell(sheetName, emptyCell, emptyEnd)
			file.SetCellValue(sheetName, emptyCell, "Нет ответов")
			row++
		} else {
			for _, ans := range report.Answers {
				file.SetCellValue(sheetName, fmt.Sprintf("A%d", row), ans.QuestionText)
				file.SetCellValue(sheetName, fmt.Sprintf("B%d", row), ans.AnswerText)

				if ans.Result != nil {
					file.SetCellValue(sheetName, fmt.Sprintf("C%d", row), *ans.Result)
				} else {
					file.SetCellValue(sheetName, fmt.Sprintf("C%d", row), "")
				}

				if ans.ImageURL != nil {
					file.SetCellValue(sheetName, fmt.Sprintf("D%d", row), *ans.ImageURL)
				} else {
					file.SetCellValue(sheetName, fmt.Sprintf("D%d", row), "")
				}

				row++
			}
		}

		reportNum++
	}

	for _, col := range []string{"A", "B", "C", "D"} {
		file.SetColWidth(sheetName, col, col, 28)
	}

	fileName := fmt.Sprintf("reports_export_%s.xlsx", time.Now().Format("2006-01-02"))

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Cache-Control", "no-cache")

	if err := file.Write(w); err != nil {
		log.Printf("Ошибка записи Excel: %v", err)
	}
}
