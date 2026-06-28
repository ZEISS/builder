package ui

import (
	"context"

	tea "charm.land/bubbletea/v2"
	"github.com/zeiss/builder/internal/config"
)

// Model represents a common interface for UI components.
type Model[T any] interface {
	Update(msg tea.Msg) (T, tea.Cmd)
	View() string
}

// Application ...
type Application interface {
	// Context returns the context of the application.
	Context() context.Context
	// Config returns the configuration of the application.
	Config() *config.Config
	// QuweueUpdateDraw adds a function to the queue to be executed in the main thread.
	QueueUpdateDraw(f func())
	Stop()
	Draw()
}
