package sites

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/home"
	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/ports"
	"github.com/zeiss/builder/pkg/utils"

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

type (
	progressDeployMsg   struct{}
	siteDeployExistsMsg struct{ site models.Site }
	siteDeployErrorMsg  struct{ err error }
)

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

var deployKeys = keyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type deploySiteModel struct {
	cfg       config.Config
	completed int
	ctx       context.Context
	err       error
	keys      keyMap
	lastKey   string
	percent   float64
	progress  progress.Model
	quitting  bool
	sitesCtrl ports.SitesController
	filesCtrl ports.FilesController
	site      models.Site
	total     int
	height    int
	width     int
}

// NewDeploy creates a new deploy model.
func NewDeploy(ctx context.Context, cfg config.Config, sitesCtrl ports.SitesController, filesCtrl ports.FilesController) *deploySiteModel {
	return &deploySiteModel{
		keys:      deployKeys,
		ctx:       ctx,
		cfg:       cfg,
		sitesCtrl: sitesCtrl,
		filesCtrl: filesCtrl,
		progress:  progress.New(progress.WithDefaultBlend()),
	}
}

// Init initializes the deploy model.
func (m *deploySiteModel) Init() tea.Cmd {
	return tea.Batch(m.getSite())
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

	case siteDeployErrorMsg:
		m.err = msg.err
		m.quitting = true
		return m, tea.Sequence(
			tea.Printf("%s %s", errorMark, m.err),
			tea.Quit,
		)

	case siteDeployExistsMsg:
		m.site = msg.site
		return m, tea.Sequence(
			m.writeFiles(),
			tea.Quit,
		)

	case progressDeployMsg:
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
	v.WindowTitle = "builder " + home.Short(m.cfg.Spec.Name)

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

func (m *deploySiteModel) writeFiles() tea.Cmd {
	cwd, _ := m.cfg.Cwd()

	path := filepath.Join(cwd, m.cfg.Spec.Sites.Path)
	files := utils.ScanDir(path, m.cfg.Spec.Sites.Ignore)

	cmds := make([]tea.Cmd, 0, len(files))
	m.total = len(files)

	for _, file := range files {
		cmds = append(cmds, func() tea.Msg {
			err := m.filesCtrl.UploadFile(m.ctx, m.site, file)
			if err != nil {
				return siteDeployErrorMsg{err: err}
			}

			return progressDeployMsg{}
		})
	}

	return tea.Sequence(cmds...)
}

func (m *deploySiteModel) getSite() tea.Cmd {
	return func() tea.Msg {
		site, err := m.sitesCtrl.GetSite(m.ctx, m.cfg.Spec.Sites.Name)

		if utilx.NotNil(err) {
			return siteDeployErrorMsg{err: err}
		}

		return siteDeployExistsMsg{site}
	}
}
