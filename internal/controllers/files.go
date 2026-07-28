package controllers

import (
	"context"

	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/ports"
)

var _ ports.FilesController = (*FilesController)(nil)

type FilesController struct {
	repo ports.FilesRepository
}

func NewFilesController(repo ports.FilesRepository) *FilesController {
	return &FilesController{
		repo: repo,
	}
}

// UploadFile is a method that uploads a file to the repository.
func (c *FilesController) UploadFile(ctx context.Context, site models.Site, cwd, file string) error {
	return c.repo.UploadFile(ctx, site, cwd, file)
}
