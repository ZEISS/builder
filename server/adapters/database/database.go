package database

import (
	"context"
	"fmt"

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
	return db.conn.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := gorm.G[models.Site](tx).Create(ctx, site)
		if err != nil {
			return err
		}

		return nil
	})
}

// GetSiteByName returns the site with the given name.
func (db *Database) GetSiteByName(ctx context.Context, site *models.Site) (models.Site, error) {
	return gorm.G[models.Site](db.conn).Preload("Ownership.User", nil).Where("name = ?", site.Name).First(ctx)
}

// GetSiteById returns the site with the given id.
func (db *Database) GetSiteById(ctx context.Context, site *models.Site) (models.Site, error) {
	return gorm.G[models.Site](db.conn).Where("id = ?", site.ID).Where("owned_by = ?", site.OwnedBy).First(ctx)
}

// ListOwners lists the owners of the site with the given name.
func (db *Database) ListOwners(ctx context.Context, name string) ([]models.User, error) {
	site, err := gorm.G[models.Site](db.conn).Preload("Ownership.User", nil).Where("name = ?", name).First(ctx)
	if err != nil {
		return nil, err
	}

	owners := make([]models.User, 0, len(site.Ownership))
	for _, ownership := range site.Ownership {
		owners = append(owners, ownership.User)
	}

	return owners, nil
}

// TransferSite transfers the site with the given name to the given owner.
func (db *Database) TransferSite(ctx context.Context, name, owner, newOwner string) error {
	return db.conn.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		rows, err := gorm.G[models.Site](tx).Where("name = ? AND owned_by = ?", name, owner).Update(ctx, "owned_by", newOwner)
		if err != nil {
			return err
		}

		if rows == 0 {
			return fmt.Errorf("site could not be transferred: no site found with name %s and owned by %s", name, owner)
		}

		return nil
	})
}

// UpdateOwnedBy updates the owner of the site with the given name.
func (db *Database) UpdateOwnedBy(ctx context.Context, name, newOwner string) error {
	return db.conn.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		rows, err := gorm.G[models.Site](tx).Where("name = ? AND owned_by IS NULL", name).Update(ctx, "owned_by", newOwner)
		if err != nil {
			return err
		}

		if rows == 0 {
			return fmt.Errorf("site owner could not be updated: %s", name)
		}

		return nil
	})
}

// UpdateSite updates the site with the given name.
func (db *Database) UpdateSite(ctx context.Context, site *models.Site) error {
	return nil
}

// DeleteSite deletes the site with the given name.
func (db *Database) DeleteSite(ctx context.Context, site *models.Site) error {
	return nil
}
