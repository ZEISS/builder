package handlers

import (
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"github.com/zeiss/builder/pkg/sbac"
	humax "github.com/zeiss/fiber-authz/v3/oas/huma"
	"gorm.io/gorm"
)

// CreateSiteInput is the input for the CreateSite operation.
type CreateSiteInput struct {
	Body struct {
		Name string `body:"name" json:"name" example:"fizzy-buzzy" doc:"The name of the site (e.g. fizzy-buzzy)."`
	}
}

// CreateSiteOutput is the output for the CreateSite operation.
type CreateSiteOutput struct {
	Body sbac.Site
}

// CreateSite creates a new site with the given name.
func (h *sitesHandler) CreateSite(ctx context.Context, input *CreateSiteInput) (*CreateSiteOutput, error) {
	email, err := humax.GetClaim(ctx, "email")
	if err != nil {
		return nil, huma.Error401Unauthorized("unauthorized", err)
	}

	site, err := h.sitesCtrl.CreateSite(ctx, input.Body.Name, email)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return nil, huma.Error409Conflict("duplicate site", err)
	}

	if err != nil {
		return nil, err
	}

	return &CreateSiteOutput{Body: site}, nil
}

// Register the handler.
func (h *sitesHandler) RegisterCreateSite(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "createSite",
		DefaultStatus: 200,
		Method:        "POST",
		Path:          "/sites",
		Summary:       "Create a new site",
		Description:   "Creates a new site in the builder. This will create a new site folder.",
		Tags:          []string{"Sites"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "The site has been created.",
			},
			"400": {
				Description: "The request was invalid.",
			},
			"409": {
				Description: "The site already exists.",
			},
			"500": {
				Description: "An internal server error occurred.",
			},
		},
	}, h.CreateSite)
}
