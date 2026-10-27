package ports

import (
	"context"

	"github.com/zeiss/builder/internal/models"
)

// SitesRepository contains the methods for sites operations.
type SitesRepository interface {
	// Create is a method that creates a new site.
	Create(ctx context.Context, site *models.Site) error
	// GetSite is a method that gets a site by name.
	GetSite(ctx context.Context, name string) (models.Site, error)
}

// SitesService contains the methods for sites operations.
type SitesService interface {
	// GetSite returns the site associated with the session.
	GetSite(ctx context.Context, site *models.Site) error
	// CreateSite creates a new site.
	CreateSite(ctx context.Context, site *models.Site) error
	// UpdateSite updates an existing site.
	UpdateSite(ctx context.Context, site *models.Site) error
	// DeleteSite deletes a site.
	DeleteSite(ctx context.Context, site *models.Site) error
}
