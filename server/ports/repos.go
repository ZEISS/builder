package ports

import (
	"context"

	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/pkg/sbac"
)

// Sites is an interface for managing sites.
type Sites interface {
	// CreateSite creates a new site with the given name.
	CreateSite(ctx context.Context, name, ownedBy string) (sbac.Site, error)
	// GetSiteByName returns the site with the given name.
	GetSiteByName(ctx context.Context, name string) (sbac.Site, error)
	// GetSiteById returns the site with the given id.
	GetSiteById(ctx context.Context, site *sbac.Site) (sbac.Site, error)
	// TransferSite transfers the site with the given name to the given owner.
	TransferSite(ctx context.Context, name, owner, newOwner string) error
	// UpdateOwnedBy changes the owner of the site with the given name.
	UpdateOwnedBy(ctx context.Context, name, newOwner string) error
	// UpdateSite updates the site with the given name.
	UpdateSite(ctx context.Context, site *sbac.Site) error
	// DeleteSite deletes the site with the given name.
	DeleteSite(ctx context.Context, site *sbac.Site) error
	// ListOwners lists the owners of the site with the given name.
	ListOwners(ctx context.Context, name string) ([]*sbac.User, error)
}

// Files is an interface for managing files.
type Files interface {
	// UploadFile uploads a file to the given path.
	UploadFile(ctx context.Context, site *sbac.Site, file *models.File) error
}

// Permissions is an interface for checking permissions.
type Permissions interface {
	// Allowed returns true if the user is allowed to perform the given permission on the site.
	Allowed(ctx context.Context, siteName, userEmail string) error
}

// Repo is an interface for managing repositories.
type Repo interface {
	// Sites is an interface for managing sites.
	Sites
	// Files is an interface for managing files.
	Files
	// Permissions is an interface for checking permissions.
	Permissions
}
