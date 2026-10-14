package models

import (
	"image"
	"slices"
	"strings"

	"github.com/zeiss/builder/internal/home"
	"github.com/zeiss/builder/internal/ui/defaults"
	"github.com/zeiss/builder/internal/ui/dialogs"
	"github.com/zeiss/builder/internal/ui/utils"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/ultraviolet/layout"
	"github.com/charmbracelet/ultraviolet/screen"
	"github.com/zeiss/pkg/utilx"
)

// UIState represents the state of the UI.
type UIState uint8

// Possible UIState values.
const (
	UIOnboarding UIState = iota
	UIInitialize
	UILanding
	UIChat
)

// UI represents the main user interface model.
type UI struct {
	defaults *defaults.UI

	sessionID       string
	sessionContinue bool

	// State holds the current UI state.
	state UIState
	// header holds the header model for the UI.
	header *Header
	// dialogs captures the dialog state for the UI.
	dialogs *dialogs.Overlay
	// width holds the width of the UI canvas.
	width int
	// height holds the height of the UI canvas.
	height int
	// layout holds the layout of the UI.
	layout uiLayout
	// keyMap holds the key map for the UI.
	keyMap KeyMap
}

// NewUI creates a new UI model with the given defaults.
func NewUI(defaults *defaults.UI, sessionID string, sessionContinue bool) *UI {
	dialogs := dialogs.NewOverlay()
	header := NewHeader(defaults)
	keyMap := DefaultKeyMap()

	return &UI{
		defaults:        defaults,
		header:          header,
		sessionID:       sessionID,
		sessionContinue: sessionContinue,
		dialogs:         dialogs,
		keyMap:          keyMap,
	}
}

// Init initializes the UI model.
func (m *UI) Init() tea.Cmd {
	return tea.Batch()
}

// Update handles incoming messages and updates the UI model.
func (m *UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if cmd := m.handleKeyPressMsg(msg); cmd != nil {
			cmds = append(cmds, cmd)
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.updateLayoutAndSize()
	}

	return m, tea.Batch(cmds...)
}

// View renders the UI model.
func (m *UI) View() tea.View {
	var v tea.View
	v.AltScreen = true
	v.WindowTitle = "Sparkles " + home.Short(m.defaults.Workspace.Cwd())

	canvas := uv.NewScreenBuffer(m.width, m.height)
	v.Cursor = m.Draw(canvas, canvas.Bounds())

	content := strings.ReplaceAll(canvas.Render(), "\r\n", "\n") // normalize newlines
	contentLines := strings.Split(content, "\n")
	for i, line := range contentLines {
		// Trim trailing spaces for concise rendering
		contentLines[i] = strings.TrimRight(line, " ")
	}

	content = strings.Join(contentLines, "\n")
	v.Content = content

	return v
}

// Draw renders the UI model to the canvas.
func (m *UI) Draw(scr uv.Screen, area uv.Rectangle) *tea.Cursor {
	layout := m.generateLayout(area.Dx(), area.Dy())

	m.layout = layout
	m.updateSize()

	// Clear the screen first
	screen.Clear(scr)

	switch m.state {
	default:
		m.drawHeader(scr, layout.header)
	}

	if m.dialogs.HasDialogs() {
		return m.dialogs.Draw(scr, scr.Bounds())
	}

	return nil
}

// handleKeyPressMsg handles key press messages.
func (m *UI) handleKeyPressMsg(msg tea.KeyPressMsg) tea.Cmd {
	var cmds []tea.Cmd

	if key.Matches(msg, m.keyMap.Commands) {
		if cmd := m.openCommandsDialog(); cmd != nil {
			cmds = append(cmds, cmd)
		}
		return tea.Batch(cmds...)
	}

	if key.Matches(msg, m.keyMap.Quit) && !m.dialogs.ContainsDialog(dialogs.QuitID) {
		if cmd := m.openQuitDialog(); cmd != nil {
			cmds = append(cmds, cmd)
		}
		return tea.Batch(cmds...)
	}

	// Route all messages to dialog if one is open.
	if m.dialogs.HasDialogs() {
		return m.handleDialogMsg(msg)
	}

	return tea.Batch(cmds...)
}

