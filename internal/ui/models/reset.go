package models

import (
	"context"
	"fmt"
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/zeiss/builder/internal/ports"
	"github.com/zeiss/pkg/utilx"
)

var errorMark = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).SetString("✗")

type (
	errorResetMsg   struct{ err error }
	successResetMsg struct{}
)

// resetKeyMap defines a set of keybindings for the reset command.
type resetKeyMap struct {
	Quit    key.Binding
	Confirm key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view.
func (k resetKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Confirm, k.Quit}
}

// FullHelp returns keybindings for the expanded help view.
func (k resetKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Confirm, k.Quit}, // second column
	}
}

var resetKeys = resetKeyMap{
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "confirm"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type resetModel struct {
	accountCtrl ports.AccountController
	keys        resetKeyMap
	help        help.Model
	quitting    bool
	err         error
	ctx         context.Context
}

// NewReset returns a new reset model.
func NewReset(ctx context.Context, accountCtrl ports.AccountController) resetModel {
	return resetModel{
		accountCtrl: accountCtrl,
		keys:        resetKeys,
		help:        help.New(),
		ctx:         ctx,
	}
}

// Init initializes the reset model.
func (r resetModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the reset model.
func (r resetModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return r, nil

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, r.keys.Confirm):
			return r, r.reset()
		case key.Matches(msg, r.keys.Quit):
			r.quitting = true
			return r, tea.Quit
		}

	case successResetMsg: // handle successful reset
		r.quitting = true
		return r, tea.Quit

	case errorResetMsg: // handle reset error
		r.err = msg.err
		r.quitting = true
		return r, tea.Sequence(
			tea.Printf("%s %s", errorMark, r.err),
			tea.Quit,
		)
	}

	return r, nil
}

// View returns the view for the reset model.
func (r resetModel) View() tea.View {
	if r.quitting {
		return tea.NewView("")
	}

	var s strings.Builder
	fmt.Fprint(&s, errorMark)
	fmt.Fprint(&s, " This is a destructive action. Please confirm by pressing enter.")
	fmt.Fprintf(&s, "\n\n %s", r.help.View(r.keys))

	return tea.NewView(s.String())
}

// reset is a command that resets the store.
func (r resetModel) reset() tea.Cmd {
	return func() tea.Msg {
		err := r.accountCtrl.Reset(r.ctx)

		if utilx.NotNil(err) {
			return errorResetMsg{err: err}
		}

		return successResetMsg{}
	}
}
