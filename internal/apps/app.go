package apps

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/zeiss/builder/internal/adapters/db"
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/services"
	"github.com/zeiss/builder/pkg/specs"
	"github.com/zeiss/builder/server/middlewares/discovery"
	"github.com/zeiss/pkg/dbx"
	"github.com/zeiss/pkg/filex"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

var (
	// ErrAppDirNotExists is returned when the app directory does not exist.
	ErrAppDirNotExists = errors.New("app directory does not exist")
)

// App is the interface for the builder application.
type App interface {
	// Path returns the path to the builder executable.
	Path() string
	// WellKnownConfig returns the well-known configuration for the app.
	WellKnownConfig() *discovery.WellKnownConfig
	// Validate validates the app specification.
	Validate() error
	// Spec returns the app specification.
	Spec() *specs.Spec
	// Account returns the current account.
	Account() *models.Account
}

var _ App = (*AppImpl)(nil)

// AppImpl is the implementation of the App interface.
type AppImpl struct {
	path            string
	force           bool
	wellKnownConfig *discovery.WellKnownConfig
	spec            *specs.Spec
	account         *models.Account
}

// NewApp returns a new App instance.
func NewApp(spec *specs.Spec, path string, force bool) *AppImpl {
	return &AppImpl{
		spec:  spec,
		path:  path,
		force: force,
	}
}

// Path returns the path to the builder executable.
func (a *AppImpl) Path() string {
	return a.path
}

// Force removes the app directory and recreates it.
func (a *AppImpl) Force() bool {
	return a.force
}

// Spec returns the app specification.
func (a *AppImpl) Spec() *specs.Spec {
	return a.spec
}

// ValidateSpec validates the app specification.
func (a *AppImpl) Validate() error {
	err := validate.Struct(a.spec)
	if err != nil {
		return err
	}

	return nil
}

// WellKnownConfig returns the well-known configuration for the app.
func (a *AppImpl) WellKnownConfig() *discovery.WellKnownConfig {
	return a.wellKnownConfig
}

// Account returns the current account.
func (a *AppImpl) Account() *models.Account {
	return a.account
}

// SetupApp sets up the app by discovering it and loading its specification.
func SetupApp(ctx context.Context) (*AppImpl, error) {
	appDir, err := config.AppDirPath()
	if err != nil {
		return nil, err
	}

	err = filex.MkdirAll(appDir, 0o755)
	if err != nil {
		return nil, err
	}

	absPath, err := config.DefaultFlags.Abs()
	if err != nil {
		return nil, err
	}

	spec, err := specs.Read(absPath)
	if err != nil {
		return nil, err
	}

	app := NewApp(spec, appDir, config.DefaultFlags.Force)

	return app, nil
}

// Initialize initializes the app.
func (a *AppImpl) Initialize(ctx context.Context) error {
	discv := discovery.NewClient()
	wellknownConfig, err := discv.Discover(ctx, config.DefaultFlags.URL)
	if err != nil {
		return err
	}
	a.wellKnownConfig = wellknownConfig

	path := a.Path()

	database, err := db.NewDB(db.WithSqlite3(path))
	if err != nil {
		return err
	}

	err = database.Migrate(ctx)
	if err != nil {
		return err
	}

	store, err := dbx.NewDatabase(database.Conn(), db.NewReadTx(), db.NewWriteTx())
	if err != nil {
		return err
	}
	srv := services.NewService(store)

	account := &models.Account{Current: true}
	err = srv.GetAccount(ctx, account)
	if err != nil {
		return err
	}

	return nil
}
