package storage

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestLocalStorageSaveAndDelete(t *testing.T) {
	tempDir := t.TempDir()
	ls := &LocalStorage{
		basePath: tempDir,
	}

	ctx := context.Background()
	folder := "covers"
	filename := "test_cover.png"
	fileContent := []byte("fake image binary content")

	// 1. Test Save
	src := bytes.NewReader(fileContent)
	url, err := ls.Save(ctx, folder, filename, src)
	if err != nil {
		t.Fatalf("expected Save to succeed, got: %v", err)
	}

	expectedURL := "/static/covers/test_cover.png"
	if url != expectedURL {
		t.Fatalf("expected URL %s, got %s", expectedURL, url)
	}

	savedPath := filepath.Join(tempDir, folder, filename)
	savedBytes, err := os.ReadFile(savedPath)
	if err != nil {
		t.Fatalf("failed to read saved file at %s: %v", savedPath, err)
	}
	if !bytes.Equal(savedBytes, fileContent) {
		t.Fatalf("saved file content mismatch")
	}

	// 2. Test Delete
	if err := ls.Delete(ctx, folder, filename); err != nil {
		t.Fatalf("expected Delete to succeed, got: %v", err)
	}

	if _, err := os.Stat(savedPath); !os.IsNotExist(err) {
		t.Fatalf("file still exists after Delete")
	}

	// Delete non-existent file should be idempotent (no error)
	if err := ls.Delete(ctx, folder, "non_existent.png"); err != nil {
		t.Fatalf("expected Delete of non-existent file to succeed silently, got: %v", err)
	}
}
