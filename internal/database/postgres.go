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

// ExecRaw executes a raw SQL string against the database (used for seed scripts).
func (s *PostgresStore) ExecRaw(ctx context.Context, sql string) error {
	_, err := s.pool.Exec(ctx, sql)
	return err
}

// Close closes the database connection pool.
func (s *PostgresStore) Close() {
	s.pool.Close()
}

// memberColumns is the shared column list for member queries.
const memberColumns = `id, profile_name, email, password_hash, favorite_console,
	COALESCE(membership_number, ''), COALESCE(address, ''), COALESCE(phone, ''),
	COALESCE(password_notes, ''), COALESCE(status, 'active'), COALESCE(late_count, 0),
	joined_at`

func scanMember(row pgx.Row) (*models.Member, error) {
	var m models.Member
	err := row.Scan(&m.ID, &m.ProfileName, &m.Email, &m.PasswordHash,
		&m.FavoriteConsole, &m.MembershipNumber, &m.Address, &m.Phone,
		&m.PasswordNotes, &m.Status, &m.LateCount, &m.JoinedAt)
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

// ListGamesWithAvailability returns games with copy counts and rental status, optionally filtered by platform.
func (s *PostgresStore) ListGamesWithAvailability(ctx context.Context, platform string) ([]GameAvailability, error) {
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
		LEFT JOIN game_copies gc ON gc.game_id = g.id`

	var args []interface{}
	if platform != "" {
		query += ` WHERE g.platform = $1`
		args = append(args, platform)
	}

	query += `
		GROUP BY g.id
		ORDER BY g.title ASC`

	rows, err := s.pool.Query(ctx, query, args...)
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

// ListPlatforms returns a summary of each platform in the catalog.
func (s *PostgresStore) ListPlatforms(ctx context.Context) ([]PlatformSummary, error) {
	query := `
		SELECT platform, COUNT(*) AS game_count
		FROM games
		GROUP BY platform
		ORDER BY platform ASC`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query platforms: %w", err)
	}
	defer rows.Close()

	var result []PlatformSummary
	for rows.Next() {
		var ps PlatformSummary
		if err := rows.Scan(&ps.Platform, &ps.GameCount); err != nil {
			return nil, fmt.Errorf("failed to scan platform summary: %w", err)
		}
		result = append(result, ps)
	}
	return result, nil
}

// GetGameDetail returns detailed info for a single game including rental stats.
func (s *PostgresStore) GetGameDetail(ctx context.Context, gameID uuid.UUID) (*GameDetail, error) {
	// Base game + copy counts.
	query := `
		SELECT g.id, g.title, g.igdb_id, g.platform, g.summary, g.cover_url, g.source_magazine, g.acquired_at,
			COUNT(gc.id) AS total_copies,
			COUNT(gc.id) FILTER (WHERE gc.status = 'available') AS available_copies
		FROM games g
		LEFT JOIN game_copies gc ON gc.game_id = g.id
		WHERE g.id = $1
		GROUP BY g.id`

	var gd GameDetail
	err := s.pool.QueryRow(ctx, query, gameID).Scan(
		&gd.Game.ID, &gd.Game.Title, &gd.Game.IgdbID, &gd.Game.Platform,
		&gd.Game.Summary, &gd.Game.CoverURL, &gd.Game.SourceMagazine, &gd.Game.AcquiredAt,
		&gd.TotalCopies, &gd.AvailableCopies,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get game detail: %w", err)
	}

	// Total rental count for this game.
	s.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM rentals r
		JOIN game_copies gc ON gc.id = r.copy_id
		WHERE gc.game_id = $1`, gameID).Scan(&gd.TotalRentals)

	// Top renter for this game.
	s.pool.QueryRow(ctx, `
		SELECT m.profile_name, COUNT(*) AS cnt
		FROM rentals r
		JOIN game_copies gc ON gc.id = r.copy_id
		JOIN members m ON m.id = r.member_id
		WHERE gc.game_id = $1
		GROUP BY m.profile_name
		ORDER BY cnt DESC
		LIMIT 1`, gameID).Scan(&gd.TopRenterName, &gd.TopRenterCount)

	// Current renter (if any copy is currently rented).
	s.pool.QueryRow(ctx, `
		SELECT m.profile_name FROM rentals r
		JOIN game_copies gc ON gc.id = r.copy_id
		JOIN members m ON m.id = r.member_id
		WHERE gc.game_id = $1 AND r.returned_at IS NULL
		LIMIT 1`, gameID).Scan(&gd.CurrentRenter)

	return &gd, nil
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

// ProcessOverdueRentals auto-returns overdue rentals and penalizes members.
func (s *PostgresStore) ProcessOverdueRentals(ctx context.Context) (int, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx,
		`SELECT r.id, r.copy_id, r.member_id, m.profile_name, g.title
		 FROM rentals r
		 JOIN members m ON m.id = r.member_id
		 JOIN game_copies gc ON gc.id = r.copy_id
		 JOIN games g ON g.id = gc.game_id
		 WHERE r.returned_at IS NULL AND r.due_at < NOW()
		 FOR UPDATE OF r`)
	if err != nil {
		return 0, fmt.Errorf("failed to query overdue rentals: %w", err)
	}

	type overdueRental struct {
		rentalID   uuid.UUID
		copyID     uuid.UUID
		memberID   uuid.UUID
		memberName string
		gameTitle  string
	}
	var overdue []overdueRental
	for rows.Next() {
		var o overdueRental
		if err := rows.Scan(&o.rentalID, &o.copyID, &o.memberID, &o.memberName, &o.gameTitle); err != nil {
			rows.Close()
			return 0, fmt.Errorf("failed to scan overdue rental: %w", err)
		}
		overdue = append(overdue, o)
	}
	rows.Close()

	if len(overdue) == 0 {
		return 0, nil
	}

	for _, o := range overdue {
		_, err = tx.Exec(ctx,
			`UPDATE rentals SET returned_at = NOW() WHERE id = $1`, o.rentalID)
		if err != nil {
			return 0, fmt.Errorf("failed to auto-return rental %s: %w", o.rentalID, err)
		}

		_, err = tx.Exec(ctx,
			`UPDATE game_copies SET status = 'available' WHERE id = $1`, o.copyID)
		if err != nil {
			return 0, fmt.Errorf("failed to mark copy available %s: %w", o.copyID, err)
		}

		_, err = tx.Exec(ctx,
			`UPDATE members SET status = 'em_debito', late_count = late_count + 1 WHERE id = $1`,
			o.memberID)
		if err != nil {
			return 0, fmt.Errorf("failed to penalize member %s: %w", o.memberID, err)
		}

		if err := s.insertActivityTx(ctx, tx, "penalty", o.memberName, o.gameTitle); err != nil {
			return 0, fmt.Errorf("failed to insert penalty activity: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to commit overdue processing: %w", err)
	}

	return len(overdue), nil
}

// GetTopShameEntries returns the top N members with the most late returns.
func (s *PostgresStore) GetTopShameEntries(ctx context.Context, limit int) ([]ShameEntry, error) {
	query := `
		SELECT profile_name, late_count
		FROM members
		WHERE late_count > 0
		ORDER BY late_count DESC, profile_name ASC
		LIMIT $1`

	rows, err := s.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query shame entries: %w", err)
	}
	defer rows.Close()

	var entries []ShameEntry
	for rows.Next() {
		var e ShameEntry
		if err := rows.Scan(&e.ProfileName, &e.LateCount); err != nil {
			return nil, fmt.Errorf("failed to scan shame entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// RedeemMember resets a member's status from 'em_debito' to 'active'.
func (s *PostgresStore) RedeemMember(ctx context.Context, memberID uuid.UUID) error {
	tag, err := s.pool.Exec(ctx,
		`UPDATE members SET status = 'active' WHERE id = $1 AND status = 'em_debito'`,
		memberID)
	if err != nil {
		return fmt.Errorf("failed to redeem member: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("member not found or not in debt: %s", memberID)
	}
	return nil
}

// GetMemberStatus returns the current status of a member.
func (s *PostgresStore) GetMemberStatus(ctx context.Context, memberID uuid.UUID) (string, error) {
	var status string
	err := s.pool.QueryRow(ctx,
		`SELECT COALESCE(status, 'active') FROM members WHERE id = $1`, memberID).Scan(&status)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("member not found: %s", memberID)
		}
		return "", fmt.Errorf("failed to get member status: %w", err)
	}
	return status, nil
}

// insertActivityTx records an activity event within an existing transaction.
func (s *PostgresStore) insertActivityTx(ctx context.Context, tx pgx.Tx, eventType, memberName, gameTitle string) error {
	_, err := tx.Exec(ctx,
		`INSERT INTO activities (id, event_type, member_name, game_title, created_at)
		 VALUES ($1, $2, $3, $4, NOW())`,
		uuid.New(), eventType, memberName, gameTitle)
	if err != nil {
		return fmt.Errorf("failed to insert activity: %w", err)
	}
	return nil
}

// InsertActivity records an activity event using the connection pool.
func (s *PostgresStore) InsertActivity(ctx context.Context, eventType, memberName, gameTitle string) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO activities (id, event_type, member_name, game_title, created_at)
		 VALUES ($1, $2, $3, $4, NOW())`,
		uuid.New(), eventType, memberName, gameTitle)
	if err != nil {
		return fmt.Errorf("failed to insert activity: %w", err)
	}
	return nil
}

// ListRecentActivities returns the N most recent activity events.
func (s *PostgresStore) ListRecentActivities(ctx context.Context, limit int) ([]ActivityEntry, error) {
	query := `
		SELECT id, event_type, member_name, game_title, created_at
		FROM activities
		ORDER BY created_at DESC
		LIMIT $1`

	rows, err := s.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query activities: %w", err)
	}
	defer rows.Close()

	var result []ActivityEntry
	for rows.Next() {
		var a ActivityEntry
		if err := rows.Scan(&a.ID, &a.EventType, &a.MemberName, &a.GameTitle, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan activity: %w", err)
		}
		result = append(result, a)
	}
	return result, nil
}

// ListMemberActiveRentals returns active rentals for a specific member.
func (s *PostgresStore) ListMemberActiveRentals(ctx context.Context, memberID uuid.UUID) ([]MemberRental, error) {
	query := `
		SELECT r.id, g.title, g.cover_url, g.platform, r.rented_at, r.due_at,
		       (r.due_at < NOW()) AS is_overdue
		FROM rentals r
		JOIN game_copies gc ON gc.id = r.copy_id
		JOIN games g ON g.id = gc.game_id
		WHERE r.member_id = $1 AND r.returned_at IS NULL
		ORDER BY r.rented_at DESC`

	rows, err := s.pool.Query(ctx, query, memberID)
	if err != nil {
		return nil, fmt.Errorf("failed to query member rentals: %w", err)
	}
	defer rows.Close()

	var result []MemberRental
	for rows.Next() {
		var mr MemberRental
		var rentedAt, dueAt time.Time
		if err := rows.Scan(&mr.RentalID, &mr.GameTitle, &mr.CoverURL, &mr.Platform,
			&rentedAt, &dueAt, &mr.IsOverdue); err != nil {
			return nil, fmt.Errorf("failed to scan member rental: %w", err)
		}
		mr.RentedAt = rentedAt.Format("02/01/2006")
		mr.DueAt = dueAt.Format("02/01/2006")
		result = append(result, mr)
	}
	return result, nil
}

// CountOnTimeReturns counts completed rentals returned before or on the due date.
func (s *PostgresStore) CountOnTimeReturns(ctx context.Context, memberID uuid.UUID) (int, error) {
	var count int
	err := s.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM rentals
		 WHERE member_id = $1 AND returned_at IS NOT NULL AND returned_at <= due_at`,
		memberID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count on-time returns: %w", err)
	}
	return count, nil
}

// ReturnGameByMember returns a game, validating that the rental belongs to the given member.
// verdict stores the member's play status in the public_legacy column.
func (s *PostgresStore) ReturnGameByMember(ctx context.Context, rentalID, memberID uuid.UUID, verdict string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var copyID uuid.UUID
	err = tx.QueryRow(ctx,
		`SELECT copy_id FROM rentals
		 WHERE id = $1 AND member_id = $2 AND returned_at IS NULL`,
		rentalID, memberID).Scan(&copyID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("rental not found or does not belong to this member")
		}
		return fmt.Errorf("failed to find rental: %w", err)
	}

	_, err = tx.Exec(ctx,
		`UPDATE rentals SET returned_at = NOW(), public_legacy = $2 WHERE id = $1`,
		rentalID, verdict)
	if err != nil {
		return fmt.Errorf("failed to update rental: %w", err)
	}

	_, err = tx.Exec(ctx, `UPDATE game_copies SET status = 'available' WHERE id = $1`, copyID)
	if err != nil {
		return fmt.Errorf("failed to update copy status: %w", err)
	}

	return tx.Commit(ctx)
}

// GetRentalGameTitle returns the game title for a rental (used for activity logging).
func (s *PostgresStore) GetRentalGameTitle(ctx context.Context, rentalID uuid.UUID) (string, error) {
	var title string
	err := s.pool.QueryRow(ctx,
		`SELECT g.title FROM rentals r
		 JOIN game_copies gc ON gc.id = r.copy_id
		 JOIN games g ON g.id = gc.game_id
		 WHERE r.id = $1`, rentalID).Scan(&title)
	if err != nil {
		return "", fmt.Errorf("failed to get rental game title: %w", err)
	}
	return title, nil
}

// ListCompletedGameIDs returns game IDs that the member has completed ("zerei").
func (s *PostgresStore) ListCompletedGameIDs(ctx context.Context, memberID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT DISTINCT gc.game_id
		 FROM rentals r
		 JOIN game_copies gc ON gc.id = r.copy_id
		 WHERE r.member_id = $1
		   AND r.returned_at IS NOT NULL
		   AND r.public_legacy = 'zerei'`, memberID)
	if err != nil {
		return nil, fmt.Errorf("failed to query completed games: %w", err)
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan completed game id: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}
