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

// SitesController contains the methods for sites operations.
type SitesController interface {
	// Create is a method that creates a new site.
	Create(ctx context.Context, site *models.Site) error
	// GetSite is a method that gets a site by name.
	GetSite(ctx context.Context, name string) (models.Site, error)
}
