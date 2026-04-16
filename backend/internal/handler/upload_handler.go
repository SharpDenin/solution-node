package handler

import (
	"backend/internal/handler/dtos/responses"
	"backend/internal/storage"
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
)

type UploadHandler struct {
	storage *storage.FileStorage
}

func NewUploadHandler(storage *storage.FileStorage) *UploadHandler {
	return &UploadHandler{storage: storage}
}

func (h *UploadHandler) UploadImage(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
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
		http.Error(w, "unsupported file format, only PNG, JPG, JPEG allowed", http.StatusBadRequest)
		return
	}

	url, err := h.storage.SaveFile(file, handler.Filename)
	if err != nil {
		http.Error(w, "save error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(responses.UploadResponse{
		URL: url,
	})
}
