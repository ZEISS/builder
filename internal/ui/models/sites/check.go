package sites

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeiss/builder/internal/ui/cmds"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	textStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render
	checkMark = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
	errorMark = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).SetString("✗")
)

type checkSiteModel struct {
	ctx         context.Context
	spinner     spinner.Model
	name        string
	initialized bool
	quitting    bool
}

// NewCheckSite creates a new check site model.
func NewCheckSite(ctx context.Context) *checkSiteModel {
	model := &checkSiteModel{ctx: ctx}

	model.resetSpinner()

	return model
}

// Init initializes the deploy model.
func (m *checkSiteModel) Init() tea.Cmd {
	return cmds.Init(m.ctx)
}

// Update handles incoming messages and updates the model accordingly.
func (m *checkSiteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg: // resize the window and progress bar
		return m, nil

	case cmds.InitMsg:
		m.initialized = true
		m.name = msg.Config.Spec.Sites.Name
		return m, cmds.SitesCheckExists(m.ctx, msg.Client, m.name)

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case cmds.ErrorMsg: // handle deployment error
		m.quitting = true
		return m, tea.Sequence(
			tea.Printf("%s %s", errorMark, msg.Err.Error()),
			tea.Quit,
		)

	case cmds.SitesCheckExistsMsg:
		m.quitting = true
		return m, tea.Sequence(
			tea.Printf("%s %s %s", checkMark, m.name, "site exists."),
			tea.Quit,
		)

	case cmds.SitesCheckDoesNotExistMsg:
		m.quitting = true
		return m, tea.Sequence(
			tea.Printf("%s %s %s", errorMark, m.name, "site does not exist."),
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

	if !m.initialized {
		fmt.Fprintf(&s, "\n %s %s\n\n", m.spinner.View(), textStyle("Initializing ..."))
		return tea.NewView("")
	}

	fmt.Fprintf(&s, "\n %s %s\n\n", m.spinner.View(), textStyle("Checking site ..."))
	return tea.NewView(s.String())
}

func (m *checkSiteModel) resetSpinner() {
	m.spinner = spinner.New()
	m.spinner.Spinner = spinner.Line
}
