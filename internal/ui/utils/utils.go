package utils

import (
	tea "charm.land/bubbletea/v2"
)

// Noop is a no-op command that does nothing.
func Noop() tea.Cmd {
	return nil
}
