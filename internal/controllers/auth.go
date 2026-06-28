package controllers

import (
	"context"

	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/ports"
)

var _ ports.DeviceAuthController = (*DeviceAuthController)(nil)

type DeviceAuthController struct {
	deviceAuthRepo ports.DeviceAuthRepository
	accountRepo    ports.AccountRepository
}

func NewDeviceAuthController(deviceAuthRepo ports.DeviceAuthRepository, accountRepo ports.AccountRepository) *DeviceAuthController {
	return &DeviceAuthController{
		deviceAuthRepo: deviceAuthRepo,
		accountRepo:    accountRepo,
	}
}

// Begin is a method that begins the device authentication process.
func (c *DeviceAuthController) Begin(ctx context.Context) (*models.DeviceAuth, error) {
	return c.deviceAuthRepo.Begin(ctx)
}

// Finish is a method that finishes the device authentication process.
func (c *DeviceAuthController) Finish(ctx context.Context, deviceAuth *models.DeviceAuth) (*models.Account, error) {
	account, err := c.deviceAuthRepo.Finish(ctx, deviceAuth)
	if err != nil {
		return nil, err
	}

	account.Current = true

	err = c.accountRepo.Create(ctx, account)
	if err != nil {
		return nil, err
	}

	return account, nil
}
