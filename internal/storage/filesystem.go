package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Storage interface {
	SaveData(dataType string, data interface{}) error
}

type FileSystem struct {
	baseDir string
}

func NewFileSystem(baseDir string) (*FileSystem, error) {
	if err:= os.MkdirAll(baseDir, 0755); err != nil {
		return nil, err
	}
	return &FileSystem{baseDir: baseDir}, nil
}

func (fs *FileSystem) SaveData(dataType string, data interface{}) error {
	timestamp := time.Now().UTC()
	filename := fmt.Sprintf("%s_%s.json",
		dataType,
		timestamp.Format("2006-01-02T15-04-05"),
	)
	path := filepath.Join(fs.baseDir, 
		timestamp.Format("2006/01/02"),
		dataType,
	)

	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(path, filename)

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, jsonData, 0644)
}