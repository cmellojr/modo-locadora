package database

import (
	"context"
	"fmt"

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

// CreateMember persists a new member in the database.
func (s *PostgresStore) CreateMember(ctx context.Context, m *models.Member) error {
	query := `
		INSERT INTO members (id, profile_name, email, password_hash, favorite_console, joined_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := s.pool.Exec(ctx, query, m.ID, m.ProfileName, m.Email, m.PasswordHash, m.FavoriteConsole, m.JoinedAt)
	if err != nil {
		return fmt.Errorf("failed to create member: %w", err)
	}
	return nil
}

// GetGameByID retrieves a game by its ID.
func (s *PostgresStore) GetGameByID(ctx context.Context, id uuid.UUID) (*models.Game, error) {
	query := `SELECT id, title, igdb_id, platform, summary FROM games WHERE id = $1`

	var g models.Game
	err := s.pool.QueryRow(ctx, query, id).Scan(&g.ID, &g.Title, &g.IgdbID, &g.Platform, &g.Summary)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Or a specific error like ErrNotFound
		}
		return nil, fmt.Errorf("failed to get game: %w", err)
	}
	return &g, nil
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
