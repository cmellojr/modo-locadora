package models

import "github.com/google/uuid"

// GameCopyStatus defines the availability status of a physical game copy.
type GameCopyStatus string

const (
	StatusAvailable GameCopyStatus = "available"
	StatusRented    GameCopyStatus = "rented"
)

// GameCopy represents a physical game copy (cartridge).
type GameCopy struct {
	ID     uuid.UUID
	GameID uuid.UUID
	Status GameCopyStatus
}
