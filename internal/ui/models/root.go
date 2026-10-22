package models

import (
	"image"
	"slices"
	"strings"

	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/oidc"
	"github.com/zeiss/builder/internal/ui/defaults"
	"github.com/zeiss/builder/internal/ui/dialogs"
	"github.com/zeiss/builder/internal/ui/status"
	"github.com/zeiss/builder/internal/ui/utils"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/ultraviolet/layout"
	"github.com/charmbracelet/ultraviolet/screen"
	"github.com/zeiss/pkg/utilx"
)

// UIState represents the state of the UI.
type UIState uint8

// Possible UIState values.
const (
	UIInitializing = iota
	UIInitialized
)

// editorHeightMargin is the vertical margin added to the textarea height to
// account for the attachments row (top) and bottom margin.
const editorHeightMargin = 2

// TextareaMaxHeight is the maximum height of the prompt textarea.
const TextareaMaxHeight = 15

// TextareaMinHeight is the minimum height of the prompt textarea.
const TextareaMinHeight = 3

// UI represents the main user interface model.
type UI struct {
	root *defaults.UI

	sessionID       string
	sessionContinue bool

	// State holds the current UI state.
	state UIState
	// header holds the header model for the UI.
	header *Header
	// editor holds the editor model for the UI.
	textarea textarea.Model
	// status holds the status model for the UI.
	status *Status
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

	// Activity spinner
	activity   spinner.Model
	activityIs bool
}

// NewUI creates a new UI model with the given defaults.
func NewUI(root *defaults.UI, sessionID string, sessionContinue bool) *UI {
	ui := &UI{}
	ui.root = root
	ui.sessionID = sessionID
	ui.sessionContinue = sessionContinue

	ta := textarea.New()
	ta.SetStyles(ui.root.Styles.Editor.Textarea)
	ta.ShowLineNumbers = false
	ta.CharLimit = -1
	ta.SetVirtualCursor(false)
	ta.DynamicHeight = true
	ta.MinHeight = TextareaMinHeight
	ta.MaxHeight = TextareaMaxHeight
	// Keep "ctrl+a" for line-start (the textarea default); bind select-all
	// to "ctrl+shift+a" instead (line-start is also available via "home").
	ta.KeyMap.LineStart = key.NewBinding(
		key.WithKeys("home", "ctrl+a"),
		key.WithHelp("home", "line start"),
	)
	ta.KeyMap.SelectAll = key.NewBinding(
		key.WithKeys("ctrl+shift+a"),
		key.WithHelp("ctrl+shift+a", "select all"),
	)
	ta.KeyMap.CopySelection = key.NewBinding()
	ta.Focus()
	ta.Placeholder = defaults.ReadyPlaceholders.Random()
	ui.textarea = ta

	dialogs := dialogs.NewOverlay()
	ui.dialogs = dialogs

	header := NewHeader(root)
	ui.header = header

	keyMap := DefaultKeyMap()
	ui.keyMap = keyMap

	status := NewStatus(root, ui)
	ui.status = status

	return ui
}

// Init initializes the UI model.
func (m *UI) Init() tea.Cmd {
	return tea.Batch(m.initialize())
}

