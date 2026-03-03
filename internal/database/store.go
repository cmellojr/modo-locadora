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

	// GetMemberByID retrieves a member by their UUID.
	GetMemberByID(ctx context.Context, id uuid.UUID) (*models.Member, error)

	// GetMemberByProfileName retrieves a member by their profile name.
	GetMemberByProfileName(ctx context.Context, name string) (*models.Member, error)

	// GetGameByID retrieves a game by its ID.
	GetGameByID(ctx context.Context, id uuid.UUID) (*models.Game, error)

	// AddGame persists a new game in the database.
	AddGame(ctx context.Context, game *models.Game) error

	// ListGames retrieves all games from the database.
	ListGames(ctx context.Context) ([]models.Game, error)

	// RegisterRental records a new rental transaction.
	RegisterRental(ctx context.Context, rental *models.Rental) error
}
