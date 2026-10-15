package dialogs

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeiss/builder/internal/ui/defaults"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	uv "github.com/charmbracelet/ultraviolet"
)

type OAuthProvider interface {
	name() string
	initiateAuth() tea.Msg
	startPolling(deviceCode string, expiresIn int) tea.Cmd
	stopPolling() tea.Msg
	// supportsCodeEntry reports whether the sign-in can finish with a
	// code pasted back from the authorization page instead of the
	// loopback redirect — for when the user declines the browser
	// permission or the browser cannot reach the callback.
	supportsCodeEntry() bool
	// submitCode exchanges a pasted code (or full callback URL) for a
	// token. Only called when supportsCodeEntry is true.
	submitCode(input string) tea.Cmd
}

// OAuthState represents the current state of the device flow.
type OAuthState int

const (
	OAuthStateInitializing OAuthState = iota
	OAuthStateDisplay
	OAuthStateSuccess
	OAuthStateSaving
	OAuthStateError
)

// OAuthID is the identifier for the model selection dialog.
const OAuthID = "oauth"

// OAuth handles the OAuth flow authentication.
type OAuth struct {
	ui            *defaults.UI
	isOnboarding  bool
	oAuthProvider OAuthProvider

	State OAuthState

	spinner spinner.Model
	help    help.Model
	keyMap  struct {
		Copy      key.Binding
		CopyURL   key.Binding
		Submit    key.Binding
		FocusCode key.Binding
		Close     key.Binding
	}

	width           int
	deviceCode      string
	userCode        string
	verificationURL string
	expiresIn       int
	interval        int
	cancelFunc      context.CancelFunc
	// codeInput is the paste field for providers whose browser flow
	// can finish with a code entered back in the terminal.
	codeInput textinput.Model
}

var _ Dialog = (*OAuth)(nil)

// newOAuth creates a new device flow component.
func newOAuth(
	ui *defaults.UI,
	isOnboarding bool,
	oAuthProvider OAuthProvider,
) (*OAuth, tea.Cmd) {
	t := ui.Styles

	m := OAuth{}
	m.ui = ui
	m.isOnboarding = isOnboarding
	m.oAuthProvider = oAuthProvider
	m.width = 0 // Set dynamically in Draw().
	m.State = OAuthStateInitializing

	m.spinner = spinner.New(
		spinner.WithSpinner(spinner.Dot),
		spinner.WithStyle(t.Dialog.OAuth.Spinner),
	)

	m.help = help.New()
	m.help.Styles = t.DialogHelpStyles()

	m.keyMap.Copy = key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "copy code"),
	)
	m.keyMap.CopyURL = key.NewBinding(
		key.WithKeys("u"),
		key.WithHelp("u", "copy url"),
	)
	m.keyMap.Submit = key.NewBinding(
		key.WithKeys("enter", "ctrl+y"),
		key.WithHelp("enter", "copy & open"),
	)
	m.keyMap.FocusCode = key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "enter code"),
	)
	m.keyMap.Close = CloseKey

	// if m.oAuthProvider.supportsCodeEntry() {
	// 	m.codeInput = textinput.New()
	// 	m.codeInput.SetVirtualCursor(false)
	// 	m.codeInput.Prompt = "> "
	// 	m.codeInput.Placeholder = "Paste the code from the browser..."
	// 	m.codeInput.SetStyles(t.TextInput)
	// 	m.codeInput.Blur()
	// }

	return &m, tea.Batch(m.spinner.Tick)
}

// ID implements Dialog.
func (m *OAuth) ID() string {
	return OAuthID
}

// HandleMsg handles messages and state transitions.
func (m *OAuth) HandleMsg(msg tea.Msg) Action {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		switch m.State {
		case OAuthStateInitializing, OAuthStateDisplay, OAuthStateSaving:
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			if cmd != nil {
				return ActionCmd{cmd}
			}
		}
	}

	return nil
}

// View renders the device flow dialog.
func (m *OAuth) Draw(scr uv.Screen, area uv.Rectangle) *tea.Cursor {
	var (
		t           = m.ui.Styles
		dialogWidth = max(0, min(60, area.Dx()-t.Dialog.View.GetHorizontalBorderSize()))
		dialogStyle = t.Dialog.View.Width(dialogWidth)
	)
	m.width = dialogWidth
	if m.isOnboarding {
		view := m.dialogContent()
		cur := m.cursor()
		if cur != nil {
			cur = adjustOnboardingInputCursor(t, cur)
			DrawOnboardingCursor(scr, area, view, cur)
		} else {
			DrawOnboarding(scr, area, view)
		}
	} else {
		view := dialogStyle.Render(m.dialogContent())
		if cur := m.cursor(); cur != nil {
			DrawCenterCursor(scr, area, view, cur)
		} else {
			DrawCenter(scr, area, view)
		}
	}
	return nil
}

// cursor reports the paste field's cursor when it is focused, so the
// terminal caret sits in the input.
func (m *OAuth) cursor() *tea.Cursor {
	return nil
}

