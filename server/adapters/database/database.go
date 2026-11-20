package database

import (
	"context"
	"fmt"

	"github.com/zeiss/builder/pkg/sbac"
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
func (db *Database) CreateSite(ctx context.Context, name, ownedBy string) (sbac.Site, error) {
	site := sbac.Site{Name: name, Slug: name}

	err := db.conn.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := gorm.G[sbac.Site](tx).Create(ctx, &site)
		if err != nil {
			return err
		}

		user := sbac.User{Email: ownedBy, Sites: []sbac.Site{site}}
		if err := tx.Where(sbac.User{Email: ownedBy}).Attrs(sbac.User{Email: ownedBy}).FirstOrCreate(&user).Error; err != nil {
			return err
		}

		role, err := gorm.G[sbac.Role](tx).Where("name = ?", sbac.RoleOwner).First(ctx)
		if err != nil {
			return err
		}

		userRole := sbac.UserRole{SiteID: site.ID, UserID: user.ID, RoleID: role.ID}
		return gorm.G[sbac.UserRole](tx).Create(ctx, &userRole)
	})

	return site, err
}

// DeleteSite deletes the site with the given name.
func (db *Database) DeleteSite(ctx context.Context, siteName string) error {
	_, err := gorm.G[sbac.Site](db.conn).Where("slug = ?", siteName).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}

// GetSiteByName returns the site with the given name.
func (db *Database) GetSiteByName(ctx context.Context, name string) (sbac.Site, error) {
	return gorm.G[sbac.Site](db.conn).Where("slug = ?", name).First(ctx)
}

// GetSiteById returns the site with the given id.
func (db *Database) GetSiteById(ctx context.Context, site *sbac.Site) (sbac.Site, error) {
	return gorm.G[sbac.Site](db.conn).Where("id = ?", site.ID).First(ctx)
}

// ListOwners lists the owners of the site with the given name.
func (db *Database) ListOwners(ctx context.Context, name string) ([]*sbac.User, error) {
	site, err := gorm.G[sbac.Site](db.conn).Preload("Users", nil).Where("slug = ?", name).First(ctx)
	if err != nil {
		return nil, err
	}

	return site.Users, nil
}

// AddOwners adds owners to the site with the given name.
func (db *Database) AddOwner(ctx context.Context, siteName, email string) error {
	return db.conn.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		site, err := gorm.G[sbac.Site](tx).Where("slug = ?", siteName).First(ctx)
		if err != nil {
			return err
		}

		user := sbac.User{Email: email, Sites: []sbac.Site{site}}
		if err := tx.Where(sbac.User{Email: email}).Attrs(sbac.User{Email: email}).FirstOrCreate(&user).Error; err != nil {
			return err
		}

		role, err := gorm.G[sbac.Role](tx).Where("name = ?", sbac.RoleOwner).First(ctx)
		if err != nil {
			return err
		}

		userRole := sbac.UserRole{SiteID: site.ID, UserID: user.ID, RoleID: role.ID}
		err = gorm.G[sbac.UserRole](tx).Create(ctx, &userRole)
		if err != nil {
			return err
		}

		return nil
	})
}

// TransferSite transfers the site with the given name to the given owner.
func (db *Database) TransferSite(ctx context.Context, name, owner, newOwner string) error {
	return db.conn.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		rows, err := gorm.G[sbac.Site](tx).Where("name = ? AND owned_by = ?", name, owner).Update(ctx, "owned_by", newOwner)
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
		rows, err := gorm.G[sbac.Site](tx).Where("name = ? AND owned_by IS NULL", name).Update(ctx, "owned_by", newOwner)
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
func (db *Database) UpdateSite(ctx context.Context, site *sbac.Site) error {
	return nil
}
