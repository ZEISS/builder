package models

import (
	"time"

	"gorm.io/gorm"
)

// Deployment represents a deployment of a site.
type Deployment struct {
	// ID is the unique identifier of the deployment.
	ID string `gorm:"primaryKey"`
	// Site is the site that this deployment is for.
	Site *Site `gorm:"foreignKey:ID"`
	// CreatedAt is the time the deployment was created.
	CreatedAt time.Time `gorm:"autoCreateTime"`
	// UpdatedAt is the time the deployment was last updated.
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	// DeletedAt is the time the deployment was deleted.
	DeletedAt gorm.DeletedAt
}
