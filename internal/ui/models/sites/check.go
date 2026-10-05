package sites

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/ports"
	"github.com/zeiss/pkg/utilx"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	textStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render
	checkMark = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
	errorMark = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).SetString("✗")
)

type (
	checkSiteExistsMsg    struct{}
	checkSiteNotExistsMsg struct{}
	siteCheckErrorMsg     struct{ err error }
)

type checkSiteModel struct {
	cfg       config.Config
	ctx       context.Context
	err       error
	spinner   spinner.Model
	quitting  bool
	sitesCtrl ports.SitesController
}

// NewCheckSite creates a new check site model.
func NewCheckSite(ctx context.Context, cfg config.Config, sitesCtrl ports.SitesController) *checkSiteModel {
	model := &checkSiteModel{
		ctx:       ctx,
		cfg:       cfg,
		sitesCtrl: sitesCtrl,
	}

	model.resetSpinner()

	return model
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

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case siteCheckErrorMsg: // handle deployment error
		m.err = msg.err
		m.quitting = true
		return m, tea.Sequence(
			tea.Printf("%s %s", errorMark, m.err),
			tea.Quit,
		)

	case checkSiteExistsMsg:
		m.quitting = true
		return m, tea.Sequence(
			tea.Printf("%s %s %s", checkMark, m.cfg.Spec.Sites.Name, "site exists."),
			tea.Quit,
		)

	case checkSiteNotExistsMsg:
		m.quitting = true
		return m, tea.Sequence(
			tea.Printf("%s %s %s", errorMark, m.cfg.Spec.Sites.Name, "site does not exist."),
			tea.Quit,
		)
	}

	return m, nil
}

// View renders the current state of the deploy model.
func (m checkSiteModel) View() tea.View {
	var s strings.Builder

	if m.quitting {
		return tea.NewView("")
	}

	fmt.Fprintf(&s, "\n %s %s\n\n", m.spinner.View(), textStyle("Checking site ..."))
	return tea.NewView(s.String())
}

func (m *checkSiteModel) getSite() tea.Cmd {
	return func() tea.Msg {
		site, err := m.sitesCtrl.GetSite(m.ctx, m.cfg.Spec.Sites.Name)

		if utilx.NotNil(err) {
			return siteCheckErrorMsg{err: err}
		}

		if utilx.Empty(site.ID) {
			return checkSiteNotExistsMsg{}
		}

		return checkSiteExistsMsg{}
	}
}

func (m *checkSiteModel) resetSpinner() {
	m.spinner = spinner.New()
	m.spinner.Spinner = spinner.Line
}
