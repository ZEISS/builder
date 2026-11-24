package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/oapi-codegen/oapi-codegen/v2/pkg/securityprovider"
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/oidc"
	"github.com/zeiss/builder/internal/ports"
	"github.com/zeiss/builder/internal/services"
	"github.com/zeiss/builder/pkg/apis"
	"github.com/zeiss/builder/pkg/sbac"
	"github.com/zeiss/builder/pkg/utils"
	discovery "github.com/zeiss/fiber-discovery"
	"github.com/zeiss/pkg/cast"
)

var (
	ErrUnimplemented = fmt.Errorf("not implemented")
	ErrSiteExists    = fmt.Errorf("site already exists")
)

var (
	_ ports.SitesRepository = (*client)(nil)
	_ ports.FilesRepository = (*client)(nil)
)

type client struct {
	apis *apis.ClientWithResponses
}

// New creates a new client.
func New(api *apis.ClientWithResponses) *client {
	return &client{
		apis: api,
	}
}

// Create is a method that creates a new site.
func (c *client) Create(ctx context.Context, site *sbac.Site) error {
	body := apis.CreateSiteJSONRequestBody{Name: site.Name}
	resp, err := c.apis.CreateSite(ctx, body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		return ErrSiteExists
	}

	if resp.StatusCode > http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// GetSite is a method that gets a site by name.
func (c *client) GetSite(ctx context.Context, siteName string) (sbac.Site, error) {
	params := &apis.GetSiteParams{Name: siteName}
	site := sbac.Site{}

	resp, err := c.apis.GetSiteWithResponse(ctx, params)
	if err != nil {
		return site, err
	}

	if resp.StatusCode() > http.StatusOK {
		return site, err
	}

	site.Name = resp.JSON200.Site.Name
	site.CreatedAt = resp.JSON200.Site.CreatedAt
	site.UpdatedAt = resp.JSON200.Site.UpdatedAt

	return site, nil
}

// UploadFile is a method that uploads a file to a site.
func (c *client) UploadFile(ctx context.Context, name, cwd, file string) error {
	body, err := os.ReadFile(filepath.Clean(filepath.Join(cwd, file)))
	if err != nil {
		return err
	}

	reader := bytes.NewReader(body)
	contentType := http.DetectContentType(body)

	params := &apis.UploadFileParams{
		FileName: file,
	}

	res, err := c.apis.UploadFileWithBody(ctx, name, params, contentType, reader)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return nil
}

// DeleteSite is a method that deletes a site by name.
func (c *client) DeleteSite(ctx context.Context, sitName string) error {
	res, err := c.apis.DeleteSiteWithResponse(ctx, sitName)
	if err != nil {
		return err
	}

	if res.StatusCode() > http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode())
	}

	return nil
}

// AddOwner is a method that adds an owner to a site.
func (c *client) AddOwner(ctx context.Context, siteName, email string) error {
	body := apis.AddOwnerJSONRequestBody{Email: email}
	resp, err := c.apis.AddOwner(ctx, siteName, body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode > http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// ListOwners is a method that lists all owners of a site.
func (c *client) ListOwners(ctx context.Context, siteName string) ([]apis.User, error) {
	var owners []apis.User

	resp, err := c.apis.ListOwnersWithResponse(ctx, siteName)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() > http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	owners = append(owners, cast.Value(resp.JSON200.Owners)...)

	return owners, nil
}

// RemoveOwner is a method that removes an owner from a site.
func (c *client) RemoveOwner(ctx context.Context, siteName, email string) error {
	body := apis.RemoveOwnerJSONRequestBody{Email: email}
	resp, err := c.apis.RemoveOwner(ctx, siteName, body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode > http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// SetupClient is a function that sets up the client for the application.
func SetupClient(ctx context.Context, srv services.Service, uri string) (*client, error) {
	wellknownCfg, err := GetWellknownConfig(ctx, uri)
	if err != nil {
		return nil, err
	}

	oidcProvider, err := oidc.SetupOidcProvider(ctx, wellknownCfg.OidcIssuer, config.DefaultClientID)
	if err != nil {
		return nil, err
	}

	account := &models.Account{Current: true}
	if err := srv.GetAccount(ctx, account); err != nil {
		return nil, err
	}

	token := account.ToToken()
	if err := oidcProvider.Refresh(ctx, token); err != nil {
		return nil, err
	}
	account.UpdateToken(token)

	err = srv.UpdateAccount(ctx, account)
	if err != nil {
		return nil, err
	}

	bearer, err := securityprovider.NewSecurityProviderBearerToken(account.IDToken)
	if err != nil {
		return nil, err
	}

	api, err := apis.NewClientWithResponses(wellknownCfg.ApiURL, apis.WithRequestEditorFn(bearer.Intercept))
	if err != nil {
		return nil, err
	}

	c := New(api)

	return c, nil
}

// GetWellknownConfig discovers the well-known configuration for the given URI.
func GetWellknownConfig(ctx context.Context, uri string) (*utils.WellknownConfig, error) {
	wellknownCfg := &utils.WellknownConfig{}

	uri, err := url.JoinPath(uri, ".well-known", "builder-configuration")
	if err != nil {
		return nil, err
	}

	discv := discovery.NewClient()
	err = discv.Discover(ctx, uri, wellknownCfg)
	if err != nil {
		return nil, err
	}

	return wellknownCfg, nil
}
