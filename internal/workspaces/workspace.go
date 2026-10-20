package workspaces

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zeiss/builder/internal/adapters/db"
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/services"
	"github.com/zeiss/builder/internal/ui/utils"

	"github.com/charmbracelet/x/ansi"
	"github.com/google/uuid"
	"github.com/zeiss/pkg/dbx"
	"github.com/zeiss/pkg/filex"
	"github.com/zeiss/pkg/utilx"
)

// Workspace represents a workspace in the workspaces view.
type Workspace interface {
	// Debug returns whether the workspace is in debug mode.
	Debug() bool
	// Yolo returns whether the workspace is in yolo mode.
	Yolo() bool
	// Cwd returns the current working directory of the workspace.
	Cwd() string
	// Data returns the data dir of the workspace.
	DataDir() string
	// LogFile returns the log file of the workspace.
	LogFile() string
	// Cleanup cleans up the workspace.
	Cleanup() error
	// Initialized returns whether the workspace is initialized.
	Initialized() bool
	// Session returns the session controller of the workspace.
	Session() (*models.Session, bool)
}

var _ Workspace = (*WorkspaceImpl)(nil)

// WorkspaceImpl is the implementation of the Workspace interface.
type WorkspaceImpl struct {
	debug       bool
	yolo        bool
	cwd         string
	dir         string
	logFile     string
	sevices     services.Service
	funcs       []func() error
	initialized bool
	session     *models.Session
}

// New returns a new WorkspaceImpl.
func New(services services.Service, debug, yolo bool, cwd, dir string, funcs ...func() error) *WorkspaceImpl {
	f := make([]func() error, 0, len(funcs))
	for _, fn := range funcs {
		f = append(f, fn)
	}

	return &WorkspaceImpl{
		debug:   debug,
		yolo:    yolo,
		cwd:     cwd,
		dir:     dir,
		sevices: services,
		funcs:   f,
	}
}

// Debug returns whether the workspace is in debug mode.
func (w *WorkspaceImpl) Debug() bool {
	return w.debug
}

// Yolo returns whether the workspace is in yolo mode.
func (w *WorkspaceImpl) Yolo() bool {
	return w.yolo
}

// Cwd returns the current working directory of the workspace.
func (w *WorkspaceImpl) Cwd() string {
	return w.cwd
}

// DataDir returns the data dir of the workspace.
func (w *WorkspaceImpl) DataDir() string {
	return w.dir
}

// LogFile returns the log file of the workspace.
func (w *WorkspaceImpl) LogFile() string {
	return w.logFile
}

// Initialized returns whether the workspace is initialized.
func (w *WorkspaceImpl) Initialized() bool {
	return w.initialized
}

// Cleanup cleans up the workspace.
func (w *WorkspaceImpl) Cleanup() error {
	for _, f := range w.funcs {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}

// Session returns the session controller of the workspace.
func (w *WorkspaceImpl) Session() (*models.Session, bool) {
	return w.session, utilx.NotNil(w.session)
}

// SetupWithProgressBar sets up the progress bar for the workspaces view.
func SetupWithProgressBar(ctx context.Context, flags config.Flags) (Workspace, error) {
	showProgressBar := utils.SupportsProgressBar()

	if showProgressBar {
		_, _ = fmt.Fprintf(os.Stderr, ansi.SetIndeterminateProgressBar)
	}

	appDir, err := config.AppDirPath()
	if err != nil {
		return nil, err
	}

	database, err := db.NewDB(db.WithSqlite3(appDir))
	if err != nil {
		return nil, err
	}

	err = database.Migrate(ctx)
	if err != nil {
		return nil, err
	}

	store, err := dbx.NewDatabase(database.Conn(), db.NewReadTx(), db.NewWriteTx())
	if err != nil {
		return nil, err
	}

	funcs := []func() error{
		func() error {
			return store.Close()
		},
	}

	srv := services.NewService(store)

	ws, err := Setup(ctx, srv, flags, funcs...)
	if err != nil {
		return nil, err
	}

	if showProgressBar {
		_, _ = fmt.Fprintf(os.Stderr, ansi.ResetProgressBar)
	}

	return ws, nil
}

// Setup is setting up a workspace.
func Setup(ctx context.Context, services services.Service, flags config.Flags, funcs ...func() error) (Workspace, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current working directory: %w", err)
	}

	abs, err := filepath.Abs(flags.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to get path to workspace directory: %w", err)
	}

	if err := filex.MkdirAll(abs, 0o700); err != nil {
		return nil, fmt.Errorf("failed to create workspace directory: %w", err)
	}

	ws := New(services, flags.Debug, flags.Yolo, cwd, abs, funcs...)

	var session *models.Session
	if utilx.NotEmpty(flags.Session) && flags.Continue {
		prevSession := &models.Session{}
		id, err := uuid.Parse(flags.Session)
		if err != nil {
			return nil, fmt.Errorf("failed to generate session id: %w", err)
		}
		prevSession.ID = id

		err = services.GetSession(ctx, prevSession)
		if err != nil {
			return nil, fmt.Errorf("failed to get session: %w", err)
		}
		session = prevSession
	}

	if utilx.Empty(session) {
		newSession, err := services.CreateSession(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create session: %w", err)
		}
		session = newSession
	}

	ws.session = session

	logFile := filepath.Join(abs, "sparkles.log")
	ws.logFile = logFile

	return ws, nil
}
