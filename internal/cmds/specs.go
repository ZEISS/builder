package cmds

import (
	"path/filepath"

	tea "charm.land/bubbletea/v2"
	"github.com/zeiss/builder/internal/ui/status"
	"github.com/zeiss/builder/pkg/specs"
)

// InitializeSpec initializes a new specification in the current directory.
func InitializeSpec(cwd string, force bool) tea.Cmd {
	return func() tea.Msg {
		example, err := specs.Example()
		if err != nil {
			return status.Error(err)
		}

		path := filepath.Join(cwd, specs.DefaultName)

		if err := specs.Write(example, path, force); err != nil {
			return status.Error(err)
		}

		return status.Success("Specification initialized")
	}
}
