package handlers

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/zeiss/builder/pkg/sbac"
)

// GetSiteInput is the input for the GetSite operation.
type GetSiteInput struct {
	// Name is the name of the site to retrieve.
	Name string `json:"name" query:"name" example:"my-site" doc:"The name of the site to retrieve."`
}

// GetSiteOutput is the output for the GetSite operation.
type GetSiteOutput struct {
	Body GetSiteOutputBody `json:"body"`
}

// GetSiteOutputBody is the body of the GetSite output.
type GetSiteOutputBody struct {
	Site sbac.Site `json:"site"`
}

// GetSite retrieves a site by its name.
func (h *sitesHandler) GetSite(ctx context.Context, input *GetSiteInput) (*GetSiteOutput, error) {
	err := h.permissionsCtrl.Allowed(ctx, input.Name, sbac.PermissionAdmin)
	if err != nil {
		return nil, err
	}

	site, err := h.sitesCtrl.GetSiteByName(ctx, input.Name)
	if err != nil {
		return nil, err
	}

	body := GetSiteOutputBody{Site: site}

	return &GetSiteOutput{Body: body}, nil
}

// RegisterGetSite registers the GetSite operation with the given Huma API.
func (h *sitesHandler) RegisterGetSite(api huma.API) {
	// Get a site
	huma.Register(api, huma.Operation{
		OperationID:   "getSite",
		DefaultStatus: 200,
		Method:        "GET",
		Path:          "/sites",
		Summary:       "Get a site",
		Description:   "Gets a site by its name.",
		Tags:          []string{"Sites"},
		Parameters: []*huma.Param{
			{
				Name:        "name",
				In:          "query",
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
	}, h.GetSite)
}
