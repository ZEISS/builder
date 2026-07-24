package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/ports"
	"github.com/zeiss/builder/pkg/apis"
)

var (
	ErrUnimplemented = fmt.Errorf("not implemented")
	ErrSiteExists    = fmt.Errorf("site already exists")
)

var _ ports.SitesRepository = (*client)(nil)

type client struct {
	apis *apis.Client
}

// New creates a new client.
func New(api *apis.Client) ports.SitesRepository {
	return &client{
		apis: api,
	}
}

// Create is a method that creates a new site.
func (c *client) Create(ctx context.Context, site *models.Site) error {
	body := apis.CreateSiteJSONRequestBody{Name: site.Name}
	resp, err := c.apis.CreateSite(ctx, body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		return ErrSiteExists
	}

	return nil
}

// UploadFile is a method that uploads a file to a site.
func (c *client) UploadFile(ctx context.Context, site *models.Site, file string) error {
	return ErrUnimplemented
}
