package models

import (
	"time"

	"github.com/google/uuid"
)

// Rental represents a rental transaction.
type Rental struct {
	ID           uuid.UUID
	MemberID     uuid.UUID
	CopyID       uuid.UUID
	RentedAt     time.Time
	DueAt        time.Time
	ReturnedAt   *time.Time
	PersonalNote string // Private notes by the member
	PublicLegacy string // Notes left on the back of the cartridge (publicly visible)
}
