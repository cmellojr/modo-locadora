package models

import (
	"time"

	"github.com/google/uuid"
)

// Member represents a user in the system.
type Member struct {
	ID              uuid.UUID
	ProfileName     string
	Email           string
	PasswordHash    string
	FavoriteConsole string
	JoinedAt        time.Time
}
