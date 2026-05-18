package storage

import (
	"context"
	"io"
	"log"
)

// StorageProvider defines the contract for persisting binary files
// such as game covers and club badges across different environments.
type StorageProvider interface {
	// Save stores the file and returns the public URL string or relative path, and an error.
	// bucketOrFolder is the sub-directory (local) or prefix (GCS) where the file will be stored.
	Save(ctx context.Context, bucketOrFolder string, filename string, src io.Reader) (string, error)
	// Delete removes the file from the storage backend.
	Delete(ctx context.Context, bucketOrFolder string, filename string) error
}

// NewStorageProvider is a factory function that returns the appropriate
// StorageProvider implementation based on the environment.
func NewStorageProvider(appEnv, bucketName string) StorageProvider {
	if appEnv == "production" {
		log.Println("System: Initializing Google Cloud Storage provider.")
		return NewGCSStorage(bucketName)
	}
	log.Println("System: Initializing Local File System storage provider.")
	return NewLocalStorage()
}
