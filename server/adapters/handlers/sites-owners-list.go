package handlers

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/zeiss/builder/pkg/sbac"
)

// ListOwnersInput is the input for the ListOwners operation.
type ListOwnersInput struct {
	Name string `json:"name" path:"name" example:"fizzy-buzzy" doc:"The name of the site to retrieve."`
}

// ListOwnersOutput is the output for the ListOwners operation.
type ListOwnersOutput struct {
	Body ListOwnersOutputBody `json:"body"`
}

// ListOwnersOutputBody is the body of the ListOwners output.
type ListOwnersOutputBody struct {
	Owners []*sbac.User `json:"owners"`
}

// ListOwners lists the owners of a site.
func (h *sitesHandler) ListOwners(ctx context.Context, input *ListOwnersInput) (*ListOwnersOutput, error) {
	err := h.permissionsCtrl.Allowed(ctx, input.Name, sbac.PermissionAdmin)
	if err != nil {
		return nil, err
	}

	owners, err := h.sitesCtrl.ListOwners(ctx, input.Name)
	if err != nil {
		return nil, err
	}

	body := ListOwnersOutputBody{Owners: owners}

	return &ListOwnersOutput{Body: body}, nil
}

// RegisterListOwners registers the ListOwners operation.
func (h *sitesHandler) RegisterListOwners(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "listOwners",
		DefaultStatus: 200,
		Method:        "GET",
		Path:          "/sites/{name}/owners",
		Summary:       "List owners of a site",
		Description:   "Lists the owners of a site by its name.",
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
				Description: "The sites have been queried.",
			},
			"404": {
				Description: "No site found.",
			},
			"500": {
				Description: "An internal server error occurred.",
			},
		},
	}, h.ListOwners)
}
