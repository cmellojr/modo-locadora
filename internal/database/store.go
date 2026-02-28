package database

import (
	"context"

	"github.com/cmellojr/modo-locadora/internal/models"
	"github.com/google/uuid"
)

// Store defines the set of operations for the database layer.
type Store interface {
	// CreateMember persists a new member in the database.
	CreateMember(ctx context.Context, member *models.Member) error

	// GetGameByID retrieves a game by its ID.
	GetGameByID(ctx context.Context, id uuid.UUID) (*models.Game, error)

	// RegisterRental records a new rental transaction.
	RegisterRental(ctx context.Context, rental *models.Rental) error
}
