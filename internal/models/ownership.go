package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Ownership represents the ownership of a resource by a user.
type Ownership struct {
	// ID is the ID of the ownership.
	ID uint8 `json:"ID" gorm:"type:integer;primaryKey;autoIncrement"`
	// UserID is the ID of the user that owns the resource.
	UserID uuid.UUID `json:"user_id" gorm:"type:integer;uniqueIndex:idx_site_owner"`
	// User is the user that owns the resource.
	User User `gorm:"foreignKey:UserID"`
	// OwnableID is the ID of the resource that is owned.
	OwnableID uuid.UUID `json:"ownable_id" gorm:"type:uuid;uniqueIndex:idx_site_owner"`
	// OwnableType is the type of the resource that is owned.
	OwnableType string `json:"ownable_type" gorm:"type:string;uniqueIndex:idx_site_owner"`

	gorm.Model
}
