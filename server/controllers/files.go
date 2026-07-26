package controllers

import (
	"context"

	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/server/ports"
)

var _ ports.Files = (*FilesController)(nil)

// FilesController is a controller for managing files.
type FilesController struct {
	db ports.Sites
	fs ports.Files
}

// NewFilesController creates a new files controller.
func NewFilesController(db ports.Sites) *FilesController {
	return &FilesController{db: db}
}

// UploadFile uploads a file to the site.
func (c *FilesController) UploadFile(ctx context.Context, site *models.Site, file *models.File) error {
	return c.fs.UploadFile(ctx, site, file)
}
