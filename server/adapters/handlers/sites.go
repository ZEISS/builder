package handlers

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/server/ports"
	"gorm.io/gorm"

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
	Body struct {
		Name string `body:"name" json:"name" example:"fizzy-buzzy" doc:"The name of the site (e.g. fizzy-buzzy)."`
	}
}

// CreateSiteOutput is the output for the CreateSite operation.
type CreateSiteOutput struct {
	Body *models.Site
}

// CreateSite creates a new site with the given name.
func (h *sitesHandler) CreateSite(ctx context.Context, input *CreateSiteInput) (*CreateSiteOutput, error) {
	site := &models.Site{
		ID:   uuid.New().String(),
		Name: input.Body.Name,
	}

	err := h.ctrl.CreateSite(ctx, site)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return nil, huma.Error409Conflict("duplicate site", err)
	}

	if err != nil {
		return nil, err
	}

	return &CreateSiteOutput{Body: site}, nil
}

// Register registers the sites handler with the given Fiber app.
func (h *sitesHandler) Register(api huma.API) {
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

	// huma.Register(api, huma.Operation{
	// 	OperationID:   "getSite",
	// 	DefaultStatus: 200,
	// 	Method:        "GET",
	// 	Path:          "/sites/{siteName}",
	// 	Summary:       "Get a site",
	// 	Description:   "Get a site by name.",
	// 	Tags:          []string{"Sites"},
	// 	Responses: map[string]*huma.Response{
	// 		"200": {
	// 			Description: "The site exists",
	// 		},
	// 		"404": {
	// 			Description: "The site was not found",
	// 		},
	// 	},
	// }, h.GetSite)

	// huma.Register(api, huma.Operation{
	// 	OperationID:   "putSiteFile",
	// 	DefaultStatus: 201,
	// 	Method:        "PUT",
	// 	Path:          "/sites/{siteName}/files/{path}",
	// 	Summary:       "Put a new file to the site",
	// 	Description:   "Put a new file to the site.",
	// 	Tags:          []string{"Sites"},
	// 	RequestBody: &huma.RequestBody{
	// 		Content: map[string]*huma.MediaType{
	// 			"application/octet-stream": {
	// 				Schema: &huma.Schema{
	// 					Type:   "object",
	// 					Format: "binary",
	// 				},
	// 			},
	// 		},
	// 	},
	// 	Responses: map[string]*huma.Response{
	// 		"201": {
	// 			Description: "The file was uploaded successfully.",
	// 		},
	// 	},
	// }, h.PutSiteFile)
}
