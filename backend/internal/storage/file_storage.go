package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type FileStorage struct {
	UploadDir string
	BaseURL   string
}

func NewFileStorage(uploadDir, baseURL string) *FileStorage {
	_ = os.MkdirAll(uploadDir, os.ModePerm)

	return &FileStorage{
		UploadDir: uploadDir,
		BaseURL:   baseURL,
	}
}

func (s *FileStorage) SaveFile(file io.Reader, originalName string) (string, error) {

	ext := filepath.Ext(originalName)
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	path := filepath.Join(s.UploadDir, fileName)

	dst, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	return s.BaseURL + "/uploads/" + fileName, nil
}
