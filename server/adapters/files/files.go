package files

import (
	"context"
	"path/filepath"

	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/server/ports"

	"github.com/gofiber/fiber/v3"
)

var _ ports.Files = (*Files)(nil)

// Files is a wrapper around the fiber storage to provide context aware storage.
type Files struct {
	storage fiber.Storage
}

// New creates a new Azure Blob Storage instance.
func New(storage fiber.Storage) *Files {
	return &Files{storage: storage}
}

// UploadFile uploads a file to the filesystem.
func (f *Files) UploadFile(ctx context.Context, site *models.Site, file *models.File) error {
	key := filepath.Join(site.Name, file.Name)

	err := f.storage.SetWithContext(ctx, key, file.Data, 0)
	if err != nil {
		return err
	}

	return nil
}
