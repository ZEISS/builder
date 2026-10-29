package dialogs

import (
	"context"
	"fmt"

	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/ui/status"
	"github.com/zeiss/builder/internal/ui/styles"
	"github.com/zeiss/builder/internal/workspaces"
	"github.com/zeiss/pkg/utilx"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

type delegateKeyMap struct {
	choose key.Binding
	remove key.Binding
}

var _ list.Item = (*SessionItem)(nil)

// SessionItem represents a session item in the list.
type SessionItem struct {
	session *models.Session
}

// Title returns the title of the session item.
func (s SessionItem) Title() string { return fmt.Sprintf("%s (%s)", s.session.Name, s.session.ID) }

// Description returns the description of the session item.
func (s SessionItem) Description() string { return s.session.Description }

// FilterValue returns the filter value of the session item.
func (s SessionItem) FilterValue() string { return s.session.Name }

// NewSessionItemDelegate returns a new session item delegate.
func NewSessionItemDelegate(keys *delegateKeyMap, styles *styles.Styles, workspace workspaces.Workspace) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var title string

		session, ok := workspace.Session()
		if !ok {
			return nil
		}

		i, ok := m.SelectedItem().(SessionItem)
		if ok {
			title = i.Title()
		} else {
			return nil
		}

		if utilx.Equal(i.session.ID, session.ID) {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch {
			case key.Matches(msg, keys.choose):
				return m.NewStatusMessage(styles.ToolCallSuccess.Render("You chose " + title))

			case key.Matches(msg, keys.remove):
				if err := workspace.DeleteSession(context.TODO(), i.session); err != nil {
					return status.Error(err)
				}

				index := m.Index()
				m.RemoveItem(index)
				if len(m.Items()) == 0 {
					keys.remove.SetEnabled(false)
				}

				return status.Success(fmt.Sprintf("Deleted %s", title))
			}
		}

		return nil
	}

	help := []key.Binding{keys.choose, keys.remove}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

// Additional short help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.choose,
		d.remove,
	}
}

// Additional full help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.choose,
			d.remove,
		},
	}
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose"),
		),
		remove: key.NewBinding(
			key.WithKeys("x", "backspace"),
			key.WithHelp("x", "delete"),
		),
	}
}
