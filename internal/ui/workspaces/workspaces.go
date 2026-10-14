package workspaces

import (
	"fmt"
	"os"

	"github.com/charmbracelet/x/ansi"
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/ui/utils"
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
}

var _ Workspace = (*WorkspaceImpl)(nil)

// WorkspaceImpl is the implementation of the Workspace interface.
type WorkspaceImpl struct {
	debug bool
	yolo  bool
	cwd   string
	dir   string
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

// CleanupFunc is a function that cleans up a workspace.
type CleanupFunc func()

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

// SetupWithProgressBar sets up the progress bar for the workspaces view
func SetupWithProgressBar(flags config.Flags) (Workspace, CleanupFunc, error) {
	showProgressBar := utils.SupportsProgressBar()

	if showProgressBar {
		_, _ = fmt.Fprintf(os.Stderr, ansi.SetIndeterminateProgressBar)
	}

	ws, cleanup, err := Setup(flags)
	if err != nil {
		return nil, nil, err
	}

	if showProgressBar {
		_, _ = fmt.Fprintf(os.Stderr, ansi.ResetProgressBar)
	}

	return ws, cleanup, nil
}

// Setup is setting up a workspace.
func Setup(flags config.Flags) (Workspace, CleanupFunc, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get current working directory: %v", err)
	}

	ws := New(flags.Debug, flags.Yolo, cwd, "")
	cleanup := func() {}

	return ws, cleanup, nil
}
