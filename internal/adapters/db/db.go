package db

import (
	"context"

	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/ports"
	"gorm.io/gorm"
)

// Database is a struct that represents the database adapter.
type Database struct {
	gorm *gorm.DB
}

var _ ports.AccountRepository = (*Database)(nil)

// RunMigrations is a helper function to run the migrations for the database.
func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Account{},
	)
}

// New creates a new instance of the Database adapter.
func New(gorm *gorm.DB) *Database {
	return &Database{
		gorm: gorm,
	}
}

// Get is a method that returns an account by ID.
func (d *Database) Get(ctx context.Context, account *models.Account) error {
	return d.gorm.WithContext(ctx).First(account).Error
}

// GetCurrent is a method that returns the current account.
func (d *Database) GetCurrent(ctx context.Context, account *models.Account) error {
	return d.gorm.WithContext(ctx).First(account).Where(&models.Account{Current: true}).Error
}

// Create is a method that creates a new account.
func (d *Database) Create(ctx context.Context, account *models.Account) error {
	return d.gorm.WithContext(ctx).FirstOrCreate(&account).Error
}

// Update is a method that updates an account.
func (d *Database) Update(ctx context.Context, account *models.Account) error {
	return d.gorm.WithContext(ctx).Save(account).Error
}

// Delete is a method that deletes an account.
func (d *Database) Delete(ctx context.Context, account *models.Account) error {
	return d.gorm.WithContext(ctx).Delete(account).Error
}

// List is a method that returns a list of all accounts.
func (d *Database) List(ctx context.Context, accounts *[]models.Account) error {
	return d.gorm.WithContext(ctx).Find(accounts).Error
}
