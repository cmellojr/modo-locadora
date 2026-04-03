package models

import (
	"time"

	"github.com/google/uuid"
)

// Club role constants.
const (
	ClubRoleAdmin  = "admin"
	ClubRoleMember = "member"
)

// Club represents a gaming community/group (turma).
type Club struct {
	ID          uuid.UUID
	Name        string
	Description string
	BadgeURL    string
	WebsiteURL  string
	CreatedBy   uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
