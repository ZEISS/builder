package ports

import (
	"context"

	"github.com/zeiss/builder/internal/models"
)

// SitesRepository contains the methods for sites operations.
type SitesRepository interface {
	// Create is a method that creates a new site.
	Create(ctx context.Context, site *models.Site) error
	// Deploy is a method that deploys a site.
	UploadFile(ctx context.Context, site *models.Site, file string) error
}

// SitesController contains the methods for sites operations.
type SitesController interface {
	// Create is a method that creates a new site.
	Create(ctx context.Context, site *models.Site) error
	// UploadFile is a method that uploads a file to a site.
	UploadFile(ctx context.Context, site *models.Site, file string) error
}
