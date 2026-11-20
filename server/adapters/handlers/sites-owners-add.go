package handlers

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/zeiss/builder/pkg/sbac"
)

// AddOwnerInput is the input for the AddOwners operation.
type AddOwnerInput struct {
	// SiteName is the name of the site to add owners to.
	SiteName string `json:"name" path:"name" example:"fizzy-buzzy" doc:"The name of the site to retrieve."`
	// Body is the owner to add.
	Body struct {
		// Email is the email of the owner to add.
		Email string `json:"email" example:"user@example.com" doc:"The email of the owner to add." required:"true" email:"true"`
	}
}

// AddOwnerOutput is the output for the AddOwner operation.
type AddOwnerOutput struct{}

// AddOwners adds owners to a site.
func (h *sitesHandler) AddOwner(ctx context.Context, input *AddOwnerInput) (*AddOwnerOutput, error) {
	err := h.permissionsCtrl.Allowed(ctx, input.SiteName, sbac.PermissionAdmin)
	if err != nil {
		return nil, err
	}

	err = h.sitesCtrl.AddOwner(ctx, input.SiteName, input.Body.Email)
	if err != nil {
		return nil, err
	}

	return &AddOwnerOutput{}, nil
}

// RegisterAddOwner registers the AddOwner operation.
func (h *sitesHandler) RegisterAddOwner(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "addOwner",
		DefaultStatus: 200,
		Method:        "PATCH",
		Path:          "/sites/{name}/owner",
		Summary:       "Add owners to a site",
		Description:   "Adds owners to a site by its name.",
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
				Description: "The owners have been added to the site.",
			},
			"401": {
				Description: "Unauthorized.",
			},
			"404": {
				Description: "No site found.",
			},
			"500": {
				Description: "An internal server error occurred.",
			},
		},
	}, h.AddOwner)
}
