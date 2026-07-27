package files

import (
	"context"
	"os"
	"path/filepath"

	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/server/configs"
	"github.com/zeiss/builder/server/ports"
	"github.com/zeiss/pkg/filex"
)

var _ ports.Files = (*Files)(nil)

// Files is the database adapter for the builder server.
type Files struct {
	cfg *configs.Config
}

// NewFiles creates a new Files instance.
func NewFiles(cfg *configs.Config) *Files {
	return &Files{cfg: cfg}
}

// UploadFile uploads a file to the filesystem.
func (f *Files) UploadFile(ctx context.Context, site *models.Site, file *models.File) error {
	path := filepath.Join(f.cfg.Flags.FilesFlags.Path, site.Name, file.Name)
	if err := filex.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	if err := os.WriteFile(path, file.Data, 0o644); err != nil {
		return err
	}

	return nil
}
