package models

import (
	"time"

	"gorm.io/gorm"
)

// Site represents a site.
type Site struct {
	// ID is the unique identifier of the site.
	ID string `path:"siteId" gorm:"primaryKey" json:"id" example:"6c785b65-e689-4653-ae78-81621226c48c"`
	// Name is the name of the site.
	Name string `body:"name" gorm:"unique" json:"name" example:"fizzy-buzzy"`
	// CreatedAt is the time the deployment was created.
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	// UpdatedAt is the time the deployment was last updated.
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	// DeletedAt is the time the deployment was deleted.
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
}
