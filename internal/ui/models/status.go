package models

import (
	"strings"
	"time"

	"charm.land/bubbles/v2/help"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/x/ansi"
	"github.com/zeiss/builder/internal/ui/defaults"
)

// InfoMsgType is the type of an info message.
type InfoMsgType int

const (
	InfoMsgTypeInfo InfoMsgType = iota
	InfoMsgTypeSuccess
	InfoMsgTypeWarn
	InfoMsgTypeError
	InfoMsgTypeUpdate
	InfoMsgTypePlan
	InfoMsgTypeYolo
)

// InfoMsg is a message to be displayed in the status bar.
type InfoMsg struct {
	Type InfoMsgType
	Msg  string
	TTL  time.Duration
}

type ClearStatusMsg struct{}

// IsEmpty checks if the [InfoMsg] is empty.
func (m InfoMsg) IsEmpty() bool {
	var zero InfoMsg
	return m == zero
}

// DefaultStatusTTL is the default time-to-live for status messages.
const DefaultStatusTTL = 5 * time.Second

// badgeLeftInset is the number of cells between the status bar's left edge
// and the mode badge.
const badgeLeftInset = 1

// Status is the status bar and help model.
type Status struct {
	ui       *defaults.UI
	hideHelp bool
	help     help.Model
	helpKm   help.KeyMap
	msg      InfoMsg
}

// NewStatus creates a new status bar and help model.
func NewStatus(ui *defaults.UI, km help.KeyMap) *Status {
	s := new(Status)
	s.ui = ui
	s.help = help.New()
	s.help.Styles = ui.Styles.Help
	s.helpKm = km

	return s
}

// Set sets the status info message.
func (s *Status) Set(msg InfoMsg) {
	s.msg = msg
}

// Clear clears the status info message.
func (s *Status) Clear() {
	s.msg = InfoMsg{}
}

// SetWidth sets the width of the status bar and help view.
func (s *Status) SetWidth(width int) {
	helpStyle := s.ui.Styles.Status.Help
	horizontalPadding := helpStyle.GetPaddingLeft() + helpStyle.GetPaddingRight()
	s.help.SetWidth(width - horizontalPadding)
}

// ShowingAll returns whether the full help view is shown.
func (s *Status) ShowingAll() bool {
	return s.help.ShowAll
}

// ToggleHelp toggles the full help view.
func (s *Status) ToggleHelp() {
	s.help.ShowAll = !s.help.ShowAll
}

// SetHideHelp sets whether the app is on the onboarding flow.
func (s *Status) SetHideHelp(hideHelp bool) {
	s.hideHelp = hideHelp
}

// Draw draws the status bar onto the screen.
func (s *Status) Draw(scr uv.Screen, area uv.Rectangle) {
	helpStyle := s.ui.Styles.Status.Help
	helpWidth := area.Dx() - helpStyle.GetPaddingLeft() - helpStyle.GetPaddingRight()
	s.help.SetWidth(max(0, helpWidth))
	helpView := helpStyle.Render(s.help.View(s.helpKm))
	uv.NewStyledString(helpView).Draw(scr, area)

	if s.msg.IsEmpty() {
		return
	}

	var indStyle lipgloss.Style
	var msgStyle lipgloss.Style

	// Mode banners show the same badge that sits next to the help hints, so
	// they honor the same left inset to keep the indicator from jumping when
	// the banner appears or expires.
	indInset := 0
	switch s.msg.Type {
	case InfoMsgTypeError:
		indStyle = s.ui.Styles.Status.ErrorIndicator
		msgStyle = s.ui.Styles.Status.ErrorMessage
	case InfoMsgTypeWarn:
		indStyle = s.ui.Styles.Status.WarnIndicator
		msgStyle = s.ui.Styles.Status.WarnMessage
	case InfoMsgTypeUpdate:
		indStyle = s.ui.Styles.Status.UpdateIndicator
		msgStyle = s.ui.Styles.Status.UpdateMessage
	case InfoMsgTypeInfo:
		indStyle = s.ui.Styles.Status.InfoIndicator
		msgStyle = s.ui.Styles.Status.InfoMessage
	case InfoMsgTypeSuccess:
		indStyle = s.ui.Styles.Status.SuccessIndicator
		msgStyle = s.ui.Styles.Status.SuccessMessage
	}

	ind := indStyle.String()
	indWidth := lipgloss.Width(ind)
	msgPad := msgStyle.GetPaddingLeft() + msgStyle.GetPaddingRight()
	avail := max(0, area.Dx()-indWidth-msgPad-indInset)

	msg := strings.Join(strings.Split(s.msg.Msg, "\n"), " ")
	msg = ansi.Truncate(msg, avail, "…")

	if w := lipgloss.Width(msg); w < avail {
		msg += strings.Repeat(" ", avail-w)
	}

	info := msgStyle.Render(msg)
	uv.NewStyledString(strings.Repeat(" ", indInset)+ind+info).Draw(scr, area)
}

// clearInfoMsgCmd returns a command that clears the info message after the
// given TTL.
func clearInfoMsgCmd(ttl time.Duration) tea.Cmd {
	return tea.Tick(ttl, func(time.Time) tea.Msg {
		return ClearStatusMsg{}
	})
}

// Info returns a command that displays an info message.
func Info(msg string) tea.Cmd {
	return func() tea.Msg {
		return InfoMsg{
			Type: InfoMsgTypeInfo,
			Msg:  msg,
		}
	}
}

// Error returns a command that displays an error message.
func Error(err error) tea.Cmd {
	return func() tea.Msg {
		return InfoMsg{
			Type: InfoMsgTypeError,
			Msg:  err.Error(),
		}
	}
}

// Warn returns a command that displays a warning message.
func Warn(msg string) tea.Cmd {
	return func() tea.Msg {
		return InfoMsg{
			Type: InfoMsgTypeWarn,
			Msg:  msg,
		}
	}
}

// Success returns a command that displays a success message.
func Success(msg string) tea.Cmd {
	return func() tea.Msg {
		return InfoMsg{
			Type: InfoMsgTypeSuccess,
			Msg:  msg,
		}
	}
}
