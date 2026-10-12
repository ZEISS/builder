package models

import (
	"context"
	"errors"

	"github.com/zeiss/builder/internal/ui/cmds"
	"github.com/zeiss/builder/pkg/specs"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var checkMark = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")

type initModel struct {
	file     string
	force    bool
	quitting bool
	ctx      context.Context
}

// NewInit returns a new init model.
func NewInit(ctx context.Context, file string, force bool) initModel {
	return initModel{
		file:  file,
		force: force,
		ctx:   ctx,
	}
}

// Init initializes the init model.
func (i initModel) Init() tea.Cmd {
	return cmds.Init(i.file, i.force)
}

// Update handles messages for the reset model.
func (m initModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil

	case cmds.InitSuccessMsg: // handle successful reset
		m.quitting = true
		return m, tea.Sequence(
			tea.Printf("%s %s successfully initialized", checkMark, m.file),
			tea.Quit,
		)

	case cmds.ErrorMsg: // handle reset error
		m.quitting = true

		if errors.Is(msg.Err, specs.ErrSpecExists) {
			return m, tea.Sequence(
				tea.Printf("%s %s: spec already exists, use --force to overwrite", errorMark, m.file),
				tea.Quit,
			)
		}

		return m, tea.Sequence(
			tea.Printf("%s %s: %s", errorMark, m.file, msg.Err.Error()),
			tea.Quit,
		)
	}

	return m, nil
}

// View returns the view for the reset model.
func (r initModel) View() tea.View {
	return tea.NewView("")
}
