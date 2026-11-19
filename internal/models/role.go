package models

import "gorm.io/gorm"

// Role is the role that the permission is assigned to.
type Role struct {
	gorm.Model
	// ID is the unique identifier for the role.
	ID uint `json:"id" gorm:"primaryKey"`
	// Name is the name of the role.
	Name string `json:"name"`
	// Permissions is the list of permissions assigned to the role.
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions"`
}
