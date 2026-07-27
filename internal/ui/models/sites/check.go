package sites

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/ports"
	"github.com/zeiss/pkg/utilx"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	checkMark = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
	errorMark = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).SetString("✗")
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Quit   key.Binding
	Accept key.Binding
}

var checkSiteKeys = keyMap{
	Accept: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "open browser"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type (
	checkSiteMsg struct {
		siteExists bool
	}
	siteCheckErrorMsg struct {
		err error
	}
)

type checkSiteModel struct {
	cfg        config.Config
	ctx        context.Context
	err        error
	siteExists bool
	lastKey    string
	quitting   bool
	keys       keyMap
	sitesCtrl  ports.SitesController
}

// NewCheckSite creates a new check site model.
func NewCheckSite(ctx context.Context, cfg config.Config, sitesCtrl ports.SitesController) *checkSiteModel {
	return &checkSiteModel{
		ctx:       ctx,
		cfg:       cfg,
		sitesCtrl: sitesCtrl,
		keys:      checkSiteKeys,
	}
}

// Init initializes the deploy model.
func (m *checkSiteModel) Init() tea.Cmd {
	return tea.Sequence(m.getSite())
}

// Update handles incoming messages and updates the model accordingly.
func (m *checkSiteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg: // resize the window and progress bar
		return m, nil

	case siteCheckErrorMsg: // handle deployment error
		m.err = msg.err
		m.quitting = true

		return m, tea.Quit

	case checkSiteMsg:
		m.siteExists = msg.siteExists
		m.quitting = true

		return m, tea.Quit

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

// View renders the current state of the deploy model.
func (m checkSiteModel) View() tea.View {
	var s strings.Builder

	if utilx.NotNil(m.err) {
		fmt.Fprintf(&s, "%s %s\n", errorMark, m.err)
	}

	if m.siteExists {
		fmt.Fprintf(&s, "%s %s exists.\n", checkMark, m.cfg.Spec.Sites.Name)
	}

	if !m.siteExists {
		fmt.Fprintf(&s, "%s %s does not exist.\n", errorMark, m.cfg.Spec.Sites.Name)
	}

	return tea.NewView(s.String())
}

func (m *checkSiteModel) getSite() tea.Cmd {
	return func() tea.Msg {
		_, err := m.sitesCtrl.GetSite(m.ctx, m.cfg.Spec.Sites.Name)

		if utilx.NotNil(err) {
			return siteCheckErrorMsg{err: err}
		}

		return checkSiteMsg{siteExists: true}
	}
}