// handleDialogMsg handles dialog messages.
func (m *UI) handleDialogMsg(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	action := m.dialogs.Update(msg)

	if utilx.Empty(action) {
		return tea.Batch(cmds...)
	}

	switch action.(type) {
	case dialogs.ActionClose:
		cmds = append(cmds, m.closeDialog())
	case dialogs.ActionQuit: // Quit the application.
		cmds = append(cmds, tea.Quit)
	default:
		cmds = append(cmds, utils.Noop())
	}

	return tea.Batch(cmds...)
}

// closeDialog closes the dialog with the given ID.
func (m *UI) closeDialog() tea.Cmd {
	if !m.dialogs.HasDialogs() {
		return utils.Noop()
	}

	m.dialogs.CloseFrontDialog()

	return utils.Noop()
}

// openCommandsDialog opens the commands dialog.
func (m *UI) openCommandsDialog() tea.Cmd {
	if m.dialogs.ContainsDialog(dialogs.CommandsID) {
		m.dialogs.BringToFront(dialogs.CommandsID)
		return utils.Noop()
	}

	commands, err := dialogs.NewCommands(m.defaults, "", false, false, false)
	if err != nil {
		return utils.Noop() // TODO: handle the errors
	}

	m.dialogs.OpenDialog(commands)

	return commands.InitialCmd()
}

// openDialog opens a dialog with the given ID.
func (m *UI) openDialog(id string) tea.Cmd {
	if id == dialogs.QuitID {
		return m.openQuitDialog()
	}

	return utils.Noop()
}

// openQuitDialog opens the quit confirmation dialog.
func (m *UI) openQuitDialog() tea.Cmd {
	if m.dialogs.ContainsDialog(dialogs.QuitID) {
		m.dialogs.BringToFront(dialogs.QuitID)
		return utils.Noop()
	}

	quitDialog := dialogs.NewQuit(m.defaults)
	m.dialogs.OpenDialog(quitDialog)

	return utils.Noop()
}

// updateSize updates the sizes of UI components based on the current layout.
func (m *UI) updateSize() {
}

// updateLayoutAndSize updates the layout and sizes of UI components.
func (m *UI) updateLayoutAndSize() {
	m.layout = m.generateLayout(m.width, m.height)
	m.updateSize()
}

// generateLayout calculates the layout rectangles for all UI components based
// on the current UI state and terminal dimensions.
func (m *UI) generateLayout(w, h int) uiLayout {
	area := image.Rect(0, 0, w, h)

	helpHeight := 1
	const landingHeaderHeight = 4

	// Add app margins
	var appRect, helpRect image.Rectangle
	layout.Vertical(
		layout.Len(area.Dy()-helpHeight),
		layout.Fill(1),
	).Split(area).Assign(&appRect, &helpRect)
	appRect.Min.Y += 1
	appRect.Max.Y -= 1
	helpRect.Min.Y -= 1
	appRect.Min.X += 1
	appRect.Max.X -= 1

	if slices.Contains([]UIState{UIOnboarding, UIInitialize}, m.state) {
		// extra padding on left and right for these states
		appRect.Min.X += 1
		appRect.Max.X -= 1
	}

	uiLayout := uiLayout{
		area:   area,
		status: helpRect,
	}

	// Handle different app states
	switch m.state {
	case UIOnboarding, UIInitialize:
		// Layout
		//
		// header
		// ------
		// main
		// ------
		// help

		var headerRect, mainRect image.Rectangle
		layout.Vertical(
			layout.Len(landingHeaderHeight),
			layout.Fill(1),
		).Split(appRect).Assign(&headerRect, &mainRect)
		uiLayout.header = headerRect
		uiLayout.main = mainRect
	}

	return uiLayout
}

// drawHeader draws the header section of the UI.
func (m *UI) drawHeader(scr uv.Screen, area uv.Rectangle) {
	m.header.Draw(
		scr,
		area,
		false,
		false,
		area.Dx(),
	)
}

// uiLayout defines the positioning of UI elements.
type uiLayout struct {
	// area is the overall available area.
	area uv.Rectangle

	// header is the header shown in special cases
	// e.x when the sidebar is collapsed
	// or when in the landing page
	// or in init/config
	header uv.Rectangle

	// main is the area for the main pane. (e.x chat, configure, landing)
	main uv.Rectangle

	// pills is the area for the pills panel.
	pills uv.Rectangle

	// editor is the area for the editor pane.
	editor uv.Rectangle

	// sidebar is the area for the sidebar.
	sidebar uv.Rectangle

	// status is the area for the status view.
	status uv.Rectangle

	// session details is the area for the session details overlay in compact mode.
	sessionDetails uv.Rectangle
}
