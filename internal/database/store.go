package database

import (
	"context"

	"github.com/cmellojr/modo-locadora/internal/models"
	"github.com/google/uuid"
)

// GameAvailability holds a game and its rental status for shelf display.
type GameAvailability struct {
	Game       models.Game
	Available  bool
	RenterName string // Non-empty when rented.
}

// ActiveRental holds rental info joined with game and member data for the admin returns page.
type ActiveRental struct {
	RentalID   uuid.UUID
	GameTitle  string
	CoverURL   string
	MemberName string
	RentedAt   string // Formatted date.
}

// Store defines the set of operations for the database layer.
type Store interface {
	// CreateMember persists a new member in the database.
	CreateMember(ctx context.Context, member *models.Member) error

	// GetMemberByID retrieves a member by their UUID.
	GetMemberByID(ctx context.Context, id uuid.UUID) (*models.Member, error)

	// GetMemberByProfileName retrieves a member by their profile name.
	GetMemberByProfileName(ctx context.Context, name string) (*models.Member, error)

	// NextMembershipNumber generates the next sequential membership number (1991-XXX).
	NextMembershipNumber(ctx context.Context) (string, error)

	// GetGameByID retrieves a game by its ID.
	GetGameByID(ctx context.Context, id uuid.UUID) (*models.Game, error)

	// AddGame persists a new game and creates one physical copy for it.
	AddGame(ctx context.Context, game *models.Game) error

	// UpdateGame updates the editable fields of an existing game.
	UpdateGame(ctx context.Context, game *models.Game) error

	// ListGames retrieves all games from the database.
	ListGames(ctx context.Context) ([]models.Game, error)

	// ListGamesWithAvailability returns all games with their rental status.
	ListGamesWithAvailability(ctx context.Context) ([]GameAvailability, error)

	// RentGame creates a rental for the given game to the given member.
	RentGame(ctx context.Context, gameID, memberID uuid.UUID) error

	// ReturnGame marks an active rental as returned.
	ReturnGame(ctx context.Context, rentalID uuid.UUID) error

	// ListActiveRentals returns all currently active (unreturned) rentals.
	ListActiveRentals(ctx context.Context) ([]ActiveRental, error)

	// RegisterRental records a new rental transaction.
	RegisterRental(ctx context.Context, rental *models.Rental) error
}
