package handler

import (
	"backend/internal/middleware"
	"backend/internal/service/checklist_service"
	"encoding/json"
	"net/http"
)

type ChecklistHandler struct {
	service *checklist_service.ChecklistService
}

func NewChecklistHandler(service *checklist_service.ChecklistService) *ChecklistHandler {
	return &ChecklistHandler{service: service}
}

func (h *ChecklistHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	checklists, err := h.service.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(checklists)
}

func (h *ChecklistHandler) GetAvailableForCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	checklists, err := h.service.GetAvailableForUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(checklists)
}
