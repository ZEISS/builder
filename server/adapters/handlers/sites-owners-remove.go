package handlers

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/zeiss/builder/pkg/sbac"
	humax "github.com/zeiss/fiber-authz/v3/oas/huma"
	"github.com/zeiss/pkg/utilx"
)

// RemoveOwnerInput is the input for the RemoveOwner operation.
type RemoveOwnerInput struct {
	// SiteName is the name of the site to remove owners from.
	SiteName string `json:"name" path:"name" example:"fizzy-buzzy" doc:"The name of the site to retrieve."`
	// Body is the owner to add.
	Body struct {
		// Email is the email of the owner to add.
		Email string `json:"email" example:"user@example.com" doc:"The email of the owner to add." required:"true" email:"true"`
	}
}

// RemoveOwnerOutput is the output for the RemoveOwner operation.
type RemoveOwnerOutput struct{}

// RemoveOwner removes owners from a site.
func (h *sitesHandler) RemoveOwner(ctx context.Context, input *RemoveOwnerInput) (*RemoveOwnerOutput, error) {
	err := h.permissionsCtrl.Allowed(ctx, input.SiteName, sbac.PermissionAdmin)
	if err != nil {
		return nil, err
	}

	email, err := humax.GetClaim(ctx, "email")
	if err != nil {
		return nil, err
	}

	if utilx.Equal(email, input.Body.Email) {
		return nil, huma.Error405MethodNotAllowed("not allowed")
	}

	err = h.sitesCtrl.RemoveOwner(ctx, input.SiteName, input.Body.Email)
	if err != nil {
		return nil, err
	}

	return &RemoveOwnerOutput{}, nil
}

// RegisterAddOwner registers the AddOwner operation.
func (h *sitesHandler) RegisterRemoveOwner(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "removeOwner",
		DefaultStatus: 200,
		Method:        "DELETE",
		Path:          "/sites/{name}/owner",
		Summary:       "Remove owners from a site",
		Description:   "Removes owners from a site by its name.",
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
				Description: "The owners have been removed from the site.",
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
	}, h.RemoveOwner)
}
