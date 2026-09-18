package account

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeiss/builder/internal/ports"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/zeiss/builder/internal/models"
)

var (
	textStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render
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

// loginKeyMap defines the keybindings for the login model.
type loginKeyMap struct {
	Quit key.Binding
	Open key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k loginKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Open, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k loginKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Open, k.Quit}, // second column
	}
}

var loginKeys = loginKeyMap{
	Open: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "open browser"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type loginModel struct {
	keys        loginKeyMap
	help        help.Model
	inputStyle  lipgloss.Style
	lastKey     string
	quitting    bool
	completing  bool
	beginning   bool
	url         string
	code        string
	err         error
	spinner     spinner.Model
	authCtrl    ports.DeviceAuthController
	accountCtrl ports.AccountController
	ctx         context.Context
}

// New creates a new login model.
func New(ctx context.Context, authCtrl ports.DeviceAuthController, accountCtrl ports.AccountController) loginModel {
	model := loginModel{
		keys:        loginKeys,
		help:        help.New(),
		ctx:         ctx,
		authCtrl:    authCtrl,
		accountCtrl: accountCtrl,
		inputStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
	}

	model.resetSpinner()
	return model
}

// Init initializes the login model.
func (m loginModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.beginAuth())
}

func (m loginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can gracefully truncate
		// its view as needed.
		m.help.SetWidth(msg.Width)
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case deviceAuthBeginMsg:
		m.beginning = true
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
		case key.Matches(msg, m.keys.Open):
			m.lastKey = "Open"
			return m, nil
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m loginModel) View() tea.View {
	var s strings.Builder

	if m.quitting {
		return tea.NewView("")
	}

	if !m.beginning {
		fmt.Fprintf(&s, "\n %s %s\n\n", m.spinner.View(), textStyle("Setting up ..."))
		return tea.NewView(s.String())
	}

	helpView := m.help.View(m.keys)
	height := 4 - strings.Count(s.String(), "\n") - strings.Count(helpView, "\n")
	fmt.Fprintf(&s, "\nTo sign in, please visit the following %s and enter the code:\n\n%s", m.url, m.code)

	return tea.NewView(s.String() + strings.Repeat("\n", height) + helpView)
}

func (l *loginModel) beginAuth() tea.Cmd {
	return func() tea.Msg {
		deviceAuth, err := l.authCtrl.Begin(l.ctx)
		if err != nil {
			return authErrorMsg{error: err}
		}

		return deviceAuthBeginMsg{deviceAuth: deviceAuth}
	}
}

func (l *loginModel) completeAuth(deviceAuth *models.DeviceAuth) tea.Cmd {
	return func() tea.Msg {
		account, err := l.authCtrl.Finish(l.ctx, deviceAuth)
		if err != nil {
			return authErrorMsg{error: err}
		}

		return deviceAuthFinishMsg{account: account}
	}
}

func (l *loginModel) resetSpinner() {
	l.spinner = spinner.New()
	l.spinner.Spinner = spinner.Line
}
