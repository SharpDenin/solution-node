package handler

import (
	"backend/internal/handler/dtos/responses"
	"backend/internal/storage"
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type UploadHandler struct {
	storage *storage.FileStorage
}

func NewUploadHandler(storage *storage.FileStorage) *UploadHandler {
	return &UploadHandler{storage: storage}
}

func (h *UploadHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "file too large", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(handler.Filename))
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		http.Error(w, "unsupported file format", http.StatusBadRequest)
		return
	}

	fileName := uuid.New().String() + ext

	url, err := h.storage.SaveFile(file, fileName)
	if err != nil {
		http.Error(w, "save error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(responses.UploadResponse{
		URL: url,
	})
}
