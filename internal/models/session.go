package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Session represents a user session, storing session-specific data.
type Session struct {
	// ID is the unique identifier for the session
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:(gen_random_uuid())"`
	// Name is the name of the session
	Name string `gorm:"type:text"`
	// Description is the description of the session
	Description string `gorm:"type:text"`
	// AccountID is the ID of the account associated with the session
	AccountID uuid.UUID `gorm:"type:uuid"`
	// Account is the account associated with the session
	Account Account `gorm:"foreignKey:AccountID"`
	// CreatedAt is the creation time of the account.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the update time of the account.
	UpdatedAt time.Time `json:"updated_at"`
	// DeletedAt is the deletion time of the account.
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
