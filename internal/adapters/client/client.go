package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/ports"
	"github.com/zeiss/builder/pkg/apis"
)

var (
	ErrCreateSite = fmt.Errorf("failed to create the site")
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

// Exists is a method that returns true if a site exists.
func (c *client) Exists(ctx context.Context, site *models.Site) (bool, error) {
	res, err := c.apis.GetSite(ctx, site.Name)
	if err != nil {
		return false, err
	}

	defer res.Body.Close()

	return res.StatusCode != http.StatusNotFound, nil
}

// Get is a method that returns a site by ID.
func (c *client) Get(_ context.Context, _ *models.Site) error {
	return nil
}

// Create is a method that creates a new site.
func (c *client) Create(ctx context.Context, site *models.Site) error {
	resp, err := c.apis.CreateSite(ctx, site.Name)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK {
		return nil
	}

	return ErrCreateSite
}

// Delete is a method that deletes a site.
func (c *client) Delete(ctx context.Context, site *models.Site) error {
	return nil
}

// UploadFile is a method that uploads a file to a site.
func (c *client) UploadFile(ctx context.Context, site *models.Site, file string) error {
	body, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(body)
	contentType := http.DetectContentType(body)

	resp, err := c.apis.PutSiteFileWithBody(ctx, site.Name, file, contentType, reader)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
