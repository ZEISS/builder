package database

import (
	"context"

	"github.com/zeiss/builder/pkg/sbac"
	"github.com/zeiss/builder/server/ports"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
func (db *Database) CreateSite(ctx context.Context, name, email string) (sbac.Site, error) {
	site := sbac.Site{Name: name, Slug: name}

	err := db.conn.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := gorm.G[sbac.Site](tx).Create(ctx, &site)
		if err != nil {
			return err
		}

		user := sbac.User{Email: email}
		if err := tx.WithContext(ctx).Where(sbac.User{Email: email}).FirstOrCreate(&user).Error; err != nil {
			return err
		}

		role, err := gorm.G[sbac.Role](tx).Where("name = ?", sbac.RoleOwner).First(ctx)
		if err != nil {
			return err
		}

		if err := tx.WithContext(ctx).Model(&user).Association("Sites").Append(&site); err != nil {
			return err
		}

		userRole := sbac.UserRole{
			UserID: user.ID,
			RoleID: role.ID,
			SiteID: site.ID,
		}
		if err := gorm.G[sbac.UserRole](tx, clause.OnConflict{DoNothing: true}).Create(ctx, &userRole); err != nil {
			return err
		}

		return nil
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

// ListOwners lists the owners of the site with the given name.
func (db *Database) ListOwners(ctx context.Context, name string) ([]*sbac.User, error) {
	site, err := gorm.G[sbac.Site](db.conn).Preload("Users", nil).Where("slug = (?)", name).First(ctx)
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

		user := sbac.User{Email: email}
		if err := tx.WithContext(ctx).Where(sbac.User{Email: email}).FirstOrCreate(&user).Error; err != nil {
			return err
		}

		role, err := gorm.G[sbac.Role](tx).Where("name = ?", sbac.RoleOwner).First(ctx)
		if err != nil {
			return err
		}

		if err := tx.WithContext(ctx).Model(&user).Association("Sites").Append(&site); err != nil {
			return err
		}

		userRole := sbac.UserRole{
			UserID: user.ID,
			RoleID: role.ID,
			SiteID: site.ID,
		}
		if err := gorm.G[sbac.UserRole](tx, clause.OnConflict{DoNothing: true}).Create(ctx, &userRole); err != nil {
			return err
		}

		return nil
	})
}

// RemoveOwner removes owners from the site with the given name.
func (db *Database) RemoveOwner(ctx context.Context, siteName, email string) error {
	return db.conn.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		site, err := gorm.G[sbac.Site](tx).Where("slug = ?", siteName).First(ctx)
		if err != nil {
			return err
		}

		user, err := gorm.G[sbac.User](tx).Where("email = ?", email).First(ctx)
		if err != nil {
			return err
		}

		err = tx.Model(&user).Association("Sites").Delete(&site)
		if err != nil {
			return err
		}

		err = tx.Model(&user).Association("Roles").Delete(&site)
		if err != nil {
			return err
		}

		return nil
	})
}

// UpdateSite updates the site with the given name.
func (db *Database) UpdateSite(ctx context.Context, site *sbac.Site) error {
	return nil
}
