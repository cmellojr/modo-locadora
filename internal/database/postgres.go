package database

import (
	"context"
	"fmt"
	"time"

	"github.com/cmellojr/modo-locadora/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresStore implements the Store interface using PostgreSQL.
type PostgresStore struct {
	pool *pgxpool.Pool
}

// NewPostgresStore creates a new PostgresStore and initializes the connection pool.
func NewPostgresStore(ctx context.Context, connString string) (*PostgresStore, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return &PostgresStore{pool: pool}, nil
}

// Close closes the database connection pool.
func (s *PostgresStore) Close() {
	s.pool.Close()
}

// memberColumns is the shared column list for member queries.
const memberColumns = `id, profile_name, email, password_hash, favorite_console,
	COALESCE(membership_number, ''), COALESCE(address, ''), COALESCE(phone, ''),
	COALESCE(password_notes, ''), joined_at`

func scanMember(row pgx.Row) (*models.Member, error) {
	var m models.Member
	err := row.Scan(&m.ID, &m.ProfileName, &m.Email, &m.PasswordHash,
		&m.FavoriteConsole, &m.MembershipNumber, &m.Address, &m.Phone,
		&m.PasswordNotes, &m.JoinedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

// CreateMember persists a new member in the database.
func (s *PostgresStore) CreateMember(ctx context.Context, m *models.Member) error {
	query := `
		INSERT INTO members (id, profile_name, email, password_hash, favorite_console, membership_number, address, phone, password_notes, joined_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := s.pool.Exec(ctx, query, m.ID, m.ProfileName, m.Email, m.PasswordHash,
		m.FavoriteConsole, m.MembershipNumber, m.Address, m.Phone, m.PasswordNotes, m.JoinedAt)
	if err != nil {
		return fmt.Errorf("failed to create member: %w", err)
	}
	return nil
}

// GetMemberByID retrieves a member by their UUID.
func (s *PostgresStore) GetMemberByID(ctx context.Context, id uuid.UUID) (*models.Member, error) {
	query := `SELECT ` + memberColumns + ` FROM members WHERE id = $1`
	m, err := scanMember(s.pool.QueryRow(ctx, query, id))
	if err != nil {
		return nil, fmt.Errorf("failed to get member by id: %w", err)
	}
	return m, nil
}

// GetMemberByProfileName retrieves a member by their profile name.
func (s *PostgresStore) GetMemberByProfileName(ctx context.Context, name string) (*models.Member, error) {
	query := `SELECT ` + memberColumns + ` FROM members WHERE profile_name = $1`
	m, err := scanMember(s.pool.QueryRow(ctx, query, name))
	if err != nil {
		return nil, fmt.Errorf("failed to get member by profile name: %w", err)
	}
	return m, nil
}

// NextMembershipNumber generates the next sequential membership number (1991-XXX).
func (s *PostgresStore) NextMembershipNumber(ctx context.Context) (string, error) {
	var seq int
	err := s.pool.QueryRow(ctx, `SELECT nextval('membership_seq')`).Scan(&seq)
	if err != nil {
		return "", fmt.Errorf("failed to get next membership number: %w", err)
	}
	return fmt.Sprintf("1991-%03d", seq), nil
}

// GetGameByID retrieves a game by its ID.
func (s *PostgresStore) GetGameByID(ctx context.Context, id uuid.UUID) (*models.Game, error) {
	query := `SELECT id, title, igdb_id, platform, summary, cover_url, source_magazine, acquired_at FROM games WHERE id = $1`

	var g models.Game
	err := s.pool.QueryRow(ctx, query, id).Scan(&g.ID, &g.Title, &g.IgdbID, &g.Platform, &g.Summary, &g.CoverURL, &g.SourceMagazine, &g.AcquiredAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get game: %w", err)
	}
	return &g, nil
}

// AddGame persists a new game and creates one physical copy for it.
func (s *PostgresStore) AddGame(ctx context.Context, g *models.Game) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	gameQuery := `
		INSERT INTO games (id, title, igdb_id, platform, summary, cover_url, source_magazine, acquired_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = tx.Exec(ctx, gameQuery, g.ID, g.Title, g.IgdbID, g.Platform, g.Summary, g.CoverURL, g.SourceMagazine, g.AcquiredAt)
	if err != nil {
		return fmt.Errorf("failed to add game: %w", err)
	}

	copyQuery := `INSERT INTO game_copies (id, game_id, status) VALUES ($1, $2, 'available')`
	_, err = tx.Exec(ctx, copyQuery, uuid.New(), g.ID)
	if err != nil {
		return fmt.Errorf("failed to create game copy: %w", err)
	}

	return tx.Commit(ctx)
}

// UpdateGame updates the editable fields of an existing game.
func (s *PostgresStore) UpdateGame(ctx context.Context, g *models.Game) error {
	query := `
		UPDATE games
		SET title = $2, platform = $3, summary = $4, cover_url = $5, source_magazine = $6
		WHERE id = $1`

	tag, err := s.pool.Exec(ctx, query, g.ID, g.Title, g.Platform, g.Summary, g.CoverURL, g.SourceMagazine)
	if err != nil {
		return fmt.Errorf("failed to update game: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("game not found: %s", g.ID)
	}
	return nil
}

// ListGames retrieves all games from the database.
func (s *PostgresStore) ListGames(ctx context.Context) ([]models.Game, error) {
	query := `SELECT id, title, igdb_id, platform, summary, cover_url, source_magazine, acquired_at FROM games ORDER BY acquired_at DESC`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query games: %w", err)
	}
	defer rows.Close()

	var games []models.Game
	for rows.Next() {
		var g models.Game
		if err := rows.Scan(&g.ID, &g.Title, &g.IgdbID, &g.Platform, &g.Summary, &g.CoverURL, &g.SourceMagazine, &g.AcquiredAt); err != nil {
			return nil, fmt.Errorf("failed to scan game: %w", err)
		}
		games = append(games, g)
	}
	return games, nil
}

// ListGamesWithAvailability returns all games with copy counts and rental status.
func (s *PostgresStore) ListGamesWithAvailability(ctx context.Context) ([]GameAvailability, error) {
	query := `
		SELECT g.id, g.title, g.igdb_id, g.platform, g.summary, g.cover_url, g.source_magazine, g.acquired_at,
			COUNT(gc.id) AS total_copies,
			COUNT(gc.id) FILTER (WHERE gc.status = 'available') AS available_copies,
			COALESCE(
				(SELECT m.profile_name FROM rentals r2
				 JOIN game_copies gc2 ON gc2.id = r2.copy_id
				 JOIN members m ON m.id = r2.member_id
				 WHERE gc2.game_id = g.id AND r2.returned_at IS NULL
				 LIMIT 1), '') AS renter_name
		FROM games g
		LEFT JOIN game_copies gc ON gc.game_id = g.id
		GROUP BY g.id
		ORDER BY g.acquired_at DESC`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query games with availability: %w", err)
	}
	defer rows.Close()

	var result []GameAvailability
	for rows.Next() {
		var ga GameAvailability
		if err := rows.Scan(
			&ga.Game.ID, &ga.Game.Title, &ga.Game.IgdbID, &ga.Game.Platform,
			&ga.Game.Summary, &ga.Game.CoverURL, &ga.Game.SourceMagazine, &ga.Game.AcquiredAt,
			&ga.TotalCopies, &ga.AvailableCopies, &ga.RenterName,
		); err != nil {
			return nil, fmt.Errorf("failed to scan game availability: %w", err)
		}
		result = append(result, ga)
	}
	return result, nil
}

// RentGame creates a rental for the given game to the given member.
func (s *PostgresStore) RentGame(ctx context.Context, gameID, memberID uuid.UUID) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Find an available copy for this game.
	var copyID uuid.UUID
	err = tx.QueryRow(ctx,
		`SELECT id FROM game_copies WHERE game_id = $1 AND status = 'available' LIMIT 1 FOR UPDATE`,
		gameID).Scan(&copyID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("no available copies for this game")
		}
		return fmt.Errorf("failed to find available copy: %w", err)
	}

	// Mark the copy as rented.
	_, err = tx.Exec(ctx, `UPDATE game_copies SET status = 'rented' WHERE id = $1`, copyID)
	if err != nil {
		return fmt.Errorf("failed to update copy status: %w", err)
	}

	// Create the rental record (3-day due date).
	rentalID := uuid.New()
	now := time.Now()
	dueAt := now.AddDate(0, 0, 3)
	_, err = tx.Exec(ctx,
		`INSERT INTO rentals (id, member_id, copy_id, rented_at, due_at) VALUES ($1, $2, $3, $4, $5)`,
		rentalID, memberID, copyID, now, dueAt)
	if err != nil {
		return fmt.Errorf("failed to create rental: %w", err)
	}

	return tx.Commit(ctx)
}

// ReturnGame marks an active rental as returned and makes the copy available.
func (s *PostgresStore) ReturnGame(ctx context.Context, rentalID uuid.UUID) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Get the copy_id from the rental.
	var copyID uuid.UUID
	err = tx.QueryRow(ctx,
		`SELECT copy_id FROM rentals WHERE id = $1 AND returned_at IS NULL`, rentalID).Scan(&copyID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("rental not found or already returned")
		}
		return fmt.Errorf("failed to find rental: %w", err)
	}

	// Mark the rental as returned.
	_, err = tx.Exec(ctx, `UPDATE rentals SET returned_at = NOW() WHERE id = $1`, rentalID)
	if err != nil {
		return fmt.Errorf("failed to update rental: %w", err)
	}

	// Mark the copy as available.
	_, err = tx.Exec(ctx, `UPDATE game_copies SET status = 'available' WHERE id = $1`, copyID)
	if err != nil {
		return fmt.Errorf("failed to update copy status: %w", err)
	}

	return tx.Commit(ctx)
}

