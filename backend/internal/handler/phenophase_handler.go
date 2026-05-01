package handler

import (
	"backend/internal/handler/dtos/requests"
	"backend/internal/handler/dtos/responses"
	"backend/internal/models"
	"backend/internal/service/phenophase_service"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type PhenophaseHandler struct {
	service *phenophase_service.PhenophaseService
}

func NewPhenophaseHandler(service *phenophase_service.PhenophaseService) *PhenophaseHandler {
	return &PhenophaseHandler{service: service}
}

func (h *PhenophaseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req requests.CreatePhenophaseRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	err := h.service.Create(
		r.Context(),
		req.Name,
		req.Description,
		req.ImageURL,
		req.OrderIndex,
		req.MinCriticalTemperature,
		req.CriticalTemperature,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *PhenophaseHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	phenophases, err := h.service.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(mapPhenophasesToResponse(phenophases))
}

func (h *PhenophaseHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	phenophase, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(mapPhenophaseToResponse(*phenophase))
}

func (h *PhenophaseHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var req requests.UpdatePhenophaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	err := h.service.Update(
		r.Context(),
		id,
		req.Name,
		req.Description,
		req.ImageURL,
		req.OrderIndex,
		req.MinCriticalTemperature,
		req.CriticalTemperature,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *PhenophaseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	err := h.service.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func mapPhenophasesToResponse(phenophases []models.Phenophase) []responses.PhenophaseResponse {
	res := make([]responses.PhenophaseResponse, 0, len(phenophases))

	for _, phenophase := range phenophases {
		res = append(res, mapPhenophaseToResponse(phenophase))
	}

	return res
}

func mapPhenophaseToResponse(phenophase models.Phenophase) responses.PhenophaseResponse {
	return responses.PhenophaseResponse{
		ID:                     phenophase.ID.String(),
		Name:                   phenophase.Name,
		Description:            phenophase.Description,
		ImageURL:               phenophase.ImageURL,
		OrderIndex:             phenophase.OrderIndex,
		MinCriticalTemperature: phenophase.MinCriticalTemperature,
		CriticalTemperature:    phenophase.CriticalTemperature,
		CreatedAt:              phenophase.CreatedAt,
	}
}
