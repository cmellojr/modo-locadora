package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// LocalStorage implements StorageProvider using the local file system.
// This is typically used for development or environments with persistent volumes.
type LocalStorage struct {
	basePath string
}

// NewLocalStorage creates a new LocalStorage provider pointing to the static assets directory.
func NewLocalStorage() *LocalStorage {
	return &LocalStorage{
		basePath: "web/static",
	}
}

// Save stores the file in the local file system under web/static/{folder}.
// It returns the relative URL path starting with /static/.
func (s *LocalStorage) Save(ctx context.Context, folder string, filename string, src io.Reader) (string, error) {
	// Ensure the directory exists within web/static
	dir := filepath.Join(s.basePath, folder)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("storage_err: failed to create directory %s: %w", dir, err)
	}

	savePath := filepath.Join(dir, filename)
	dst, err := os.Create(savePath)
	if err != nil {
		return "", fmt.Errorf("storage_err: failed to create file %s: %w", savePath, err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("storage_err: failed to write file %s: %w", savePath, err)
	}

	// Return the relative URL for the web interface
	return fmt.Sprintf("/static/%s/%s", folder, filename), nil
}

// Delete removes the file from the local file system.
func (s *LocalStorage) Delete(ctx context.Context, folder string, filename string) error {
	filePath := filepath.Join(s.basePath, folder, filename)
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("storage_err: failed to delete file %s: %w", filePath, err)
	}
	return nil
}
