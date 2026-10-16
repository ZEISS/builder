package dialogs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/browser"
	"github.com/zeiss/builder/internal/clipboard"
	"github.com/zeiss/builder/internal/oidc"
	"github.com/zeiss/builder/internal/ui/defaults"
	"github.com/zeiss/pkg/utilx"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	uv "github.com/charmbracelet/ultraviolet"
)

// OAuthState represents the current state of the device flow.
type OAuthState int

const (
	OAuthStateInitializing OAuthState = iota
	OAuthStateAuthenticating
	OAuthStateSuccess
	OAuthStateSaving
	OAuthStateError
)

// OAuthID is the identifier for the model selection dialog.
const OAuthID = "oauth"

// OAuth handles the OAuth flow authentication.
type OAuth struct {
	ui           *defaults.UI
	isOnboarding bool
	oidcProvider oidc.OidcProvider
	name         string

	// State is the current state of the OAuth flow.
	State OAuthState

	spinner spinner.Model
	help    help.Model

	keyMap struct {
		Copy      key.Binding
		CopyURL   key.Binding
		Submit    key.Binding
		FocusCode key.Binding
		Close     key.Binding
	}

	width      int
	token      *oidc.Token
	deviceAuth *oidc.DeviceAuth

	cancelFunc context.CancelFunc
	codeInput  textinput.Model
}

var _ Dialog = (*OAuth)(nil)

// NewOAuth creates a new OAuth dialog.
func NewOAuth(ui *defaults.UI, oidcProvider oidc.OidcProvider, name string) (*OAuth, tea.Cmd) {
	t := ui.Styles

	m := &OAuth{}
	m.ui = ui
	m.oidcProvider = oidcProvider
	m.name = name
	m.width = 0
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

	return m, tea.Batch(m.spinner.Tick, Begin(m.oidcProvider))
}

// ID implements Dialog.
func (m *OAuth) ID() string {
	return OAuthID
}

// HandleMsg handles messages and state transitions.
func (m *OAuth) HandleMsg(msg tea.Msg) Action {
	switch msg := msg.(type) {

	case spinner.TickMsg: // spinner tick
		spinner, cmd := m.spinner.Update(msg)
		m.spinner = spinner
		return ActionCmd{tea.Batch(cmd, m.spinner.Tick)}

	case tea.KeyPressMsg:
		if key.Matches(msg, m.keyMap.Close) { // close dialog
			return ActionClose{}
		}

		if key.Matches(msg, m.keyMap.Submit) && m.State == OAuthStateAuthenticating {
			return ActionCmd{CopyCodeAndOpenURL(m.deviceAuth)}
		}

		if key.Matches(msg, m.keyMap.Submit) && m.State == OAuthStateSuccess {
			return ActionCmd{CloseAction()}
		}

	case ActionOAuthInitialized: // initiate oauth
		m.deviceAuth = msg.DeviceAuth
		m.State = OAuthStateAuthenticating
		return ActionCmd{Finish(m.oidcProvider, m.deviceAuth)}

	case ActionOAuthErrored: // oauth error
		m.State = OAuthStateError
		return ActionCmd{tea.Batch(defaults.Error(msg.Error))}

	case ActionOAuthCompleted: // oauth completed
		m.token = msg.Token
		m.State = OAuthStateSuccess
		return ActionCmd{defaults.Success("OAuth completed successfully")}

	default:
		return nil
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
		dialogTitle  = fmt.Sprintf("Let's authenticate with %s", m.name)
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

	case OAuthStateAuthenticating:
		// Render each text segment with its own style. Wrapping the
		// whole concatenation in a single style would lose the text
		// color after enterKeyStyle's reset code.
		instructionText := instructionStyle.Render("Press ") +
			enterKeyStyle.Render("enter") +
			instructionStyle.Render(" to open the browser and authenticate.")
		if m.deviceAuth.UserCode != "" {
			instructionText = instructionStyle.Render("Press ") +
				enterKeyStyle.Render("enter") +
				instructionStyle.Render(" to copy the code below and open the browser.")
		}
		instructions := lipgloss.NewStyle().
			Width(innerWidth).
			Padding(0, 1).
			Render(instructionText)

		elements := []string{"", instructions, ""}
		if m.deviceAuth.UserCode != "" {
			codeBox := lipgloss.NewStyle().
				Width(innerWidth).
				Height(7).
				Align(lipgloss.Center, lipgloss.Center).
				Background(t.Dialog.OAuth.UserCodeBg).
				Render(
					t.Dialog.OAuth.UserCode.Render(m.deviceAuth.UserCode),
				)
			elements = append(elements, codeBox, "")
		}

		link := linkStyle.Hyperlink(m.deviceAuth.VerificationURI, "id=oauth-verify").Render(m.deviceAuth.VerificationURI)
		url := statusTextStyle.
			Width(innerWidth).
			Padding(0, 1).
			Render("Browser not opening? Pay a visit to:\n" + link)

		waitingMsg := "Waiting for browser authentication..."
		if m.deviceAuth.UserCode != "" {
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
		submit := m.keyMap.Submit
		submit = key.NewBinding(
			key.WithKeys("enter", "ctrl+y"),
			key.WithHelp("enter", "open browser"),
		)

		h := []key.Binding{
			m.keyMap.CopyURL,
			submit,
			m.keyMap.Close,
		}

		return h
	}
}

// Begin initiates the OAuth device flow.
func Begin(oidcProvider oidc.OidcProvider) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		deviceAuth, err := oidcProvider.Begin(ctx)
		if err != nil {
			return ActionOAuthErrored{Error: fmt.Errorf("failed to initiate device auth: %w", err)}
		}

		return ActionOAuthInitialized{DeviceAuth: deviceAuth}
	}
}

// Finish completes the OAuth device flow.
func Finish(oidcProvider oidc.OidcProvider, deviceAuth *oidc.DeviceAuth) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		token, err := oidcProvider.Finish(ctx, deviceAuth)
		if err != nil {
			return ActionOAuthErrored{Error: fmt.Errorf("failed to finish device auth: %w", err)}
		}

		return ActionOAuthCompleted{Token: token}
	}
}

// CopyCodeAndOpenURL copies the code and opens the URL in the browser.
func CopyCodeAndOpenURL(deviceAuth *oidc.DeviceAuth) tea.Cmd {
	return func() tea.Msg {
		if utilx.Empty(deviceAuth.UserCode) {
			if err := browser.OpenURL(deviceAuth.VerificationURI); err != nil {
				return ActionOAuthErrored{fmt.Errorf("failed to open browser: %w", err)}
			}
			return nil
		}

		cmd := clipboard.CopyToClipboardWithCallback(
			deviceAuth.UserCode,
			"Code copied and URL opened",
			func() tea.Msg {
				if err := browser.OpenURL(deviceAuth.VerificationURI); err != nil {
					return ActionOAuthErrored{fmt.Errorf("failed to open browser: %w", err)}
				}
				return nil
			},
		)

		return ActionCmd{Cmd: cmd}
	}
}
