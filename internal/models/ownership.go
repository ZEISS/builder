package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Ownership represents the ownership of a resource by a user.
type Ownership struct {
	// UserID is the ID of the user who owns the resource.
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	// User is the user who owns the resource.
	User User `json:"user" gorm:"foreignKey:UserID"`
	// OwnableID is the ID of the resource that is owned.
	OwnableID uuid.UUID `json:"ownable_id" gorm:"type:uuid;not null"`
	// OwnableType is the type of the resource that is owned.
	OwnableType string `json:"ownable_type" gorm:"not null"`

	gorm.Model
}
