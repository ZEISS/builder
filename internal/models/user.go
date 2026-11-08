package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User is the model for a user.
type User struct {
	// ID is the unique identifier of the user.
	ID uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	// Email is the email of the user.
	Email string `json:"email" gorm:"unique" validate:"required,email"`
	// Ownerships is the list of projects that the user owns.
	Ownerships []Ownership

	gorm.Model
}
