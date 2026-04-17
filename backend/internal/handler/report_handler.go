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
	// НОВЫЙ ПАРАМЕТР
	if v := query.Get("user_name"); v != "" {
		filters.UserName = &v
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
	if v := query.Get("user_name"); v != "" {
		filters.UserName = &v
	}

	// Получаем данные
	reports, err := h.reportService.ExportReports(r.Context(), filters)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка получения данных: %v", err), 500)
		return
	}

	if len(reports) == 0 {
		http.Error(w, "Нет отчётов по заданным фильтрам", 404)
		return
	}

	// Создаём Excel-файл
	file := excelize.NewFile()
	sheetName := "Отчёты"
	if err := file.SetSheetName("Sheet1", sheetName); err != nil {
		http.Error(w, "Ошибка переименования листа", 500)
		return
	}

	// Стиль для заголовка отчёта
	headerStyle, err := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 12},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#E0E0E0"}, Pattern: 1},
	})
	if err != nil {
		http.Error(w, "Ошибка создания стиля", 500)
		return
	}

	row := 1 // начинаем с первой строки
	reportNum := 1

	for _, report := range reports {
		if reportNum > 1 {
			row++ // пустая строка между отчётами
		}

		// ----- Заголовок отчёта (объединённый) -----
		headerText := fmt.Sprintf("Отчёт №%d   Дата: %s   Ответственный: %s",
			reportNum,
			report.ReportDate.Format("02.01.2006"),
			report.ResponsibleName,
		)
		if report.Place != "" {
			headerText += fmt.Sprintf("   Место: %s", report.Place)
		}

		startCell := fmt.Sprintf("A%d", row)
		endCell := fmt.Sprintf("C%d", row)
		if err := file.MergeCell(sheetName, startCell, endCell); err != nil {
			http.Error(w, "Ошибка объединения ячеек", 500)
			return
		}
		file.SetCellValue(sheetName, startCell, headerText)
		file.SetCellStyle(sheetName, startCell, endCell, headerStyle)
		row++

		// ----- Шапка колонок для текущего отчёта -----
		headers := []string{"Вопрос", "Ответ", "Изображение"}
		for i, h := range headers {
			col := string(rune('A' + i))
			file.SetCellValue(sheetName, fmt.Sprintf("%s%d", col, row), h)
		}
		row++

		// ----- Ответы -----
		if len(report.Answers) == 0 {
			emptyCell := fmt.Sprintf("A%d", row)
			emptyEnd := fmt.Sprintf("C%d", row)
			file.MergeCell(sheetName, emptyCell, emptyEnd)
			file.SetCellValue(sheetName, emptyCell, "Нет ответов")
			row++
		} else {
			for _, ans := range report.Answers {
				file.SetCellValue(sheetName, fmt.Sprintf("A%d", row), ans.QuestionText)
				file.SetCellValue(sheetName, fmt.Sprintf("B%d", row), ans.AnswerText)
				if ans.ImageURL != nil {
					file.SetCellValue(sheetName, fmt.Sprintf("C%d", row), *ans.ImageURL)
				} else {
					file.SetCellValue(sheetName, fmt.Sprintf("C%d", row), "")
				}
				row++
			}
		}
		reportNum++
	}

	// Автоширина колонок
	for _, col := range []string{"A", "B", "C"} {
		file.SetColWidth(sheetName, col, col, 25)
	}

	// ✅ Устанавливаем заголовки ТОЛЬКО после успешного создания файла
	fileName := fmt.Sprintf("export_%s.xlsx", time.Now().Format("2006-01-02"))
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Cache-Control", "no-cache")

	// Записываем файл в ответ
	if err := file.Write(w); err != nil {
		// Заголовки уже отправлены, ошибку можно только залогировать
		log.Printf("Ошибка записи Excel: %v", err)
	}
}
