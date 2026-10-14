package defaults

import "github.com/zeiss/builder/internal/ui/styles"

// UI returns the default UI configurations.
type UI struct {
	// Styles holds the theme styles for the UI.
	Styles *styles.Styles
}

// DefaultUI returns the default UI configurations using the
// theme from config (or the default Charmtone theme if unset).
func DefaultUI() *UI {
	var s styles.Styles

	return &UI{
		Styles: &s,
	}
}
