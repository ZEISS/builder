package store

import (
	"context"
	"os"
	"path/filepath"

	"github.com/katallaxie/pkg/filex"
	"github.com/zeiss/builder/server/ports"
)

var _ ports.Sites = (*fsImpl)(nil)

type fsImpl struct {
	root string
}

// NewFS returns a new fs adapter.
func NewFS(root string) *fsImpl {
	return &fsImpl{root: root}
}

// CreateSite implements ports.Sites.
func (f *fsImpl) CreateSite(_ context.Context, name string) error {
	path := filepath.Join(f.root, name)
	if err := os.MkdirAll(path, 0o755); err != nil {
		return err
	}

	return nil
}

// GetSite implements ports.Sites.
func (f *fsImpl) GetSite(_ context.Context, name string) (bool, error) {
	path := filepath.Join(f.root, name)
	return filex.FileExists(path)
}

// DeleteObject implements ports.Sites.
func (f *fsImpl) DeleteObject(_ context.Context, _, _ string) error {
	return nil
}

// PutObject implements ports.Sites.
func (f *fsImpl) PutObject(_ context.Context, site, name string, content []byte) error {
	path := filepath.Join(f.root, site, name)
	if err := filex.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	if err := os.WriteFile(path, content, 0o644); err != nil {
		return err
	}

	return nil
}
