package dialogs

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/progress"
	tea "charm.land/bubbletea/v2"
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/zeiss/builder/internal/ui/defaults"
)

const DeployID = "deploy"

// DeployDialog represents a deploy configuration dialog.
type DeployDialog struct {
	ui  *defaults.UI
	url string

	keyMap struct {
		Submit key.Binding
		Close  key.Binding
	}

	help     help.Model
	progress progress.Model

	windowWidth int
}

var _ Dialog = (*DeployDialog)(nil)

// NewDeploy creates a new DeployDialog.
func NewDeploy(ui *defaults.UI) *DeployDialog {
	t := ui.Styles
	m := &DeployDialog{}
	m.ui = ui

	progress := progress.New(progress.WithDefaultBlend())
	m.progress = progress

	help := help.New()
	help.Styles = t.DialogHelpStyles()

	m.help = help

	m.keyMap.Submit = key.NewBinding(
		key.WithKeys("enter", "ctrl+y"),
		key.WithHelp("enter", "confirm"),
	)
	m.keyMap.Close = CloseKey

	return m
}

// ID implements Dialog.
func (m *DeployDialog) ID() string {
	return DeployID
}

// HandleMsg implements Dialog.
func (m *DeployDialog) HandleMsg(msg tea.Msg) Action {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keyMap.Close):
			return ActionClose{}
		default:
			var cmd tea.Cmd
			return ActionCmd{cmd}
		}
	}
	return nil
}

// Cursor returns the cursor for the dialog.
func (m *DeployDialog) Cursor() *tea.Cursor {
	return nil
}

// Draw implements [Dialog].
func (m *DeployDialog) Draw(scr uv.Screen, area uv.Rectangle) *tea.Cursor {
	t := m.ui.Styles
	width := max(0, min(defaultDialogMaxWidth, area.Dx()-t.Dialog.View.GetHorizontalBorderSize()))
	innerWidth := width - m.ui.Styles.Dialog.View.GetHorizontalFrameSize()

	rc := NewRenderContext(t, width)
	rc.Title = "Deploy to Builder"
	rc.Help = renderDialogHelp(t, &m.help, m, innerWidth)

	rc.AddPart(renderProgress(t, &m.progress, innerWidth))

	view := rc.Render()

	cur := m.Cursor()
	DrawCenterCursor(scr, area, view, cur)
	return cur
}

// ShortHelp returns the short help view.
func (m *DeployDialog) ShortHelp() []key.Binding {
	return []key.Binding{
		m.keyMap.Submit,
		m.keyMap.Close,
	}
}

// FullHelp returns the full help view.
func (m *DeployDialog) FullHelp() [][]key.Binding {
	return [][]key.Binding{m.ShortHelp()}
}
