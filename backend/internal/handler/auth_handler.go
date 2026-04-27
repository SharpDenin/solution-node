package handler

import (
	"backend/internal/handler/dtos/requests"
	"backend/internal/middleware"
	"backend/internal/service/auth_service"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type AuthHandler struct {
	authService *auth_service.AuthService
}

func NewAuthHandler(authService *auth_service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req requests.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.authService.Register(
		r.Context(),
		req.FullName,
		req.Login,
		req.Password,
		req.Role,
		req.Position,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req requests.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(
		r.Context(),
		req.FullName,
		req.Login,
		req.Password,
	)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.authService.GetCurrentUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.authService.GetAllUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *AuthHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var req requests.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.authService.UpdateUser(r.Context(), id, req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := h.authService.DeleteUser(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) RestoreUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := h.authService.RestoreUser(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) BlockUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := h.authService.BlockUser(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) UnblockUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := h.authService.UnblockUser(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
