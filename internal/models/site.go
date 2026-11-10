package models

import (
	"time"

	"github.com/google/uuid"
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
	OwnedBy *string `json:"ownedBy" gorm:"type:string;uniqueIndex:idx_owned_by;default:NULL" example:"indy@jones.com"`
	// Ownership is the list of ownerships of the site.
	Ownership []Ownership `gorm:"polymorphic:Ownable;polymorphicType:OwnableType;polymorphicId:OwnableID;polymorphicValue:site"`
}

// BeforeCreate is a GORM hook that is called before a site is created.
func (s *Site) BeforeCreate(tx *gorm.DB) error {
	s.ID = uuid.New().String()
	return nil
}
