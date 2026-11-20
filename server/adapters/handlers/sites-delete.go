package handlers

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/zeiss/builder/pkg/sbac"
)

// DeleteSiteInput is the input for the DeleteSite operation.
type DeleteSiteInput struct {
	// Name is the name of the site to delete.
	Name string `json:"name" path:"name" example:"fizzy-buzzy" doc:"The name of the site to delete."`
}

// DeleteSiteOutput is the output for the DeleteSite operation.
type DeleteSiteOutput struct{}

// DeleteSite deletes a site by its name.
func (h *sitesHandler) DeleteSite(ctx context.Context, input *DeleteSiteInput) (*DeleteSiteOutput, error) {
	err := h.permissionsCtrl.Allowed(ctx, input.Name, sbac.PermissionAdmin)
	if err != nil {
		return nil, err
	}

	err = h.sitesCtrl.DeleteSite(ctx, input.Name)
	if err != nil {
		return nil, err
	}

	return &DeleteSiteOutput{}, nil
}

// RegisterDeleteSite registers the DeleteSite operation with the given Huma API.
func (h *sitesHandler) RegisterDeleteSite(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "deleteSite",
		DefaultStatus: 200,
		Method:        "DELETE",
		Path:          "/sites/{name}",
		Summary:       "Delete a site",
		Description:   "Deletes a site by its name.",
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
				Description: "The site has been deleted.",
			},
			"404": {
				Description: "No site found.",
			},
			"500": {
				Description: "An internal server error occurred.",
			},
		},
	}, h.DeleteSite)
}
