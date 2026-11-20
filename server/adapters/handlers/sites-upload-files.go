package handlers

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/pkg/sbac"
)

// UploadFileInput is the input for the UploadFile operation.
type UploadFileInput struct {
	// Name is the name of the site.
	SiteName string `path:"name" json:"name" example:"fizzy-buzzy" doc:"The name of the site."`
	// Filename is the name of the file to upload.
	Filename string `query:"fileName" json:"fileName" example:"index.html" doc:"The filename of the file."`
	// RawBody is the raw body of the file to upload.
	RawBody []byte
}

// UploadFileOutput is the output of the UploadFile operation.
type UploadFileOutput struct {
	Body string `json:"body"`
}

// UploadFile handles the upload of a file to the site.
func (h *sitesHandler) UploadFile(ctx context.Context, input *UploadFileInput) (*UploadFileOutput, error) {
	err := h.permissionsCtrl.Allowed(ctx, input.SiteName, sbac.PermissionAdmin)
	if err != nil {
		return nil, huma.Error401Unauthorized("not authorized", err)
	}

	file := &models.File{
		Name: input.Filename,
		Data: input.RawBody,
	}

	site, err := h.sitesCtrl.GetSiteByName(ctx, input.SiteName)
	if err != nil {
		fmt.Printf("failed to get site: %v\n", err)
		return nil, err
	}

	err = h.filesCtrl.UploadFile(ctx, &site, file)
	if err != nil {
		fmt.Println("failed to upload file", err)
		return nil, err
	}

	fmt.Println(&site, file)

	return &UploadFileOutput{}, nil
}

// RegisterUploadFile registers the routes for the UploadFile operation.
func (h *sitesHandler) RegisterUploadFile(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID:   "uploadFile",
		DefaultStatus: 200,
		Method:        "POST",
		Path:          "/sites/{name}/files/upload",
		Summary:       "This endpoint uploads a new file to the site.",
		Description:   "Upload a new file to the site.",
		Tags:          []string{"Sites"},
		MaxBodyBytes:  10 * 1024 * 1024, // 10 MiB,
		Parameters: []*huma.Param{
			{
				Name:        "name",
				In:          "path",
				Description: "The name of the site.",
				Required:    true,
				Example:     "fizzy-buzzy",
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
			"200": {
				Description: "The file was uploaded successfully.",
			},
		},
	}, h.UploadFile)
}
