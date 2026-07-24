package controllers

import (
	"context"

	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/server/ports"
)

var _ ports.Sites = (*SitesController)(nil)

// SitesController is a controller for managing sites.
type SitesController struct {
	db ports.Sites
}

// NewSitesController creates a new SitesController.
func NewSitesController(db ports.Sites) *SitesController {
	return &SitesController{db: db}
}

// CreateSite creates a new site with the given name.
func (c *SitesController) CreateSite(ctx context.Context, site *models.Site) error {
	err := c.db.CreateSite(ctx, site)
	if err != nil {
		return err
	}

	return nil
}

// GetSite returns the site with the given name.
func (c *SitesController) GetSite(ctx context.Context, site *models.Site) (models.Site, error) {
	return c.db.GetSite(ctx, site)
}

// UpdateSite updates the site with the given name.
func (c *SitesController) UpdateSite(ctx context.Context, site *models.Site) error {
	err := c.db.UpdateSite(ctx, site)
	if err != nil {
		return err
	}

	return nil
}

// DeleteSite deletes the site with the given name.
func (c *SitesController) DeleteSite(ctx context.Context, site *models.Site) error {
	err := c.db.DeleteSite(ctx, site)
	if err != nil {
		return err
	}

	return nil
}
