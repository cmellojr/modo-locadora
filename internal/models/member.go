package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	MemberStatusActive   = "active"
	MemberStatusEmDebito = "em_debito"
)

// Member represents a user in the system.
type Member struct {
	ID               uuid.UUID
	ProfileName      string
	Email            string
	PasswordHash     string
	FavoriteConsole  string
	MembershipNumber string
	Address          string
	Phone            string
	PasswordNotes    string
	Status           string // "active" or "em_debito"
	LateCount        int
	JoinedAt         time.Time
}
