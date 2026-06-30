package sites

import (
	"context"
	"fmt"

	"github.com/katallaxie/pkg/utilx"
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/ports"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

var createSiteKeys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type (
	createSiteMsg      struct{}
	createSiteErrorMsg struct {
		err error
	}
)

type createSiteModel struct {
	cfg       *config.Config
	ctx       context.Context
	err       error
	lastKey   string
	quitting  bool
	keys      keyMap
	sitesCtrl ports.SitesController
}

// NewCreateSite creates a new create site model.
func NewCreateSite(ctx context.Context, cfg *config.Config, sitesCtrl ports.SitesController) *createSiteModel {
	return &createSiteModel{
		ctx:       ctx,
		cfg:       cfg,
		sitesCtrl: sitesCtrl,
		keys:      createSiteKeys,
	}
}

// Init initializes the deploy model.
func (m *createSiteModel) Init() tea.Cmd {
	return tea.Batch(m.createSite())
}

// Update handles incoming messages and updates the model accordingly.
func (m *createSiteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil

	case createSiteErrorMsg:
		m.err = msg.err
		m.quitting = true
		return m, tea.Quit

	case createSiteMsg:
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
func (m createSiteModel) View() tea.View {
	s := fmt.Sprintf("%s %s successfully created.\n", checkMark, m.cfg.Spec.Deploy.Site)

	if utilx.NotNil(m.err) {
		s = fmt.Sprintf("%s %s\n", errorMark, m.err)
	}

	return tea.NewView(s)
}

func (m *createSiteModel) createSite() tea.Cmd {
	return func() tea.Msg {
		site := &models.Site{Name: m.cfg.Spec.Deploy.Site}
		err := m.sitesCtrl.Create(m.ctx, site)
		if utilx.NotNil(err) {
			return createSiteErrorMsg{err: err}
		}

		return createSiteMsg{}
	}
}
