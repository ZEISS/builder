package database

import (
	"context"

	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/server/ports"

	"gorm.io/gorm"
)

var _ ports.Sites = (*Database)(nil)

// Database is the database adapter for the builder server.
type Database struct {
	conn *gorm.DB
}

// NewDatabase creates a new Database instance.
func NewDatabase(conn *gorm.DB) *Database {
	return &Database{conn: conn}
}

// CreateSite creates a new site with the given name.
func (db *Database) CreateSite(ctx context.Context, site *models.Site) error {
	return db.conn.Transaction(func(tx *gorm.DB) error {
		err := gorm.G[models.Site](tx).Create(ctx, site)
		if err != nil {
			return err
		}
		return nil
	})
}

// GetSiteByName returns the site with the given name.
func (db *Database) GetSiteByName(ctx context.Context, site *models.Site) (models.Site, error) {
	return gorm.G[models.Site](db.conn).Where("name = ?", site.Name).First(ctx)
}

// GetSiteById returns the site with the given id.
func (db *Database) GetSiteById(ctx context.Context, site *models.Site) (models.Site, error) {
	return gorm.G[models.Site](db.conn).Where("id = ?", site.ID).First(ctx)
}

// UpdateSite updates the site with the given name.
func (db *Database) UpdateSite(ctx context.Context, site *models.Site) error {
	return nil
}

// DeleteSite deletes the site with the given name.
func (db *Database) DeleteSite(ctx context.Context, site *models.Site) error {
	return nil
}
