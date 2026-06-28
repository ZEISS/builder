package controllers

import (
	"context"

	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/ports"
)

var _ ports.SitesController = (*SitesController)(nil)

type SitesController struct {
	sitesRepo ports.SitesRepository
}

func NewSitesController(sitesRepo ports.SitesRepository) *SitesController {
	return &SitesController{
		sitesRepo: sitesRepo,
	}
}

// Get is a method that returns a site by ID.
func (c *SitesController) Get(ctx context.Context, site *models.Site) error {
	return c.sitesRepo.Get(ctx, site)
}

// Create is a method that creates a new site.
func (c *SitesController) Create(ctx context.Context, site *models.Site) error {
	return c.sitesRepo.Create(ctx, site)
}

// Delete is a method that deletes a site.
func (c *SitesController) Delete(ctx context.Context, site *models.Site) error {
	return c.sitesRepo.Delete(ctx, site)
}

// Exists is a method that returns true if a site exists.
func (c *SitesController) Exists(ctx context.Context, site *models.Site) (bool, error) {
	return c.sitesRepo.Exists(ctx, site)
}

// UploadFile is a method that uploads a file to a site.
func (c *SitesController) UploadFile(ctx context.Context, site *models.Site, file string) error {
	return c.sitesRepo.UploadFile(ctx, site, file)
}
