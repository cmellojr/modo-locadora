package database

import (
	"context"
	"time"

	"github.com/cmellojr/modo-locadora/internal/models"
	"github.com/google/uuid"
)

// GameAvailability holds a game and its copy/rental status for shelf display.
type GameAvailability struct {
	Game            models.Game
	TotalCopies     int
	AvailableCopies int
	RenterName      string // Non-empty when all copies are rented.
}

// ActiveRental holds rental info joined with game and member data for the admin returns page.
type ActiveRental struct {
	RentalID   uuid.UUID
	GameTitle  string
	CoverURL   string
	MemberName string
	RentedAt   string // Formatted date.
}

// ShameEntry holds data for the "Painel da Vergonha" (Wall of Shame).
type ShameEntry struct {
	ProfileName string
	LateCount   int
}

// PlatformSummary holds summary data for a console/platform in the shelf.
type PlatformSummary struct {
	Platform  string
	GameCount int
}

// ActivityEntry holds data for the "Aconteceu na Locadora" feed.
type ActivityEntry struct {
	ID         uuid.UUID
	EventType  string // "penalty", "redemption", "new_game", "prestige"
	MemberName string
	GameTitle  string
	CreatedAt  time.Time
}

// MemberRental holds a member's active rental for the carteirinha self-return.
type MemberRental struct {
	RentalID  uuid.UUID
	GameTitle string
	CoverURL  string
	Platform  string
	RentedAt  string // Formatted date.
	DueAt     string // Formatted date.
	IsOverdue bool
}

// GameDetail holds detailed info for a single game page.
type GameDetail struct {
	Game            models.Game
	TotalCopies     int
	AvailableCopies int
	TotalRentals    int
	TopRenterName   string
	TopRenterCount  int
	CurrentRenter   string
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

	// ListGamesWithAvailability returns games with their rental status, optionally filtered by platform.
	ListGamesWithAvailability(ctx context.Context, platform string) ([]GameAvailability, error)

	// ListPlatforms returns a summary of each platform in the catalog.
	ListPlatforms(ctx context.Context) ([]PlatformSummary, error)

	// GetGameDetail returns detailed info for a single game including rental stats.
	GetGameDetail(ctx context.Context, gameID uuid.UUID) (*GameDetail, error)

	// RentGame creates a rental for the given game to the given member.
	RentGame(ctx context.Context, gameID, memberID uuid.UUID) error

	// ReturnGame marks an active rental as returned.
	ReturnGame(ctx context.Context, rentalID uuid.UUID) error

	// ListActiveRentals returns all currently active (unreturned) rentals.
	ListActiveRentals(ctx context.Context) ([]ActiveRental, error)

	// RegisterRental records a new rental transaction.
	RegisterRental(ctx context.Context, rental *models.Rental) error

	// UpdateMemberNotes saves the member's password notebook text.
	UpdateMemberNotes(ctx context.Context, memberID uuid.UUID, notes string) error

	// GetMemberRentalStats returns counts of active and overdue rentals for a member.
	GetMemberRentalStats(ctx context.Context, memberID uuid.UUID) (activeCount, overdueCount int, err error)

	// ProcessOverdueRentals auto-returns overdue rentals and penalizes members.
	ProcessOverdueRentals(ctx context.Context) (int, error)

	// GetTopShameEntries returns the top N members with the most late returns.
	GetTopShameEntries(ctx context.Context, limit int) ([]ShameEntry, error)

	// RedeemMember resets a member's status from 'em_debito' to 'active'.
	RedeemMember(ctx context.Context, memberID uuid.UUID) error

	// GetMemberStatus returns the member's current status.
	GetMemberStatus(ctx context.Context, memberID uuid.UUID) (string, error)

	// InsertActivity records an event in the activities feed.
	InsertActivity(ctx context.Context, eventType, memberName, gameTitle string) error

	// ListRecentActivities returns the N most recent activity events.
	ListRecentActivities(ctx context.Context, limit int) ([]ActivityEntry, error)

	// ListMemberActiveRentals returns all active (unreturned) rentals for a specific member.
	ListMemberActiveRentals(ctx context.Context, memberID uuid.UUID) ([]MemberRental, error)

	// CountOnTimeReturns counts how many on-time returns a member has made.
	CountOnTimeReturns(ctx context.Context, memberID uuid.UUID) (int, error)

	// ReturnGameByMember returns a game for a specific member (validates ownership).
	// verdict stores the member's play status ("zerei", "joguei_um_pouco", "desisti").
	ReturnGameByMember(ctx context.Context, rentalID, memberID uuid.UUID, verdict string) error

	// GetRentalGameTitle returns the game title for a rental (used for activity logging).
	GetRentalGameTitle(ctx context.Context, rentalID uuid.UUID) (string, error)

	// ListCompletedGameIDs returns game IDs that the member has completed ("zerei").
	ListCompletedGameIDs(ctx context.Context, memberID uuid.UUID) ([]uuid.UUID, error)
}
