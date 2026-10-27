package dialogs

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"

	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/ui/styles"
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
func NewSessionItemDelegate(keys *delegateKeyMap, styles *styles.Styles) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var title string

		if i, ok := m.SelectedItem().(SessionItem); ok {
			title = i.Title()
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch {
			case key.Matches(msg, keys.choose):
				return m.NewStatusMessage(styles.ToolCallSuccess.Render("You chose " + title))

			case key.Matches(msg, keys.remove):
				index := m.Index()
				m.RemoveItem(index)
				if len(m.Items()) == 0 {
					keys.remove.SetEnabled(false)
				}
				return m.NewStatusMessage(styles.ToolCallSuccess.Render("Deleted " + title))
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
