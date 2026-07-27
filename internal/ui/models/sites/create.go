package sites

import (
	"context"
	"fmt"

	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/ports"
	"github.com/zeiss/pkg/utilx"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

type (
	createSiteMsg      struct{}
	createSiteErrorMsg struct{ err error }
)

type createSiteModel struct {
	cfg       config.Config
	ctx       context.Context
	err       error
	lastKey   string
	quitting  bool
	keys      keyMap
	sitesCtrl ports.SitesController
}

// NewCreateSite creates a new create site model.
func NewCreateSite(ctx context.Context, cfg config.Config, sitesCtrl ports.SitesController) *createSiteModel {
	return &createSiteModel{
		cfg:       cfg,
		sitesCtrl: sitesCtrl,
		ctx:       ctx,
	}
}

// Init initializes the deploy model.
func (m *createSiteModel) Init() tea.Cmd {
	return tea.Sequence(m.createSite())
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
	s := fmt.Sprintf("%s %s successfully created.\n", checkMark, m.cfg.Spec.Sites.Name)

	if utilx.NotNil(m.err) {
		s = fmt.Sprintf("%s %s\n", errorMark, m.err)
	}

	return tea.NewView(s)
}

// createSite creates a new site.
func (m *createSiteModel) createSite() tea.Cmd {
	return func() tea.Msg {
		site := &models.Site{Name: m.cfg.Spec.Sites.Name}
		err := m.sitesCtrl.Create(m.ctx, site)

		if utilx.NotNil(err) {
			return createSiteErrorMsg{err: err}
		}

		return createSiteMsg{}
	}
}
