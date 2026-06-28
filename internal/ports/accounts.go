package ports

import (
	"context"

	"github.com/zeiss/builder/internal/models"
)

// AccountRepository is an interface that defines the methods for account operations.
type AccountRepository interface {
	// GetCurrent is a method that returns the current account.
	GetCurrent(ctx context.Context, account *models.Account) error
	// Get is a method that returns an account by ID.
	Get(ctx context.Context, account *models.Account) error
	// Create is a method that creates a new account.
	Create(ctx context.Context, account *models.Account) error
	// Delete is a method that deletes an account.
	Delete(ctx context.Context, account *models.Account) error
	// Update is a method that updates an account.
	Update(ctx context.Context, account *models.Account) error
	// List is a method that returns a list of accounts.
	List(ctx context.Context, accounts *[]models.Account) error
}

// AccountController is an interface that defines the methods for account operations.
type AccountController interface {
	// GetCurrent is a method that returns the current account.
	GetCurrent(ctx context.Context, account *models.Account) error
	// Get is a method that returns an account by ID.
	Get(ctx context.Context, account *models.Account) error
	// Create is a method that creates a new account.
	Create(ctx context.Context, account *models.Account) error
	// Delete is a method that deletes an account.
	Delete(ctx context.Context, account *models.Account) error
	// Update is a method that updates an account.
	Update(ctx context.Context, account *models.Account) error
	// List is a method that returns a list of accounts.
	List(ctx context.Context, accounts *[]models.Account) error
}
