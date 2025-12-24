package internal

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"projct/model"
)

func Reader(filesInfo chan model.FileInfo, filename os.DirEntry, ctx context.Context) error {
	path := filepath.Join("input", filename.Name())
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Error opening - %s Error: %v", filename, err)
	}
	defer f.Close()
	t, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("Error reading - %s Error: %v", filename, err)
	}
	fileInfo := model.FileInfo{
		Info: f,
		Path: path,
		Data: t,
	}
	filesInfo <- fileInfo
	return nil
}