// Update handles incoming messages and updates the UI model.
func (m *UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.updateLayoutAndSize()
	case status.ClearStatusMsg:
		m.status.Clear()
	case status.InfoMsg:
		m.status.Set(msg)
		ttl := msg.TTL
		if ttl <= 0 {
			ttl = status.DefaultStatusTTL
		}
		cmds = append(cmds, status.ClearInfoMsgCmd(ttl))

	case dialogs.ActionCmd:
		if msg.Cmd != nil {
			cmds = append(cmds, msg.Cmd)
		}

	case spinner.TickMsg:
		if m.dialogs.HasDialogs() {
			if cmd := m.handleDialogMsg(msg); cmd != nil {
				cmds = append(cmds, cmd)
			}
		}

	case tea.KeyPressMsg:
		if cmd := m.handleKeyPressMsg(msg); cmd != nil {
			cmds = append(cmds, cmd)
		}

	default:
		if m.dialogs.HasDialogs() {
			if cmd := m.handleDialogMsg(msg); cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

// View renders the UI model.
func (m *UI) View() tea.View {
	var v tea.View
	v.AltScreen = true
	v.WindowTitle = "Sparkles " + m.root.Workspace.Cwd()

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
	case UIInitialized:
		m.drawHeader(scr, layout.header)
		main := uv.NewStyledString(m.landingView())
		main.Draw(scr, layout.main)

		editor := uv.NewStyledString(m.renderEditorView(scr.Bounds().Dx()))
		editor.Draw(scr, layout.editor)
	default:
		m.drawHeader(scr, layout.header)
	}

	m.status.Draw(scr, layout.status)

	if m.dialogs.HasDialogs() {
		return m.dialogs.Draw(scr, scr.Bounds())
	}

	return nil
}

// ShortHelp implements [help.KeyMap].
func (m *UI) ShortHelp() []key.Binding {
	k := &m.keyMap
	binds := []key.Binding{
		k.Quit,
		k.Commands,
		k.Help,
	}

	return binds
}

// FullHelp implements [help.KeyMap].
func (m *UI) FullHelp() [][]key.Binding {
	k := &m.keyMap
	help := k.Help
	help.SetHelp("ctrl+g", "less")

	binds := [][]key.Binding{
		{help, k.Quit, k.Commands, k.Help},
	}

	return binds
}

func (m *UI) initialize() tea.Cmd {
	var cmds []tea.Cmd

	if _, ok := m.root.Workspace.Session(); !ok {
		cmds = append(cmds, m.openConfigDialog())
	}

	return tea.Batch(cmds...)
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

	m.state = UIInitialized

	if utilx.Empty(action) {
		return tea.Batch(cmds...)
	}

	switch msg := action.(type) {
	// Handle dialog actions.
	case dialogs.ActionClose:
		cmds = append(cmds, m.closeDialog())

	case dialogs.ActionOpenDialog: // Open a dialog.
		m.dialogs.CloseDialog(dialogs.CommandsID)
		cmds = append(cmds, m.openDialog(msg.DialogID))

	case dialogs.ActionInitiateBuilder: // Initiate the builder.
		cmds = append(cmds, m.openOAuthBuilderDialog(msg))

	case dialogs.ActionQuit: // Quit the application.
		cmds = append(cmds, tea.Quit)

	case dialogs.ActionCmd: // Run a command.
		cmds = append(cmds, msg.Cmd)

	case dialogs.ActionToggleHelp:
		m.status.ToggleHelp()
		m.dialogs.CloseDialog(dialogs.CommandsID)

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

// openConfigDialog opens the config dialog.
func (m *UI) openConfigDialog() tea.Cmd {
	if m.dialogs.ContainsDialog(dialogs.ConfigID) {
		m.dialogs.BringToFront(dialogs.ConfigID)
		return utils.Noop()
	}

	config, err := dialogs.NewConfigDialog(m.root, "")
	if err != nil {
		return utils.Noop() // TODO: handle the errors
	}

	m.dialogs.OpenDialog(config)

	return utils.Noop()
}

// openOAuthBuilderDialog opens the OAuth builder dialog.
func (m *UI) openOAuthBuilderDialog(msg dialogs.ActionInitiateBuilder) tea.Cmd {
	if m.dialogs.ContainsDialog(dialogs.ConfigID) {
		m.dialogs.CloseDialog(dialogs.ConfigID)
	}

	oidcProvider := oidc.New(msg.URL, config.DefaultClientID)
	authDialog, cmd := dialogs.NewOAuth(m.root, oidcProvider, "Builder by Zeiss")

	if m.dialogs.ContainsDialog(authDialog.ID()) {
		m.dialogs.BringToFront(authDialog.ID())
		return nil
	}
	m.dialogs.OpenDialogWithGrace(authDialog)

	return tea.Batch(cmd)
}

// openCommandsDialog opens the commands dialog.
func (m *UI) openCommandsDialog() tea.Cmd {
	if m.dialogs.ContainsDialog(dialogs.CommandsID) {
		m.dialogs.BringToFront(dialogs.CommandsID)
		return utils.Noop()
	}

	commands, err := dialogs.NewCommands(m.root, "", false, false, false)
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

	if id == dialogs.ConfigID {
		return m.openConfigDialog()
	}

	if id == dialogs.DeployID {
		return m.openDeployDialog()
	}

	if id == dialogs.SessionsID {
		return m.openSessionsDialog()
	}

	return utils.Noop()
}

// openQuitDialog opens the quit confirmation dialog.
func (m *UI) openQuitDialog() tea.Cmd {
	if m.dialogs.ContainsDialog(dialogs.QuitID) {
		m.dialogs.BringToFront(dialogs.QuitID)
		return utils.Noop()
	}

	quitDialog := dialogs.NewQuit(m.root)
	m.dialogs.OpenDialog(quitDialog)

	return utils.Noop()
}

// openDeployDialog opens the deploy dialog.
func (m *UI) openDeployDialog() tea.Cmd {
	if m.dialogs.ContainsDialog(dialogs.DeployID) {
		m.dialogs.BringToFront(dialogs.DeployID)
		return utils.Noop()
	}

	deployDialog := dialogs.NewDeploy(m.root)
	m.dialogs.OpenDialog(deployDialog)

	return utils.Noop()
}

// openSessionsDialog opens the sessions dialog.
func (m *UI) openSessionsDialog() tea.Cmd {
	if m.dialogs.ContainsDialog(dialogs.SessionsID) {
		m.dialogs.BringToFront(dialogs.SessionsID)
		return utils.Noop()
	}

	sessionsDialog, err := dialogs.NewSessions(m.root, m.sessionID)
	if err != nil {
		return status.Error(err)
	}

	m.dialogs.OpenDialog(sessionsDialog)

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
	editorHeight := m.textarea.Height() + editorHeightMargin
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

	if slices.Contains([]UIState{UIInitializing}, m.state) {
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
	case UIInitializing:
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

	case UIInitialized:
		// Layout
		//
		// header
		// ------
		// main
		// ------
		// editor
		// ------
		// help
		var headerRect, mainRect image.Rectangle
		layout.Vertical(
			layout.Len(landingHeaderHeight),
			layout.Fill(1),
		).Split(appRect).Assign(&headerRect, &mainRect)
		var editorRect image.Rectangle
		layout.Vertical(
			layout.Len(mainRect.Dy()-editorHeight),
			layout.Fill(1),
		).Split(mainRect).Assign(&mainRect, &editorRect)
		// Remove extra padding from editor (but keep it for header and main)
		editorRect.Min.X -= 1
		editorRect.Max.X += 1
		uiLayout.header = headerRect
		uiLayout.main = mainRect
		uiLayout.editor = editorRect
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

// landingView renders the landing page.
func (m *UI) landingView() string {
	width := m.layout.main.Dx()

	content := lipgloss.JoinHorizontal(lipgloss.Left, "Info")

	return lipgloss.NewStyle().
		Width(width).
		Height(m.layout.main.Dy() - 1).
		PaddingTop(1).
		Render(
			lipgloss.JoinVertical(lipgloss.Left, "", content),
		)
}

// renderEditorView renders the editor view with attachments if any.
func (m *UI) renderEditorView(_ int) string {
	var attachmentsView string
	return strings.Join([]string{
		attachmentsView,
		m.textarea.View(),
		"", // margin at bottom of editor
	}, "\n")
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
