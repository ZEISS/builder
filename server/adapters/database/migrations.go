package database

import (
	"github.com/zeiss/builder/internal/models"

	fiber_goth "github.com/zeiss/fiber-goth/v3/adapters"
	"gorm.io/gorm"
)

// RunMigrations is a helper function to run the migrations for the database.
func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&fiber_goth.GothAccount{},
		&fiber_goth.GothSession{},
		&fiber_goth.GothUser{},
		&fiber_goth.GothCsrfToken{},
		&fiber_goth.GothVerificationToken{},
		&models.Deployment{},
		&models.Site{},
	)
}
