package build

import (
	"context"
	"strings"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/glamour/v2"
)

const (
	defaultWidth = 80 // bootstrap width before the first WindowSizeMsg
	minViewport  = 5  // rows reserved above the prompt for streaming tokens
	popoverCap   = 6  // max rows the popover may claim
)

// Model represents the build model.
type Model struct {
	// ctx is the context for the build model.
	ctx context.Context
	// cancel is the cancel function for the build context.
	cancel context.CancelFunc
	// renderer is the terminal renderer for the build output.
	renderer *glamour.TermRenderer
	// width is the width of the terminal window.
	width int
	// spinner is the spinner used to indicate build progress.
	spinner spinner.Model
}

// New returns a new build model.
func New(ctx context.Context) Model {
	r, _ := glamour.NewTermRenderer(glamour.WithStandardStyle("dark"), glamour.WithWordWrap(defaultWidth-4))

	sp := spinner.New()
	sp.Spinner = spinner.MiniDot
	sp.Style = styleSpinner

	pctx, cancel := context.WithCancel(ctx)

	m := Model{
		ctx:      pctx,
		cancel:   cancel,
		renderer: r,
		spinner:  sp,
	}

	return m
}

// Init initializes the model state.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles incoming messages and updates the model state.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

// View renders the model state to the terminal.
func (m Model) View() tea.View {
	var s strings.Builder

	return tea.NewView(s.String())
}