func (m *OAuth) dialogContent() string {
	t := m.ui.Styles

	switch m.State {
	case OAuthStateInitializing, OAuthStateSaving:
		return m.innerDialogContent()

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

func (m *OAuth) headerContent() string {
	var (
		t            = m.ui.Styles
		titleStyle   = t.Dialog.Title
		textStyle    = t.Dialog.PrimaryText
		dialogStyle  = t.Dialog.View.Width(m.width)
		headerOffset = titleStyle.GetHorizontalFrameSize() + dialogStyle.GetHorizontalFrameSize()
		dialogTitle  = fmt.Sprintf("Let’s authenticate with %s", m.oAuthProvider.name())
	)
	if m.isOnboarding {
		return textStyle.Render(dialogTitle)
	}
	return defaults.DialogTitle(t, titleStyle.Render(dialogTitle), m.width-headerOffset, t.Dialog.TitleGradFromColor, t.Dialog.TitleGradToColor)
}

func (m *OAuth) innerDialogContent() string {
	var (
		t                = m.ui.Styles
		instructionStyle = t.Dialog.OAuth.Instructions
		enterKeyStyle    = t.Dialog.OAuth.Enter
		successStyle     = t.Dialog.OAuth.Success
		linkStyle        = t.Dialog.OAuth.Link
		errorStyle       = t.Dialog.OAuth.ErrorText
		statusTextStyle  = t.Dialog.OAuth.StatusText
	)

	// innerWidth is the dialog's content area: total width minus the
	// View frame (border). Every block sizes to this so nothing gets
	// re-wrapped when the dialog frame renders it.
	innerWidth := m.width - t.Dialog.View.GetHorizontalFrameSize()

	switch m.State {
	case OAuthStateInitializing:
		return lipgloss.NewStyle().
			Width(innerWidth).
			Align(lipgloss.Center).
			Render(
				successStyle.Render(m.spinner.View()) +
					statusTextStyle.Render("Initializing..."),
			)

	case OAuthStateDisplay:
		// Render each text segment with its own style. Wrapping the
		// whole concatenation in a single style would lose the text
		// color after enterKeyStyle's reset code.
		instructionText := instructionStyle.Render("Press ") +
			enterKeyStyle.Render("enter") +
			instructionStyle.Render(" to open the browser and authenticate.")
		if m.userCode != "" {
			instructionText = instructionStyle.Render("Press ") +
				enterKeyStyle.Render("enter") +
				instructionStyle.Render(" to copy the code below and open the browser.")
		}
		instructions := lipgloss.NewStyle().
			Width(innerWidth).
			Padding(0, 1).
			Render(instructionText)

		elements := []string{"", instructions, ""}
		if m.userCode != "" {
			codeBox := lipgloss.NewStyle().
				Width(innerWidth).
				Height(7).
				Align(lipgloss.Center, lipgloss.Center).
				Background(t.Dialog.OAuth.UserCodeBg).
				Render(
					t.Dialog.OAuth.UserCode.Render(m.userCode),
				)
			elements = append(elements, codeBox, "")
		}

		link := linkStyle.Hyperlink(m.verificationURL, "id=oauth-verify").Render(m.verificationURL)
		url := statusTextStyle.
			Width(innerWidth).
			Padding(0, 1).
			Render("Browser not opening? Pay a visit to:\n" + link)

		waitingMsg := "Waiting for browser authentication..."
		if m.userCode != "" {
			waitingMsg = "Verifying..."
		}
		waiting := statusTextStyle.
			Width(innerWidth).
			Padding(0, 1).
			Render(
				successStyle.Render(m.spinner.View()) + statusTextStyle.Render(waitingMsg),
			)

		elements = append(elements, url, "", waiting, "")

		return lipgloss.JoinVertical(lipgloss.Left, elements...)

	case OAuthStateSuccess:
		return successStyle.
			Width(innerWidth).
			Padding(1).
			Render("Authentication successful!")

	case OAuthStateSaving:
		return lipgloss.NewStyle().
			Width(innerWidth).
			Align(lipgloss.Center).
			Render(
				successStyle.Render(m.spinner.View()) +
					statusTextStyle.Render(" Setting up Sparkles ..."),
			)

	case OAuthStateError:
		return errorStyle.
			Width(innerWidth).
			Padding(1).
			Render("Authentication failed.")

	default:
		return ""
	}
}

// FullHelp returns the full help view.
func (m *OAuth) FullHelp() [][]key.Binding {
	return [][]key.Binding{m.ShortHelp()}
}

// ShortHelp returns the full help view.
func (m *OAuth) ShortHelp() []key.Binding {
	switch m.State {
	case OAuthStateError:
		return []key.Binding{m.keyMap.Close}

	case OAuthStateSuccess:
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("enter", "ctrl+y", "esc"),
				key.WithHelp("enter", "finish"),
			),
		}

	case OAuthStateSaving:
		// No actionable keys while the save completes.
		return nil

	default:
		return nil
	}
}
