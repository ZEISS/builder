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
	sitesCtrl ports.Sites
	filesCtrl ports.Files
}

// NewSitesHandler creates a new SitesHandler with the given Sites controller.
func NewSitesHandler(sitesCtrl ports.Sites, filesCtrl ports.Files) *sitesHandler {
	return &sitesHandler{sitesCtrl: sitesCtrl, filesCtrl: filesCtrl}
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

	err := h.sitesCtrl.CreateSite(ctx, site)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return nil, huma.Error409Conflict("duplicate site", err)
	}

	if err != nil {
		return nil, err
	}

	return &CreateSiteOutput{Body: site}, nil
}

// UploadFileInput is the input for the UploadFile operation.
type UploadFileInput struct {
	Site struct {
		ID string `path:"siteId"`
	} `json:"site"`
	Body struct {
		Name        string `body:"fileName" json:"name" example:"fizzy-buzzy" doc:"The name of the site (e.g. fizzy-buzzy)."`
		File        []byte `body:"file" json:"file" example:"fizzy-buzzy" doc:"The name of the site (e.g. fizzy-buzzy)."`
		ContentType string `body:"contentType" json:"content_type" example:"text/plain" doc:"The content type of the file."`
	}
}

// UploadFileOutput is the output of the UploadFile operation.
type UploadFileOutput struct {
	Body string `json:"body"`
}

// UploadFile handles the upload of a file to the site.
func (h *sitesHandler) UploadFile(ctx context.Context, input *UploadFileInput) (*UploadFileOutput, error) {
	file := &models.File{
		Name:        input.Body.Name,
		Data:        input.Body.File,
		ContentType: input.Body.ContentType,
	}

	site, err := h.sitesCtrl.GetSite(ctx, &models.Site{
		ID: input.Site.ID,
	})
	if err != nil {
		return nil, err
	}

	err = h.filesCtrl.UploadFile(ctx, &site, file)
	if err != nil {
		return nil, err
	}

	return &UploadFileOutput{}, nil
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

	huma.Register(api, huma.Operation{
		OperationID:   "uploadFile",
		DefaultStatus: 200,
		Method:        "POST",
		Path:          "/sites/{siteId}/files/upload",
		Summary:       "This endpoint uploads a new file to the site.",
		Description:   "Upload a new file to the site.",
		Tags:          []string{"Sites"},
		Parameters: []*huma.Param{
			{
				Name:        "siteId",
				In:          "path",
				Description: "The ID of the site.",
				Required:    true,
				Example:     "1bd3edd9-0310-4bca-b88a-2b16d8ca7116",
				Schema: &huma.Schema{
					Type: "string",
				},
			},
			{
				Name:        "fileName",
				In:          "query",
				Description: "The name of the file.",
				Required:    true,
				Schema: &huma.Schema{
					Type: "string",
				},
			},
		},
		RequestBody: &huma.RequestBody{
			Required: true,
			Content: map[string]*huma.MediaType{
				"application/octet-stream": {
					Schema: &huma.Schema{
						Type:   "string",
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
	}, h.UploadFile)
}
