package models

import "github.com/google/uuid"

// Game represents a game metadata in the system.
type Game struct {
	ID       uuid.UUID
	Title    string
	IgdbID   string
	Platform string
	Summary  string
}
