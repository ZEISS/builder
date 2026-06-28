package ports

import (
	"context"

	"github.com/zeiss/builder/internal/models"
)

// DeviceAuthRepository is an interface that defines the methods for device authentication operations.
type DeviceAuthRepository interface {
	// Begin is a method that begins the device authentication process.
	Begin(ctx context.Context) (*models.DeviceAuth, error)
	// Finish is a method that finishes the device authentication process.
	// TODO: could be returned in a token structure
	Finish(ctx context.Context, deviceAuth *models.DeviceAuth) (*models.Account, error)
}

// DeviceAuthController is an interface that defines the methods for device authentication operations.
type DeviceAuthController interface {
	// Begin is a method that begins the device authentication process.
	Begin(ctx context.Context) (*models.DeviceAuth, error)
	// Finish is a method that finishes the device authentication process.
	Finish(ctx context.Context, deviceAuth *models.DeviceAuth) (*models.Account, error)
}
