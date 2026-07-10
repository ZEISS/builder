package sites

import (
	"context"
	"strings"

	"github.com/katallaxie/pkg/conv"
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/home"
	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/ports"
	"github.com/zeiss/builder/pkg/utils"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/progress"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

const (
	padding  = 2
	maxWidth = 80
)

type (
	progressDeployMsg struct{}
	deployErrorMsg    struct {
		err error
	}
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

type deployModel struct {
	cfg       *config.Config
	completed int
	ctx       context.Context
	err       error
	keys      keyMap
	lastKey   string
	percent   float64
	progress  progress.Model
	quitting  bool
	sitesCtrl ports.SitesController
	total     int
	height    int
	width     int
}

// NewDeploy creates a new deploy model.
func NewDeploy(ctx context.Context, cfg *config.Config, sitesCtrl ports.SitesController) *deployModel {
	return &deployModel{
		keys:      deployKeys,
		ctx:       ctx,
		cfg:       cfg,
		sitesCtrl: sitesCtrl,
		progress:  progress.New(progress.WithDefaultBlend()),
	}
}

// Init initializes the deploy model.
func (m *deployModel) Init() tea.Cmd {
	return tea.Sequence(m.writeFiles(), tea.Quit)
}

// Update handles incoming messages and updates the model accordingly.
func (m *deployModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg: // resize the window and progress bar
		m.width, m.height = msg.Width, msg.Height

		m.progress.SetWidth(msg.Width - padding*2 - 4)
		if m.progress.Width() > maxWidth {
			m.progress.SetWidth(maxWidth)
		}

		return m, nil

	case deployErrorMsg: // handle deployment error
		m.err = msg.err
		m.quitting = true

		return m, tea.Sequence(
			tea.Printf("%s %s", errorMark, m.err),
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
func (m deployModel) View() tea.View {
	var v tea.View
	v.WindowTitle = "builder " + home.Short(m.cfg.Spec.Name)

	pad := strings.Repeat(" ", padding)
	v.Content = "\n" +
		pad + "Deploying.." + "\n\n" +
		pad + conv.String(m.total) + " files" + "\n\n" +
		pad + m.progress.ViewAs(m.percent) + "\n\n" +
		pad + helpStyle("Press q to quit")

	return v
}

func (m *deployModel) writeFiles() tea.Cmd {
	files := utils.ScanDir(m.cfg.Spec.Sites.Path, m.cfg.Spec.Sites.Ignore)
	cmds := make([]tea.Cmd, 0, len(files))
	m.total = len(files)
	site := &models.Site{
		Name: m.cfg.Spec.Sites.Name,
	}

	for _, file := range files {
		cmds = append(cmds, func() tea.Msg {
			err := m.sitesCtrl.UploadFile(m.ctx, site, file)
			if err != nil {
				return deployErrorMsg{err: err}
			}

			return progressDeployMsg{}
		})
	}

	return tea.Sequence(cmds...)
}
