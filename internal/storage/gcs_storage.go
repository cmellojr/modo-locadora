package storage

import (
	"context"
	"fmt"
	"io"
	"log"

	"cloud.google.com/go/storage"
)

// GCSStorage implements StorageProvider using Google Cloud Storage.
// This provider is intended for serverless production environments like Cloud Run.
type GCSStorage struct {
	client     *storage.Client
	bucketName string
}

// NewGCSStorage creates a new GCSStorage provider for the given bucket.
func NewGCSStorage(bucketName string) *GCSStorage {
	if bucketName == "" {
		log.Fatal("storage_err: STORAGE_BUCKET_NAME environment variable is required for production storage")
	}

	// Initialize the client once to be reused across requests.
	// Note: In a production app, we might want to pass a context here,
	// but for simplicity in this factory we use Background.
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("storage_err: failed to create GCS client: %v", err)
	}

	return &GCSStorage{
		client:     client,
		bucketName: bucketName,
	}
}

// Save uploads the file to Google Cloud Storage and returns the public URL.
func (s *GCSStorage) Save(ctx context.Context, folder string, filename string, src io.Reader) (string, error) {
	objectPath := fmt.Sprintf("%s/%s", folder, filename)
	wc := s.client.Bucket(s.bucketName).Object(objectPath).NewWriter(ctx)

	if _, err := io.Copy(wc, src); err != nil {
		return "", fmt.Errorf("storage_err: failed to upload to bucket %s: %w", s.bucketName, err)
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("storage_err: failed to finalize upload to %s: %w", objectPath, err)
	}

	// Return the public URL for the uploaded object.
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", s.bucketName, objectPath), nil
}

// Delete removes the object from Google Cloud Storage.
func (s *GCSStorage) Delete(ctx context.Context, folder string, filename string) error {
	objectPath := fmt.Sprintf("%s/%s", folder, filename)
	if err := s.client.Bucket(s.bucketName).Object(objectPath).Delete(ctx); err != nil {
		return fmt.Errorf("storage_err: failed to delete object %s from GCS: %w", objectPath, err)
	}
	return nil
}
