package models

import (
	tea "charm.land/bubbletea/v2"
	"github.com/zeiss/builder/internal/ui/defaults"
)

// UI represents the main user interface model.
type UI struct {
	defaults *defaults.UI

	sessionID       string
	sessionContinue bool
}

// NewUI creates a new UI model with the given defaults.
func NewUI(defaults *defaults.UI, sessionID string, sessionContinue bool) *UI {
	return &UI{defaults: defaults, sessionID: sessionID, sessionContinue: sessionContinue}
}

// Init initializes the UI model.
func (m *UI) Init() tea.Cmd {
	return tea.Batch()
}

// Update handles incoming messages and updates the UI model.
func (m *UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:

	}
	return m, nil
}

// View renders the UI model.
func (m *UI) View() tea.View {
	return tea.NewView("")
}
