package models

import (
	"time"

	"gorm.io/gorm"
)

// Site represents a site.
type Site struct {
	// ID is the unique identifier of the site.
	ID string `path:"siteId" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id" example:"6c785b65-e689-4653-ae78-81621226c48c"`
	// Name is the name of the site.
	Name string `body:"name" gorm:"unique;uniqueIndex:idx_owned_by,sort:desc" json:"name" example:"fizzy-buzzy"`
	// CreatedAt is the time the deployment was created.
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	// UpdatedAt is the time the deployment was last updated.
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	// DeletedAt is the time the deployment was deleted.
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
	// Owner is the email of the account that owns the site.
	OwnedBy *string `json:"ownedBy" gorm:"type:string;uniqueIndex:idx_owned_by" example:"indy@jones.com"`
}