// ListActiveRentals returns all currently active (unreturned) rentals.
func (s *PostgresStore) ListActiveRentals(ctx context.Context) ([]ActiveRental, error) {
	query := `
		SELECT r.id, g.title, g.cover_url, m.profile_name, r.rented_at
		FROM rentals r
		JOIN game_copies gc ON gc.id = r.copy_id
		JOIN games g ON g.id = gc.game_id
		JOIN members m ON m.id = r.member_id
		WHERE r.returned_at IS NULL
		ORDER BY r.rented_at DESC`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query active rentals: %w", err)
	}
	defer rows.Close()

	var result []ActiveRental
	for rows.Next() {
		var ar ActiveRental
		var rentedAt time.Time
		if err := rows.Scan(&ar.RentalID, &ar.GameTitle, &ar.CoverURL, &ar.MemberName, &rentedAt); err != nil {
			return nil, fmt.Errorf("failed to scan rental: %w", err)
		}
		ar.RentedAt = rentedAt.Format("02/01/2006")
		result = append(result, ar)
	}
	return result, nil
}

// RegisterRental records a new rental transaction.
func (s *PostgresStore) RegisterRental(ctx context.Context, r *models.Rental) error {
	query := `
		INSERT INTO rentals (id, member_id, copy_id, rented_at, due_at, returned_at, personal_note, public_legacy)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := s.pool.Exec(ctx, query, r.ID, r.MemberID, r.CopyID, r.RentedAt, r.DueAt, r.ReturnedAt, r.PersonalNote, r.PublicLegacy)
	if err != nil {
		return fmt.Errorf("failed to register rental: %w", err)
	}
	return nil
}

// UpdateMemberNotes saves the member's password notebook text.
func (s *PostgresStore) UpdateMemberNotes(ctx context.Context, memberID uuid.UUID, notes string) error {
	tag, err := s.pool.Exec(ctx, `UPDATE members SET password_notes = $1 WHERE id = $2`, notes, memberID)
	if err != nil {
		return fmt.Errorf("failed to update member notes: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("member not found: %s", memberID)
	}
	return nil
}

// GetMemberRentalStats returns counts of active and overdue rentals for a member.
func (s *PostgresStore) GetMemberRentalStats(ctx context.Context, memberID uuid.UUID) (activeCount, overdueCount int, err error) {
	query := `
		SELECT
			COUNT(*) FILTER (WHERE returned_at IS NULL) AS active,
			COUNT(*) FILTER (WHERE returned_at IS NULL AND due_at < NOW()) AS overdue
		FROM rentals WHERE member_id = $1`
	err = s.pool.QueryRow(ctx, query, memberID).Scan(&activeCount, &overdueCount)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get rental stats: %w", err)
	}
	return activeCount, overdueCount, nil
}
