package dialogs

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/zeiss/builder/internal/ui/defaults"
)

const (
	onboardingConnectPlaceholder = "https://builder.internal"
)

const ConfigID = "builder"

const defaultConfigDialogMaxWidth = 73

// ConfigDialog represents a builder configuration dialog.
type ConfigDialog struct {
	ui  *defaults.UI
	url string

	keyMap struct {
		Submit key.Binding
		Close  key.Binding
	}

	input textinput.Model
	help  help.Model
}

var _ Dialog = (*ConfigDialog)(nil)

// NewConfigDialog creates a new ConfigDialog.
func NewConfigDialog(ui *defaults.UI, url string) (*ConfigDialog, error) {
	t := ui.Styles
	m := &ConfigDialog{}
	m.ui = ui
	m.url = url

	help := help.New()
	help.Styles = t.DialogHelpStyles()

	m.help = help

	m.input = textinput.New()
	m.input.SetVirtualCursor(false)
	m.input.SetValue(m.url)
	m.input.Placeholder = onboardingConnectPlaceholder
	m.input.SetStyles(ui.Styles.TextInput)
	m.input.Focus()

	m.keyMap.Submit = key.NewBinding(
		key.WithKeys("enter", "ctrl+y"),
		key.WithHelp("enter", "confirm"),
	)
	m.keyMap.Close = CloseKey

	return m, nil
}

// ID implements Dialog.
func (m *ConfigDialog) ID() string {
	return ConfigID
}

// HandleMsg implements Dialog.
func (m *ConfigDialog) HandleMsg(msg tea.Msg) Action {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keyMap.Submit):
			return ActionInitiateBuilder{
				URL: m.input.Value(),
			}
		case key.Matches(msg, m.keyMap.Close):
			return ActionClose{}
		default:
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return ActionCmd{cmd}
		}
	}
	return nil
}

// Cursor returns the cursor for the dialog.
func (m *ConfigDialog) Cursor() *tea.Cursor {
	return InputCursor(m.ui.Styles, m.input.Cursor())
}

// Draw implements [Dialog].
func (m *ConfigDialog) Draw(scr uv.Screen, area uv.Rectangle) *tea.Cursor {
	t := m.ui.Styles
	width := max(0, min(defaultConfigDialogMaxWidth, area.Dx()-t.Dialog.View.GetHorizontalBorderSize()))
	innerWidth := width - t.Dialog.View.GetHorizontalFrameSize()
	m.input.SetWidth(dialogInputTextWidth(t, m.input, innerWidth))

	rc := NewRenderContext(t, width)
	rc.Title = "Connect to Builder"

	titleText := t.Dialog.PrimaryText.Render("To start, let's connect to your builder.")
	rc.AddPart(titleText)

	inputView := t.Dialog.InputPrompt.Render(m.input.View())
	rc.AddPart(inputView)

	rc.Help = renderDialogHelp(t, &m.help, m, innerWidth)

	cur := m.Cursor()

	rc.Title = ""
	rc.TitleInfo = ""
	rc.IsOnboarding = true
	view := rc.Render()
	cur = adjustOnboardingInputCursor(t, cur)
	DrawOnboardingCursor(scr, area, view, cur)

	return cur
}

// ShortHelp returns the short help view.
func (m *ConfigDialog) ShortHelp() []key.Binding {
	return []key.Binding{
		m.keyMap.Submit,
		m.keyMap.Close,
	}
}

// FullHelp returns the full help view.
func (m *ConfigDialog) FullHelp() [][]key.Binding {
	return [][]key.Binding{m.ShortHelp()}
}
