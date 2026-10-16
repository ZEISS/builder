package models

import (
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/lipgloss/v2"
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/x/ansi"
	"github.com/zeiss/builder/internal/ui/defaults"
)

// badgeLeftInset is the number of cells between the status bar's left edge
// and the mode badge.
const badgeLeftInset = 1

// Status is the status bar and help model.
type Status struct {
	ui       *defaults.UI
	hideHelp bool
	help     help.Model
	helpKm   help.KeyMap
	msg      defaults.InfoMsg
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
func (s *Status) Set(msg defaults.InfoMsg) {
	s.msg = msg
}

// Clear clears the status info message.
func (s *Status) Clear() {
	s.msg = defaults.InfoMsg{}
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
	case defaults.InfoMsgTypeError:
		indStyle = s.ui.Styles.Status.ErrorIndicator
		msgStyle = s.ui.Styles.Status.ErrorMessage
	case defaults.InfoMsgTypeWarn:
		indStyle = s.ui.Styles.Status.WarnIndicator
		msgStyle = s.ui.Styles.Status.WarnMessage
	case defaults.InfoMsgTypeUpdate:
		indStyle = s.ui.Styles.Status.UpdateIndicator
		msgStyle = s.ui.Styles.Status.UpdateMessage
	case defaults.InfoMsgTypeInfo:
		indStyle = s.ui.Styles.Status.InfoIndicator
		msgStyle = s.ui.Styles.Status.InfoMessage
	case defaults.InfoMsgTypeSuccess:
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
