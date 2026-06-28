package ui

import tea "charm.land/bubbletea/v2"

// Notification represents a desktop notification request.
type Notification struct {
	Title   string
	Message string
}

// Backend defines the interface for sending desktop notifications.
type Backend interface {
	Send(n Notification) tea.Cmd
}
