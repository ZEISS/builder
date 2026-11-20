package ports

import (
	"context"

	"github.com/zeiss/builder/pkg/sbac"
)

// SitesRepository contains the methods for sites operations.
type SitesRepository interface {
	// Create is a method that creates a new site.
	Create(ctx context.Context, site *sbac.Site) error
	// GetSite is a method that gets a site by name.
	GetSite(ctx context.Context, name string) (sbac.Site, error)
	// DeleteSite is a method that deletes a site.
	DeleteSite(ctx context.Context, name string) error
}

// SitesService contains the methods for sites operations.
type SitesService interface {
	// GetSite returns the site associated with the session.
	GetSite(ctx context.Context, site *sbac.Site) error
	// CreateSite creates a new site.
	CreateSite(ctx context.Context, site *sbac.Site) error
	// UpdateSite updates an existing site.
	UpdateSite(ctx context.Context, site *sbac.Site) error
	// DeleteSite deletes a site.
	DeleteSite(ctx context.Context, site *sbac.Site) error
	// UpdateSiteOwnedBy updates the owner of a site.
	UpdateSiteOwnedBy(ctx context.Context, name, owner string) error
}
