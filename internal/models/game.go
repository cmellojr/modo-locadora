package models

import (
	"time"

	"github.com/google/uuid"
)

// Game represents a game metadata in the system.
type Game struct {
	ID             uuid.UUID
	Title          string
	IgdbID         string
	Platform       string
	Summary        string
	CoverURL       string
	SourceMagazine string
	CoverDisplay   string // CSS object-fit value: "cover", "contain", "fill"
	AcquiredAt     time.Time
}
