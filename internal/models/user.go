package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User is the model for a user.
type User struct {
	// ID is the unique identifier of the user.
	ID uuid.UUID `json:"ID" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	// Email is the email of the user.
	Email string `json:"email" gorm:"unique" validate:"required,email"`

	gorm.Model
}

// BeforeCreate is a GORM hook that is called before a user is created.
func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}
