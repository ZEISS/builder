package ui

import (
	"log/slog"

	tea "charm.land/bubbletea/v2"
)

// BellBackend is a notification backend that triggers the terminal bell character (\x07).
type BellBackend struct{}

// NewBellBackend creates a new bell notification backend.
func NewBellBackend() *BellBackend {
	return &BellBackend{}
}

// Send returns a [tea.Cmd] that triggers the terminal bell character (\x07).
// The terminal will emit an audible beep or visual flash based on user
// configuration. No message text is displayed.
func (b *BellBackend) Send(n Notification) tea.Cmd {
	slog.Debug("Sending bell notification", "title", n.Title, "message", n.Message)

	return tea.Raw("\x07")
}
