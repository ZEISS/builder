package dialogs

import tea "charm.land/bubbletea/v2"

// Action represents an action taken in a dialog after handling a message.
type Action any

// ActionClose is a message to close the current dialog.
type ActionClose struct{}

// ActionQuit is a message to quit the application.
type ActionQuit = tea.QuitMsg

// ActionOpenDialog is a message to open a dialog.
type ActionOpenDialog struct {
	DialogID string
}
