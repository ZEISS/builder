package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeiss/builder/internal/ports"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/zeiss/builder/internal/models"
)

var (
	checkMark = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
	errorMark = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).SetString("✗")
)

type (
	deviceAuthBeginMsg struct {
		deviceAuth *models.DeviceAuth
	}

	deviceAuthFinishMsg struct {
		account *models.Account
	}

	authErrorMsg struct {
		error error
	}
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Quit   key.Binding
	Accept key.Binding
	Down   key.Binding
	Up     key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Accept, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Accept, k.Quit}, // second column
	}
}

var keys = keyMap{
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

type loginModel struct {
	keys        keyMap
	help        help.Model
	inputStyle  lipgloss.Style
	lastKey     string
	quitting    bool
	completing  bool
	url         string
	code        string
	err         error
	authCtrl    ports.DeviceAuthController
	accountCtrl ports.AccountController
	ctx         context.Context
}

func New(ctx context.Context, authCtrl ports.DeviceAuthController, accountCtrl ports.AccountController) loginModel {
	return loginModel{
		keys:        keys,
		help:        help.New(),
		ctx:         ctx,
		authCtrl:    authCtrl,
		accountCtrl: accountCtrl,
		inputStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
	}
}

func (m loginModel) Init() tea.Cmd {
	return tea.Batch(m.beginAuth())
}

func (m loginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can gracefully truncate
		// its view as needed.
		m.help.SetWidth(msg.Width)
		return m, nil

	case deviceAuthBeginMsg:
		m.completing = true
		m.url = msg.deviceAuth.VerificationURI
		m.code = msg.deviceAuth.UserCode
		return m, tea.Batch(m.completeAuth(msg.deviceAuth))

	case deviceAuthFinishMsg:
		m.completing = false
		m.quitting = true
		return m, tea.Sequence(
			tea.Printf("%s %s", checkMark, "device auth completed."),
			tea.Quit,
		)

	case authErrorMsg: // handle deployment error
		m.err = msg.error
		m.quitting = true
		return m, tea.Sequence(
			tea.Printf("%s %s", errorMark, m.err),
			tea.Quit,
		)

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.Accept):
			m.lastKey = "Accept"
			return m, nil
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *loginModel) beginAuth() tea.Cmd {
	return func() tea.Msg {
		deviceAuth, err := m.authCtrl.Begin(m.ctx)
		if err != nil {
			return authErrorMsg{error: err}
		}

		return deviceAuthBeginMsg{deviceAuth: deviceAuth}
	}
}

func (m *loginModel) completeAuth(deviceAuth *models.DeviceAuth) tea.Cmd {
	return func() tea.Msg {
		account, err := m.authCtrl.Finish(m.ctx, deviceAuth)
		if err != nil {
			return authErrorMsg{error: err}
		}

		return deviceAuthFinishMsg{account: account}
	}
}

func (m loginModel) View() tea.View {
	if m.quitting {
		return tea.NewView("")
	}

	prompt := ""
	if m.completing {
		prompt += fmt.Sprintf("To sign in, please visit the following %s and enter the code:\n\n%s", m.url, m.code)
	}

	helpView := m.help.View(m.keys)
	height := 8 - strings.Count(prompt, "\n") - strings.Count(helpView, "\n")

	return tea.NewView(prompt + strings.Repeat("\n", height) + helpView)
}
