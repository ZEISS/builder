package workspaces

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/x/ansi"
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/ui/utils"
	"github.com/zeiss/pkg/filex"
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
	Cleanup()
}

var _ Workspace = (*WorkspaceImpl)(nil)

// WorkspaceImpl is the implementation of the Workspace interface.
type WorkspaceImpl struct {
	debug   bool
	yolo    bool
	cwd     string
	dir     string
	logFile string
}

// New returns a new WorkspaceImpl.
func New(debug, yolo bool, cwd, dir string) *WorkspaceImpl {
	return &WorkspaceImpl{
		debug: debug,
		yolo:  yolo,
		cwd:   cwd,
		dir:   dir,
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

// Cleanup cleans up the workspace.
func (w *WorkspaceImpl) Cleanup() {}

// SetupWithProgressBar sets up the progress bar for the workspaces view.
func SetupWithProgressBar(flags config.Flags) (Workspace, error) {
	showProgressBar := utils.SupportsProgressBar()

	if showProgressBar {
		_, _ = fmt.Fprintf(os.Stderr, ansi.SetIndeterminateProgressBar)
	}

	ws, err := Setup(flags)
	if err != nil {
		return nil, err
	}

	if showProgressBar {
		_, _ = fmt.Fprintf(os.Stderr, ansi.ResetProgressBar)
	}

	return ws, nil
}

// Setup is setting up a workspace.
func Setup(flags config.Flags) (Workspace, error) {
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

	ws := New(flags.Debug, flags.Yolo, cwd, "")

	logFile := filepath.Join(abs, "sparkles.log")
	ws.logFile = logFile

	return ws, nil
}
