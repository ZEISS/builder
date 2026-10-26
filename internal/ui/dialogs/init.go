package dialogs

import (
	"strings"

	"github.com/zeiss/builder/internal/ui/defaults"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	uv "github.com/charmbracelet/ultraviolet"
)

// ActionInitErrored is an action that is dispatched when the init flow encounters an error.
type ActionInitErrored struct{ Error error }

// ActionInitCompleted is an action that is dispatched when the init flow completes successfully.
type ActionInitCompleted struct{}

// InitState represents the current state of the init flow.
type InitState int

const (
	InitStateInitialized InitState = iota
	InitStateStarted
	InitStateFinished
	InitStateError
)

// InitID is the identifier for the init dialog.
const InitID = "init"

// InitDialog handles the initialization flow for the UI.
type InitDialog struct {
	ui *defaults.UI

	state   InitState
	spinner spinner.Model
	help    help.Model

	keyMap struct {
		Confirm key.Binding
		Close   key.Binding
	}

	width int
}

var _ Dialog = (*InitDialog)(nil)

// NewInit creates a new init dialog.
func NewInit(ui *defaults.UI) (*InitDialog, tea.Cmd) {
	t := ui.Styles

	m := &InitDialog{}
	m.ui = ui
	m.width = 0
	m.state = InitStateInitialized

	m.spinner = spinner.New(
		spinner.WithSpinner(spinner.Dot),
		spinner.WithStyle(t.Dialog.OAuth.Spinner),
	)

	m.help = help.New()
	m.help.Styles = t.DialogHelpStyles()
	m.keyMap.Confirm = key.NewBinding(
		key.WithKeys("enter", "ctrl+y"),
		key.WithHelp("enter", "confirm"),
	)
	m.keyMap.Close = CloseKey

	return m, tea.Batch(m.spinner.Tick)
}

// ID implements Dialog.
func (m *InitDialog) ID() string {
	return InitID
}

// HandleMsg handles messages and state transitions.
func (m *InitDialog) HandleMsg(msg tea.Msg) Action {
	switch msg := msg.(type) {
	case spinner.TickMsg: // spinner tick
		spinner, cmd := m.spinner.Update(msg)
		m.spinner = spinner
		return ActionCmd{cmd}

	case tea.KeyPressMsg:
		if key.Matches(msg, m.keyMap.Close) { // close dialog
			return ActionClose{}
		}
		if key.Matches(msg, m.keyMap.Confirm) { // confirm
			return ActionCmd{}
		}
	}

	return ActionNoop{}
}

// View renders the device flow dialog.
func (m *InitDialog) Draw(scr uv.Screen, area uv.Rectangle) *tea.Cursor {
	var (
		t           = m.ui.Styles
		dialogWidth = max(0, min(60, area.Dx()-t.Dialog.View.GetHorizontalBorderSize()))
		dialogStyle = t.Dialog.View.Width(m.width)
	)
	m.width = dialogWidth

	view := dialogStyle.Render(m.dialogContent())
	if cur := m.cursor(); cur != nil {
		DrawCenterCursor(scr, area, view, cur)
	} else {
		DrawCenter(scr, area, view)
	}

	return nil
}

// cursor reports the paste field's cursor when it is focused, so the
// terminal caret sits in the input.
func (m *InitDialog) cursor() *tea.Cursor {
	return nil
}

func (m *InitDialog) dialogContent() string {
	t := m.ui.Styles

	switch m.state {
	default:
		innerWidth := m.width - t.Dialog.View.GetHorizontalFrameSize()
		elements := []string{
			m.headerContent(),
			m.innerDialogContent(),
			renderDialogHelp(t, &m.help, m, innerWidth),
		}
		return strings.Join(elements, "\n")
	}
}

func (m *InitDialog) headerContent() string {
	var (
		t            = m.ui.Styles
		titleStyle   = t.Dialog.Title
		dialogStyle  = t.Dialog.View.Width(m.width)
		headerOffset = titleStyle.GetHorizontalFrameSize() + dialogStyle.GetHorizontalFrameSize()
		dialogTitle  = "Let's initialize ... "
	)

	return defaults.DialogTitle(t, titleStyle.Render(dialogTitle), m.width-headerOffset, t.Dialog.TitleGradFromColor, t.Dialog.TitleGradToColor)
}

func (m *InitDialog) innerDialogContent() string {
	var (
		t                = m.ui.Styles
		instructionStyle = t.Dialog.OAuth.Instructions
		enterKeyStyle    = t.Dialog.OAuth.Enter
		// successStyle     = t.Dialog.OAuth.Success
		// linkStyle        = t.Dialog.OAuth.Link
		// errorStyle       = t.Dialog.OAuth.ErrorText
		// statusTextStyle  = t.Dialog.OAuth.StatusText
	)

	// innerWidth is the dialog's content area: total width minus the
	// View frame (border). Every block sizes to this so nothing gets
	// re-wrapped when the dialog frame renders it.
	innerWidth := m.width - t.Dialog.View.GetHorizontalFrameSize()

	switch m.state {
	case InitStateInitialized:
		instructionText := instructionStyle.Render("Press ") +
			enterKeyStyle.Render("enter") +
			instructionStyle.Render(" to initialize your project.")

		instructions := lipgloss.NewStyle().
			Width(innerWidth).
			Padding(1, 1).
			Render(instructionText)

		elements := []string{
			instructions,
		}
		return lipgloss.JoinVertical(lipgloss.Left, elements...)
	default:
		return ""
	}
}

// FullHelp returns the full help view.
func (m *InitDialog) FullHelp() [][]key.Binding {
	return [][]key.Binding{m.ShortHelp()}
}

// ShortHelp returns the full help view.
func (m *InitDialog) ShortHelp() []key.Binding {
	switch m.state {
	default:
		submit := key.NewBinding(
			key.WithKeys("enter", "ctrl+y"),
			key.WithHelp("enter", "confirm"),
		)

		h := []key.Binding{
			submit,
			m.keyMap.Close,
		}

		return h
	}
}
