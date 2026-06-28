package ports

import "context"

// Sites is an interface for managing sites.
// This interface wraps the S3 interface and provides additional site-specific functionality.
type Sites interface {
	// CreateSite creates a new site with the given name.
	CreateSite(ctx context.Context, name string) error
	// GetSite returns the site with the given name.
	GetSite(ctx context.Context, name string) (bool, error)
	// PutObject uploads a file to the S3 bucket with the given name.
	PutObject(ctx context.Context, site, name string, content []byte) error
	// DeleteObject deletes the object with the given name from the S3 bucket with the given name.
	DeleteObject(ctx context.Context, site, name string) error
}

// Repo is an interface for managing repositories.
type Repo interface {
	// Sites is an interface for managing sites.
	Sites
}
