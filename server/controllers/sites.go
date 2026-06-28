package controllers

import (
	"context"

	"github.com/zeiss/builder/server/ports"
)

var _ ports.Sites = (*SitesController)(nil)

// SitesController is a controller for managing sites.
type SitesController struct {
	storage ports.Sites
}

// NewSitesController creates a new SitesController.
func NewSitesController(storage ports.Sites) *SitesController {
	return &SitesController{storage: storage}
}

// CreateSite creates a new site with the given name.
func (c *SitesController) CreateSite(ctx context.Context, name string) error {
	err := c.storage.CreateSite(ctx, name)
	if err != nil {
		return err
	}

	return nil
}

// GetSite returns the site with the given name.
func (c *SitesController) GetSite(ctx context.Context, name string) (bool, error) {
	return c.storage.GetSite(ctx, name)
}

// PutObject uploads a file to the S3 bucket with the given name.
func (c *SitesController) PutObject(ctx context.Context, site, name string, content []byte) error {
	err := c.storage.PutObject(ctx, site, name, content)
	if err != nil {
		return err
	}

	return nil
}

// DeleteObject deletes the object with the given name from the S3 bucket with the given name.
func (c *SitesController) DeleteObject(_ context.Context, _, _ string) error {
	return nil
}
