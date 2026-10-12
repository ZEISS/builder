package ports

import (
	"context"
)

// FilesRepository is an interface for the files repository.
type FilesRepository interface {
	// UploadFile uploads a file to the repository.
	UploadFile(ctx context.Context, name, cwd, file string) error
}

// FilesController is an interface for the files controller.
type FilesController interface {
	// UploadFile uploads a file to the repository.
	UploadFile(ctx context.Context, name, cwd, file string) error
}
