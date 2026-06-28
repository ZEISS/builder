package handlers

import (
	"context"
	"net/url"

	"github.com/zeiss/builder/server/ports"

	"github.com/danielgtaylor/huma/v2"
)

type sitesHandler struct {
	ctrl ports.Sites
}

// NewSitesHandler creates a new SitesHandler with the given Sites controller.
func NewSitesHandler(ctrl ports.Sites) *sitesHandler {
	return &sitesHandler{ctrl: ctrl}
}

// CreateSiteInput is the input for the CreateSite operation.
type CreateSiteInput struct {
	Name string `path:"siteName" json:"name" example:"fizz-buzz"`
}

// CreateSiteOutput is the output for the PutSite operation.
type CreateSiteOutput struct{}

// CreateSite creates a new site with the given name.
func (h *sitesHandler) CreateSite(ctx context.Context, input *CreateSiteInput) (*CreateSiteOutput, error) {
	err := h.ctrl.CreateSite(ctx, input.Name)
	if err != nil {
		return nil, err
	}

	return &CreateSiteOutput{}, nil
}

// GetSiteInput is the input for the GetSite operation.
type GetSiteInput struct {
	Name string `path:"siteName" json:"name" example:"fizz-buzz"`
}

// GetSiteOutput is the output for the GetSite operation.
type GetSiteOutput struct{}

// GetSite gets a site by name.
func (h *sitesHandler) GetSite(ctx context.Context, input *GetSiteInput) (*GetSiteOutput, error) {
	exists, _ := h.ctrl.GetSite(ctx, input.Name)
	if !exists {
		return nil, huma.Error404NotFound("site not found")
	}

	return &GetSiteOutput{}, nil
}

// PutSiteFileInput is the input for the PutSiteFile operation.
type PutSiteFileInput struct {
	SiteName string `path:"siteName" json:"siteName" example:"fizz-buzz" required:"true"`
	Path     string `path:"path" json:"path" example:"index.html" required:"true"`
	RawBody  []byte `contentType:"application/octet-stream" body:"body" json:"body" example:""`
}

// PutSiteFileOutput is the output for the PutSiteFile operation.
type PutSiteFileOutput struct{}

// PutSiteFile uploads a file to the site with the given name.
func (h *sitesHandler) PutSiteFile(ctx context.Context, input *PutSiteFileInput) (*PutSiteFileOutput, error) {
	path, err := url.PathUnescape(input.Path)
	if err != nil {
		return nil, err
	}

	err = h.ctrl.PutObject(ctx, input.SiteName, path, input.RawBody)
	if err != nil {
		return nil, err
	}

	return &PutSiteFileOutput{}, nil
}

// Register registers the sites handler with the given Fiber app.
func (h *sitesHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "createSite",
		DefaultStatus: 200,
		Method:        "POST",
		Path:          "/sites/{siteName}",
		Summary:       "Create a new site",
		Description:   "Creates a new site in the builder. This will create a new site folder.",
		Tags:          []string{"Sites"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "The site has been created.",
			},
			"500": {
				Description: "An internal server error occurred.",
			},
		},
	}, h.CreateSite)

	huma.Register(api, huma.Operation{
		OperationID:   "getSite",
		DefaultStatus: 200,
		Method:        "GET",
		Path:          "/sites/{siteName}",
		Summary:       "Get a site",
		Description:   "Get a site by name.",
		Tags:          []string{"Sites"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "The site exists",
			},
			"404": {
				Description: "The site was not found",
			},
		},
	}, h.GetSite)

	huma.Register(api, huma.Operation{
		OperationID:   "putSiteFile",
		DefaultStatus: 201,
		Method:        "PUT",
		Path:          "/sites/{siteName}/files/{path}",
		Summary:       "Put a new file to the site",
		Description:   "Put a new file to the site.",
		Tags:          []string{"Sites"},
		RequestBody: &huma.RequestBody{
			Content: map[string]*huma.MediaType{
				"application/octet-stream": {
					Schema: &huma.Schema{
						Type:   "object",
						Format: "binary",
					},
				},
			},
		},
		Responses: map[string]*huma.Response{
			"201": {
				Description: "The file was uploaded successfully.",
			},
		},
	}, h.PutSiteFile)
}
