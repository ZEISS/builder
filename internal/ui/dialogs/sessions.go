package dialogs

import (
	"context"
	"image"
	"time"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/ui/defaults"
)

// ListItemGenerator is an interface for generating list items.
type ListItemGenerator[S ~[]R, R any] func(s S) []R

func sessionListGenerator(sessions ...*models.Session) []list.Item {
	items := make([]list.Item, 0, len(sessions))
	for _, s := range sessions {
		item := SessionItem{
			session: s,
		}
		items = append(items, item)
	}

	return items
}

// SessionsID is the identifier for the session selector dialog.
const SessionsID = "sessions"

const sessionDoubleClickThreshold = 400 * time.Millisecond

type sessionsMode uint8

const (
	sessionsModeNormal sessionsMode = iota
	sessionsModeDeleting
	sessionsModeUpdating
)

// Session is a session selector dialog.
type Session struct {
	ui *defaults.UI

	help               help.Model
	input              textinput.Model
	selectedSessionInx int
	sessions           []*models.Session

	sessionsMode  sessionsMode
	bodyArea      image.Rectangle
	mouseScrolled bool
	lastClickTime time.Time
	lastClickID   string

	list list.Model

	keyMap struct {
		Select        key.Binding
		Next          key.Binding
		Previous      key.Binding
		UpDown        key.Binding
		Delete        key.Binding
		Rename        key.Binding
		ConfirmRename key.Binding
		CancelRename  key.Binding
		ConfirmDelete key.Binding
		CancelDelete  key.Binding
		Close         key.Binding
	}

	delegateKeys *delegateKeyMap
}

var _ Dialog = (*Session)(nil)

// NewSessions creates a new Session dialog.
func NewSessions(ui *defaults.UI, selectedSessionID string) (*Session, error) {
	s := new(Session)
	s.sessionsMode = sessionsModeNormal
	s.ui = ui

	delegateKeys := newDelegateKeyMap()
	s.delegateKeys = delegateKeys

	sessions, err := ui.Workspace.ListSessions(context.Background())
	if err != nil {
		return nil, err
	}
	s.sessions = sessions

	help := help.New()
	help.Styles = ui.Styles.DialogHelpStyles()

	s.help = help

	delegate := NewSessionItemDelegate(delegateKeys, ui.Styles)
	items := sessionListGenerator(sessions...)

	sessionList := list.New(items, delegate, 0, 0)
	sessionList.SetShowTitle(false)
	sessionList.SetShowStatusBar(false)
	s.list = sessionList

	// s.list = list.NewFilterableList(sessionItems(ui.Styles, sessionsModeNormal, s.sessions...)...)
	// s.list.Focus()
	// s.list.SetSelected(s.selectedSessionInx)

	s.input = textinput.New()
	s.input.SetVirtualCursor(false)
	s.input.Placeholder = "Enter session name"
	s.input.SetStyles(ui.Styles.TextInput)
	s.input.Focus()

	s.keyMap.Select = key.NewBinding(
		key.WithKeys("enter", "tab", "ctrl+y"),
		key.WithHelp("enter", "choose"),
	)
	s.keyMap.Next = key.NewBinding(
		key.WithKeys("down", "ctrl+n"),
		key.WithHelp("↓", "next item"),
	)
	s.keyMap.Previous = key.NewBinding(
		key.WithKeys("up", "ctrl+p"),
		key.WithHelp("↑", "previous item"),
	)
	s.keyMap.UpDown = key.NewBinding(
		key.WithKeys("up", "down"),
		key.WithHelp("↑↓", "choose"),
	)
	s.keyMap.Delete = key.NewBinding(
		key.WithKeys("ctrl+x"),
		key.WithHelp("ctrl+x", "delete"),
	)
	s.keyMap.Rename = key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("ctrl+r", "rename"),
	)
	s.keyMap.ConfirmRename = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "confirm"),
	)
	s.keyMap.CancelRename = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	)
	s.keyMap.ConfirmDelete = key.NewBinding(
		key.WithKeys("y", "enter"),
		key.WithHelp("y", "delete"),
	)
	s.keyMap.CancelDelete = key.NewBinding(
		key.WithKeys("n", "esc"),
		key.WithHelp("n", "cancel"),
	)
	s.keyMap.Close = CloseKey

	return s, nil
}

// ID implements Dialog.
func (s *Session) ID() string {
	return SessionsID
}

// HandleMsg implements Dialog.
func (s *Session) HandleMsg(msg tea.Msg) Action {
	var cmds []tea.Cmd

	// This will also call our delegate's update function.
	newListModel, cmd := s.list.Update(msg)
	s.list = newListModel
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

// Draw implements [Dialog].
func (s *Session) Draw(scr uv.Screen, area uv.Rectangle) *tea.Cursor {
	t := s.ui.Styles

	width := max(0, min(defaultDialogMaxWidth, area.Dx()-t.Dialog.View.GetHorizontalBorderSize()))
	height := max(0, min(defaultDialogHeight, area.Dy()-t.Dialog.View.GetVerticalBorderSize()))
	innerWidth := width - t.Dialog.View.GetHorizontalFrameSize()
	innerHeight := height - t.Dialog.View.GetVerticalFrameSize()

	s.list.SetSize(innerWidth, innerHeight)

	var cur *tea.Cursor
	rc := NewRenderContext(t, width)
	rc.Title = "Sessions"

	bodyView := t.Dialog.List.Height(s.list.Height()).Render(s.list.View())
	bodyView = defaults.JoinScrollbar(t, bodyView, innerWidth, innerHeight, s.list.Height(), s.list.Index())
	rc.AddPart(bodyView)

	view := rc.Render()

	DrawCenterCursor(scr, area, view, cur)
	return cur
}

// ShortHelp implements [help.KeyMap].
func (s *Session) ShortHelp() []key.Binding {
	switch s.sessionsMode {
	case sessionsModeDeleting:
		return []key.Binding{
			s.keyMap.ConfirmDelete,
			s.keyMap.CancelDelete,
		}
	case sessionsModeUpdating:
		return []key.Binding{
			s.keyMap.ConfirmRename,
			s.keyMap.CancelRename,
		}
	default:
		return []key.Binding{
			s.keyMap.UpDown,
			s.keyMap.Rename,
			s.keyMap.Delete,
			s.keyMap.Select,
			s.keyMap.Close,
		}
	}
}

// FullHelp implements [help.KeyMap].
func (s *Session) FullHelp() [][]key.Binding {
	m := [][]key.Binding{}
	slice := []key.Binding{
		s.keyMap.UpDown,
		s.keyMap.Rename,
		s.keyMap.Delete,
		s.keyMap.Select,
		s.keyMap.Close,
	}

	switch s.sessionsMode {
	case sessionsModeDeleting:
		slice = []key.Binding{
			s.keyMap.ConfirmDelete,
			s.keyMap.CancelDelete,
		}
	case sessionsModeUpdating:
		slice = []key.Binding{
			s.keyMap.ConfirmRename,
			s.keyMap.CancelRename,
		}
	default:
	}
	for i := 0; i < len(slice); i += 4 {
		end := min(i+4, len(slice))
		m = append(m, slice[i:end])
	}
	return m
}
