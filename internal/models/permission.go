package models

import "gorm.io/gorm"

// Permission represents a permission on a site.
type Permission struct {
	gorm.Model
	// ID is the unique identifier for the permission.
	ID int `json:"id" gorm:"primaryKey"`
	// Name is the name of the permission.
	Name string `json:"name" gorm:"uniqueIndex:idx_permission_name"`
}
