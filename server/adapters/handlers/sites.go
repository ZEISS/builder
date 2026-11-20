package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/zeiss/builder/server/ports"

	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
	humax "github.com/zeiss/fiber-authz/v3/oas/huma"
)

type sitesHandler struct {
	sitesCtrl       ports.Sites
	filesCtrl       ports.Files
	permissionsCtrl ports.Permissions
}

// NewSitesHandler creates a new SitesHandler with the given Sites controller.
func NewSitesHandler(sitesCtrl ports.Sites, filesCtrl ports.Files, permissionsCtrl ports.Permissions) *sitesHandler {
	return &sitesHandler{
		sitesCtrl:       sitesCtrl,
		filesCtrl:       filesCtrl,
		permissionsCtrl: permissionsCtrl,
	}
}

// UpdateSiteOwnedByInput is the input for the UpdateSiteOwnedBy operation.
type UpdateSiteOwnedByInput struct {
	// Name is the name of the site.
	Name string `path:"name" json:"name" example:"fizzy-buzzy" doc:"The name of the site."`
}

// UpdateSiteOwnedByOutput is the output of the UpdateSiteOwnedBy operation.
type UpdateSiteOwnedByOutput struct{}

// UpdateSiteOwnedBy handles the updating of a site's owner.
func (h *sitesHandler) UpdateSiteOwnedBy(ctx context.Context, input *UpdateSiteOwnedByInput) (*UpdateSiteOwnedByOutput, error) {
	token, ok := ctx.Value(humax.OidcTokenKey).(*jwt.Token)
	if !ok {
		return nil, errors.New("no token found")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims format")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("invalid token claims format")
	}

	err := h.sitesCtrl.UpdateOwnedBy(ctx, input.Name, email)
	if err != nil {
		return nil, errors.New("failed to claim site")
	}

	return &UpdateSiteOwnedByOutput{}, nil
}

// Register registers the sites handler with the given Huma API.
func Register(api huma.API, h *sitesHandler) {
	huma.Register(api, huma.Operation{
		OperationID:   "patchSiteOwnedBy",
		DefaultStatus: http.StatusOK,
		Method:        "PATCH",
		Path:          "/site/{name}/owner",
		Summary:       "Change the owner of a site",
		Description:   "Changes the owner of a site to the given email.",
		Tags:          []string{"Sites"},
		Parameters: []*huma.Param{
			{
				Name:        "name",
				In:          "path",
				Description: "The name of the site.",
				Required:    true,
				Schema: &huma.Schema{
					Type: "string",
				},
			},
		},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "The owner of the site has been updated.",
			},
			"400": {
				Description: "The request was invalid.",
			},
			"401": {
				Description: "Unauthorized.",
			},
			"403": {
				Description: "Forbidden.",
			},
			"404": {
				Description: "No site found.",
			},
			"500": {
				Description: "An internal server error occurred.",
			},
		},
	}, h.UpdateSiteOwnedBy)
}
