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
	ErrUnimplemented = fmt.Errorf("not implemented")
	ErrSiteExists    = fmt.Errorf("site already exists")
)

var _ ports.SitesRepository = (*client)(nil)

type client struct {
	apis *apis.ClientWithResponses
}

// New creates a new client.
func New(api *apis.ClientWithResponses) ports.SitesRepository {
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

// GetSite is a method that gets a site by name.
func (c *client) GetSite(ctx context.Context, name string) (models.Site, error) {
	params := &apis.GetSiteParams{Name: name}
	site := models.Site{}

	resp, err := c.apis.GetSiteWithResponse(ctx, params)
	if err != nil {
		return site, err
	}

	if resp.StatusCode() > http.StatusOK {
		return site, err
	}

	site.Name = resp.JSON200.Site.Name
	site.ID = resp.JSON200.Site.Id
	site.CreatedAt = resp.JSON200.Site.CreatedAt
	site.UpdatedAt = resp.JSON200.Site.UpdatedAt

	return site, nil
}

// UploadFile is a method that uploads a file to a site.
func (c *client) UploadFile(ctx context.Context, site *models.Site, file string) error {
	body, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(body)
	contentType := http.DetectContentType(body)

	params := &apis.UploadFileParams{
		FileName: file,
	}

	resp, err := c.apis.UploadFileWithBody(ctx, site.Name, params, contentType, reader)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
