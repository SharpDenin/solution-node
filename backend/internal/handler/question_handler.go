package handler

import (
	"backend/internal/handler/dtos/requests"
	"backend/internal/service/question_service"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type QuestionHandler struct {
	questionService *question_service.QuestionService
}

func NewQuestionHandler(questionService *question_service.QuestionService) *QuestionHandler {
	return &QuestionHandler{
		questionService: questionService,
	}
}

func (h *QuestionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateQuestionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	question, err := h.questionService.Create(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusCreated, question)
}

func (h *QuestionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "invalid question id", http.StatusBadRequest)
		return
	}

	question, err := h.questionService.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, question)
}

func (h *QuestionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	questions, err := h.questionService.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, questions)
}

func (h *QuestionHandler) GetByChecklist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	checklistID, err := uuid.Parse(vars["id"]) // было "checklistId"
	if err != nil {
		http.Error(w, "invalid checklist id", http.StatusBadRequest)
		return
	}

	questions, err := h.questionService.GetByChecklist(r.Context(), checklistID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, questions)
}

func (h *QuestionHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "invalid question id", http.StatusBadRequest)
		return
	}

	var req requests.UpdateQuestionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	question, err := h.questionService.Update(r.Context(), id, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, question)
}

func (h *QuestionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "invalid question id", http.StatusBadRequest)
		return
	}

	if err := h.questionService.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
