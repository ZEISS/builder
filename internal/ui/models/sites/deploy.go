package sites

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeiss/builder/internal/ports"
	"github.com/zeiss/builder/internal/ui/cmds"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/progress"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/zeiss/pkg/conv"
	"github.com/zeiss/pkg/utilx"
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

const (
	padding  = 2
	maxWidth = 80
)

// deployKeyMap defines the keybindings for the deploy model.
type deployKeyMap struct {
	Quit   key.Binding
	Accept key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k deployKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Accept, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k deployKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Accept, k.Quit}, // second column
	}
}

var deployKeys = deployKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type deploySiteModel struct {
	path      string
	id        string
	name      string
	ignore    []string
	completed int
	ctx       context.Context
	err       error
	keys      deployKeyMap
	lastKey   string
	percent   float64
	progress  progress.Model
	quitting  bool
	filesCtrl ports.FilesController
	total     int
	height    int
	width     int
}

// NewDeploy creates a new deploy model.
func NewDeploy(ctx context.Context, name, id, path string, ignore []string, filesCtrl ports.FilesController) *deploySiteModel {
	return &deploySiteModel{
		keys:      deployKeys,
		ctx:       ctx,
		name:      name,
		path:      path,
		id:        id,
		ignore:    ignore,
		filesCtrl: filesCtrl,
		progress:  progress.New(progress.WithDefaultBlend()),
	}
}

// Init initializes the deploy model.
func (m *deploySiteModel) Init() tea.Cmd {
	return cmds.ScanDir(m.path, m.ignore)
}

// Update handles incoming messages and updates the model accordingly.
func (m *deploySiteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg: // resize the window and progress bar
		m.width, m.height = msg.Width, msg.Height

		m.progress.SetWidth(msg.Width - padding*2 - 4)
		if m.progress.Width() > maxWidth {
			m.progress.SetWidth(maxWidth)
		}

		return m, nil

	case cmds.ErrorMsg:
		m.err = msg.Err
		m.quitting = true
		return m, tea.Sequence(
			tea.Printf("%s %s", errorMark, m.err),
			tea.Quit,
		)

	case cmds.SitesDeployScanSuccessMsg:
		m.total = msg.Total
		return m, cmds.DeploySite(m.ctx, m.id, m.path, msg.Files, m.filesCtrl)

	case cmds.SitesDeployFileSuccessMsg:
		m.completed++
		m.percent = float64(m.completed) / float64(m.total)

		if m.completed == m.total {
			m.quitting = true
			return m, tea.Quit
		}

		return m, nil

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
func (m deploySiteModel) View() tea.View {
	var v tea.View
	v.WindowTitle = "builder "

	var s strings.Builder
	if utilx.NotNil(m.err) {
		fmt.Fprintf(&s, "%s %s\n", errorMark, m.err)
	}

	if !m.quitting {
		pad := strings.Repeat(" ", padding)
		v.Content = "\n" +
			pad + "Deploying.." + "\n\n" +
			pad + conv.String(m.total) + " files" + "\n\n" +
			pad + m.progress.ViewAs(m.percent) + "\n\n" +
			pad + helpStyle("Press q to quit")
	}

	if m.quitting {
		fmt.Fprintf(&s, "%s %d files uploaded\n", checkMark, m.total)
	}

	v.Content = s.String()
	return v
}
