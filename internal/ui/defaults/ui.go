package defaults

import (
	"github.com/zeiss/builder/internal/ui/styles"
	"github.com/zeiss/builder/internal/workspaces"
)

// UI returns the default UI configurations.
type UI struct {
	// Styles holds the theme styles for the UI.
	Styles *styles.Styles
	// Workspace holds the workspace configuration for the UI.
	Workspace workspaces.Workspace
}

// DefaultUI returns the default UI configurations using the
// theme from config (or the default Charmtone theme if unset).
func DefaultUI(workspace workspaces.Workspace) *UI {
	s := styles.CharmtonePantera()

	return &UI{
		Workspace: workspace,
		Styles:    &s,
	}
}
