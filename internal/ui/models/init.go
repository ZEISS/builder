package models

import (
	"context"

	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/pkg/specs"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var checkMark = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")

type (
	errorInitMsg   struct{ err error }
	successInitMsg struct{}
)

type initModel struct {
	cfg      config.Config
	quitting bool
	err      error
	ctx      context.Context
}

// NewInit returns a new init model.
func NewInit(ctx context.Context, cfg config.Config) initModel {
	return initModel{
		cfg: cfg,
		ctx: ctx,
	}
}

// Init initializes the init model.
func (i initModel) Init() tea.Cmd {
	return i.init()
}

// Update handles messages for the reset model.
func (i initModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return i, nil

	case successInitMsg: // handle successful reset
		i.quitting = true
		return i, tea.Sequence(
			tea.Printf("%s successfully initialized", checkMark),
			tea.Quit,
		)

	case errorInitMsg: // handle reset error
		i.err = msg.err
		i.quitting = true
		return i, tea.Sequence(
			tea.Printf("%s %s", errorMark, i.err),
			tea.Quit,
		)
	}

	return i, nil
}

// View returns the view for the reset model.
func (r initModel) View() tea.View {
	return tea.NewView("")
}

// init is a command that initializes a project.
func (i initModel) init() tea.Cmd {
	return func() tea.Msg {
		example, err := specs.Example()
		if err != nil {
			return errorInitMsg{err: err}
		}

		if err := specs.Write(example, i.cfg.File, i.cfg.Flags.Force); err != nil {
			return errorInitMsg{err: err}
		}

		return successInitMsg{}
	}
}
