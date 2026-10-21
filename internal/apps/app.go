package apps

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/pkg/specs"
	"github.com/zeiss/builder/server/middlewares/discovery"
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
	// WellKnownConfig returns the well-known config for the app.
	WellKnownConfig() *discovery.WellKnownConfig
}

// AppImpl is the implementation of the App interface.
type AppImpl struct {
	path            string
	force           bool
	wellKnownConfig *discovery.WellKnownConfig
}

// NewApp returns a new App instance.
func NewApp(wellKnownConfig *discovery.WellKnownConfig, path string, force bool) *AppImpl {
	return &AppImpl{
		wellKnownConfig: wellKnownConfig,
		path:            path,
		force:           force,
	}
}

// Path returns the path to the builder executable.
func (a *AppImpl) Path() (string, error) {
	appDir, err := config.AppDirPath()
	if err != nil {
		return "", err
	}

	err = filex.MkdirAll(appDir, 0o755)
	if err != nil {
		return "", err
	}

	return appDir, nil
}

// Spec returns the app specification.
func (a *AppImpl) LoadSpec() (*specs.Spec, error) {
	absPath, err := config.DefaultFlags.Abs()
	if err != nil {
		return nil, err
	}

	spec, err := specs.Read(absPath)
	if err != nil {
		return nil, err
	}

	return spec, nil
}

// Force removes the app directory and recreates it.
func (a *AppImpl) Force() bool {
	return a.force
}

// SpecPath returns the path to the app specification.
func (a *AppImpl) SpecPath() (string, error) {
	absPath, err := config.DefaultFlags.Abs()
	if err != nil {
		return "", err
	}

	return absPath, nil
}

// ValidateSpec validates the app specification.
func (a *AppImpl) ValidateSpec() error {
	spec, err := a.LoadSpec()
	if err != nil {
		return err
	}

	err = validate.Struct(spec)
	if err != nil {
		return err
	}

	return nil
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

	discv := discovery.NewClient()

	wellknownConfig, err := discv.Discover(ctx, config.DefaultFlags.URL)
	if err != nil {
		return nil, err
	}

	app := NewApp(wellknownConfig, appDir, config.DefaultFlags.Force)

	return app, nil
}
