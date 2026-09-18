package account

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/ports"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type (
	fetchAccountsMsg struct {
		accounts []models.Account
	}
	switchErrorMsg struct {
		error error
	}
)

// switchKeyMap defines the keybindings for the login model.
type switchKeyMap struct {
	Quit   key.Binding
	Accept key.Binding
	Down   key.Binding
	Up     key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k switchKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Accept, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k switchKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Accept, k.Quit}, // second column
	}
}

var switchKeys = switchKeyMap{
	Accept: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "open browser"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
}

type switchModel struct {
	keys        switchKeyMap
	help        help.Model
	inputStyle  lipgloss.Style
	lastKey     string
	quitting    bool
	err         error
	accounts    []models.Account
	cursor      int
	accountCtrl ports.AccountController
	ctx         context.Context
}

func NewSwitch(ctx context.Context, accountCtrl ports.AccountController) switchModel {
	return switchModel{
		keys:        switchKeys,
		help:        help.New(),
		ctx:         ctx,
		accountCtrl: accountCtrl,
		inputStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
	}
}

func (m switchModel) Init() tea.Cmd {
	return tea.Batch(m.fetchAccounts())
}

func (m switchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can gracefully truncate
		// its view as needed.
		m.help.SetWidth(msg.Width)
		return m, nil

	case switchErrorMsg: // handle deployment error
		m.err = msg.error
		m.quitting = true
		return m, tea.Sequence(
			tea.Printf("%s %s", errorMark, m.err),
			tea.Quit,
		)

	case fetchAccountsMsg:
		m.accounts = msg.accounts
		return m, nil

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.Accept):
			m.lastKey = "Accept"
			return m, tea.Quit
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keys.Down):
			m.cursor++
			if m.cursor >= len(m.accounts) {
				m.cursor = 0
			}
		case key.Matches(msg, m.keys.Up):
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.accounts) - 1
			}
		}
	}

	return m, nil
}

func (m switchModel) View() tea.View {
	if m.quitting {
		return tea.NewView("")
	}

	s := strings.Builder{}
	s.WriteString("Switch your current account:\n\n")

	for i := range m.accounts {
		if m.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(m.accounts[i].Email)
		if m.accounts[i].Current {
			fmt.Fprintf(&s, " (%s)", checkMark)
		}

		s.WriteString("\n")
	}

	helpView := m.help.View(m.keys)

	return tea.NewView(s.String() + helpView)
}

func (m switchModel) fetchAccounts() tea.Cmd {
	return func() tea.Msg {
		var accounts []models.Account
		err := m.accountCtrl.List(m.ctx, &accounts)
		if err != nil {
			return switchErrorMsg{error: err}
		}

		return fetchAccountsMsg{accounts: accounts}
	}
}
