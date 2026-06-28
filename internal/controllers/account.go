package controllers

import (
	"context"

	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/ports"
)

var _ ports.AccountController = (*AccountController)(nil)

// AccountController is a controller for managing accounts.
type AccountController struct {
	account ports.AccountRepository
}

// NewAccountController creates a new AccountController.
func NewAccountController(account ports.AccountRepository) *AccountController {
	return &AccountController{
		account: account,
	}
}

// Get is a method that returns an account by ID.
func (c *AccountController) Get(ctx context.Context, account *models.Account) error {
	return c.account.Get(ctx, account)
}

// GetCurrent is a method that returns the current account.
func (c *AccountController) GetCurrent(ctx context.Context, account *models.Account) error {
	return c.account.GetCurrent(ctx, account)
}

// Create is a method that creates a new account.
func (c *AccountController) Create(ctx context.Context, account *models.Account) error {
	return c.account.Create(ctx, account)
}

// Update is a method that updates an existing account.
func (c *AccountController) Update(ctx context.Context, account *models.Account) error {
	return c.account.Update(ctx, account)
}

// Delete is a method that deletes an account.
func (c *AccountController) Delete(ctx context.Context, account *models.Account) error {
	return c.account.Delete(ctx, account)
}

// List is a method that returns a list of accounts.
func (c *AccountController) List(ctx context.Context, accounts *[]models.Account) error {
	return c.account.List(ctx, accounts)
}
