package ports

import (
	"context"

	"github.com/zeiss/builder/internal/models"
)

// Sites is an interface for managing sites.
type Sites interface {
	// CreateSite creates a new site with the given name.
	CreateSite(ctx context.Context, site *models.Site) error
	// GetSiteByName returns the site with the given name.
	GetSiteByName(ctx context.Context, site *models.Site) (models.Site, error)
	// GetSiteById returns the site with the given id.
	GetSiteById(ctx context.Context, site *models.Site) (models.Site, error)
	// TransferSite transfers the site with the given name to the given owner.
	TransferSite(ctx context.Context, name, owner, newOwner string) error
	// UpdateOwnedBy changes the owner of the site with the given name.
	UpdateOwnedBy(ctx context.Context, name, newOwner string) error
	// UpdateSite updates the site with the given name.
	UpdateSite(ctx context.Context, site *models.Site) error
	// DeleteSite deletes the site with the given name.
	DeleteSite(ctx context.Context, site *models.Site) error
	// ListOwners lists the owners of the site with the given name.
	ListOwners(ctx context.Context, name string) ([]models.User, error)
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
