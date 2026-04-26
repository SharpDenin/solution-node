package handler

import (
	"backend/internal/handler/dtos/requests"
	"backend/internal/handler/dtos/responses"
	"backend/internal/models"
	"backend/internal/service/variety_service"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type VarietyHandler struct {
	service *variety_service.VarietyService
}

func NewVarietyHandler(service *variety_service.VarietyService) *VarietyHandler {
	return &VarietyHandler{service: service}
}

func (h *VarietyHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req requests.CreateVarietyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	err := h.service.Create(
		r.Context(),
		req.Name,
		req.Description,
		req.Priority,
		req.ImageURL,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *VarietyHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	varieties, err := h.service.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(mapVarietiesToResponse(varieties))
}

func (h *VarietyHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	variety, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(mapVarietyToResponse(*variety))
}

func (h *VarietyHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var req requests.UpdateVarietyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	err := h.service.Update(
		r.Context(),
		id,
		req.Name,
		req.Description,
		req.Priority,
		req.ImageURL,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *VarietyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := h.service.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func mapVarietiesToResponse(varieties []models.Variety) []responses.VarietyResponse {
	res := make([]responses.VarietyResponse, 0, len(varieties))

	for _, variety := range varieties {
		res = append(res, mapVarietyToResponse(variety))
	}

	return res
}

func mapVarietyToResponse(variety models.Variety) responses.VarietyResponse {
	return responses.VarietyResponse{
		ID:          variety.ID.String(),
		Name:        variety.Name,
		Description: variety.Description,
		Priority:    variety.Priority,
		ImageURL:    variety.ImageURL,
		CreatedAt:   variety.CreatedAt,
	}
}
