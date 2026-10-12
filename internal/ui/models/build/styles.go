package build

import "charm.land/lipgloss/v2"

var builderColor = lipgloss.Color("208")

var (
	// Accent for "this is yours / an action": textarea marker, confirmations,
	// turn spinner, selected popover row.
	styleBuilder = lipgloss.NewStyle().Foreground(builderColor)

	// Neutrals for secondary copy (banner, status bar, summaries, popover descriptions).
	styleDim    = lipgloss.NewStyle().Faint(true)
	styleStatus = lipgloss.NewStyle().Faint(true)

	// Warn/error break the single accent rule on purpose: terminal convention
	// yellow and red, which fighting costs more than it gains.
	styleWarn  = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	styleError = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

	// Confirmations use the accent, "something good happened", never loud.
	styleOK = lipgloss.NewStyle().Foreground(builderColor)

	// Prompt marker and spinner share the accent: both mark live activity.
	stylePrompt  = lipgloss.NewStyle().Foreground(builderColor)
	styleSpinner = lipgloss.NewStyle().Foreground(builderColor)

	// User's echoed line: bold default, distinct from assistant markdown
	// without painting the user's text orange.
	styleUser = lipgloss.NewStyle().Bold(true)

	// Backend label: connected is quiet (bold, no colour); disconnected shouts
	// with yellow + `!` so the state survives colour stripped terminals.
	styleBackendOK   = lipgloss.NewStyle().Bold(true)
	styleBackendWarn = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("11"))

	// Popover: no backgrounds or marker column. Selected row is bold+accent;
	// the "current" entry (e.g. active profile) is bold, no colour.
	stylePopoverRow      = lipgloss.NewStyle()
	stylePopoverCurrent  = lipgloss.NewStyle().Bold(true)
	stylePopoverSelected = lipgloss.NewStyle().Bold(true).Foreground(builderColor)

	// Queued prompt panel: a faint rounded box above the prompt holding a prompt
	// the user lined up mid turn. Structural framing, not a highlight, so it
	// stays faint rather than taking the accent. Width is set per render.
	styleQueued = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1).Faint(true)
)
