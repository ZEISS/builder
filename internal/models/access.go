package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Access is a many-to-many relationship between users and sites.
type Access struct {
	// ID is the unique identifier for the access.
	ID uint `gorm:"primaryKey"`
	// UserID is the ID of the user that owns the resource.
	UserID uuid.UUID `gorm:"primaryKey"`
	// SiteID is the ID of the site that is accessible.
	SiteID uuid.UUID `gorm:"primaryKey"`
	// RoleID is the ID of the role that is accessible.
	RoleID uint `gorm:"primaryKey"`
	// User is the user that owns the resource.
	User User `gorm:"constraint:OnDelete:CASCADE;"`
	// Site is the site that is accessible.
	Site Site `gorm:"constraint:OnDelete:CASCADE;"`
	// Role is the role that is accessible.
	Role Role `gorm:"constraint:OnDelete:CASCADE;"`

	gorm.Model
}
