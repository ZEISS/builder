package auth

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

type switchModel struct {
	keys        keyMap
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
		keys:        keys,
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
