package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	MemberStatusActive   = "active"
	MemberStatusInDebt = "in_debt"
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
	Status           string // "active" or "in_debt"
	LateCount        int
	JoinedAt         time.Time
}

// MemberTitle represents a member's earned progression title.
type MemberTitle struct {
	Label    string // Portuguese display label
	BadgeCSS string // CSS class for the badge color
}

// ComputeMemberTitle determines the member title based on rental history.
// completedGames = count of distinct games with "completed" verdict.
// onTimeReturns  = count of rentals returned on or before due date.
func ComputeMemberTitle(completedGames, onTimeReturns int) MemberTitle {
	switch {
	case completedGames >= 5:
		return MemberTitle{Label: "Dono da Calcada", BadgeCSS: "is-title-legend"}
	case onTimeReturns >= 25:
		return MemberTitle{Label: "Socio Ouro", BadgeCSS: "is-title-gold"}
	case onTimeReturns >= 10:
		return MemberTitle{Label: "Socio Prata", BadgeCSS: "is-title-silver"}
	default:
		return MemberTitle{Label: "Socio Novato", BadgeCSS: "is-title-novice"}
	}
}
