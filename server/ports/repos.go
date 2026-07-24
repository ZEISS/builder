package ports

import (
	"context"

	"github.com/zeiss/builder/internal/models"
)

// Sites is an interface for managing sites.
type Sites interface {
	// CreateSite creates a new site with the given name.
	CreateSite(ctx context.Context, site *models.Site) error
	// GetSite returns the site with the given name.
	GetSite(ctx context.Context, site *models.Site) (models.Site, error)
	// UpdateSite updates the site with the given name.
	UpdateSite(ctx context.Context, site *models.Site) error
	// DeleteSite deletes the site with the given name.
	DeleteSite(ctx context.Context, site *models.Site) error
}

// Files is an interface for managing files.
type Files interface {
	// UploadFile uploads a file to the given path.
	UploadFile(ctx context.Context, site *models.Site, file *models.File) error
}

// Repo is an interface for managing repositories.
type Repo interface {
	// Sites is an interface for managing sites.
	Sites
}
