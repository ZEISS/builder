package cmds

import (
	tea "charm.land/bubbletea/v2"
	"github.com/zeiss/builder/pkg/specs"
)

// InitSuccessMsg is a message that indicates the initialization was successful.
type InitSuccessMsg struct{}

// Init initializes the application.
func Init(path string, force bool) tea.Cmd {
	return func() tea.Msg {
		example, err := specs.Example()
		if err != nil {
			return ErrorMsg{Err: err}
		}

		if err := specs.Write(example, path, force); err != nil {
			return ErrorMsg{Err: err}
		}

		return InitSuccessMsg{}
	}
}
