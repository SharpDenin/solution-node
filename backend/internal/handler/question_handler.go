package handler

import (
	"backend/internal/handler/dtos/requests"
	"backend/internal/service/question_service"
	"encoding/json"
	"net/http"
	"strings"
)

type QuestionHandler struct {
	service *question_service.QuestionService
}

func NewQuestionHandler(s *question_service.QuestionService) *QuestionHandler {
	return &QuestionHandler{service: s}
}

func (h *QuestionHandler) Create(w http.ResponseWriter, r *http.Request) {

	var req requests.CreateQuestionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", 400)
		return
	}

	err := h.service.Create(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *QuestionHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	res, err := h.service.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(res)
}

func (h *QuestionHandler) Update(w http.ResponseWriter, r *http.Request) {

	var req requests.UpdateQuestionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", 400)
		return
	}

	err := h.service.Update(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *QuestionHandler) Delete(w http.ResponseWriter, r *http.Request) {

	id := strings.TrimPrefix(r.URL.Path, "/delete-questions/")

	err := h.service.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
}
