package models

import (
	"context"

	"github.com/zeiss/builder/internal/ui/cmds"

	tea "charm.land/bubbletea/v2"
)

type checkModel struct {
	file     string
	quitting bool
	err      error
	ctx      context.Context
}

// NewCheck returns a new check model.
func NewCheck(ctx context.Context, file string) checkModel {
	return checkModel{
		file: file,
		ctx:  ctx,
	}
}

// Init initializes the check model.
func (m checkModel) Init() tea.Cmd {
	return cmds.Check(m.ctx, m.file)
}

// Update handles messages for the check model.
func (m checkModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil

	case cmds.CheckSuccessMsg:
		return m, tea.Sequence(
			tea.Printf("%s %s: %s", checkMark, m.file, "check passed"),
			tea.Quit,
		)

	case cmds.ErrorMsg: // handle reset error
		m.quitting = true
		return m, tea.Sequence(
			tea.Printf("%s %s: %s", errorMark, m.file, msg.Err.Error()),
			tea.Quit,
		)
	}

	return m, nil
}

// View returns the view for the check model.
func (m checkModel) View() tea.View {
	return tea.NewView("")
}
