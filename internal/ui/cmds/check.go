package cmds

import (
	"context"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/zeiss/builder/pkg/specs"

	tea "charm.land/bubbletea/v2"
)

// CheckSuccessMsg is sent when the check command succeeds.
type CheckSuccessMsg struct{}

// Checks the spec and provides feedback
func Check(ctx context.Context, path string) tea.Cmd {
	return func() tea.Msg {
		f, err := os.ReadFile(path)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		spec := &specs.Spec{}
		if err := spec.UnmarshalYAML(f); err != nil {
			return ErrorMsg{Err: err}
		}

		validate := validator.New(validator.WithRequiredStructEnabled())
		err = validate.StructCtx(ctx, spec)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		return CheckSuccessMsg{}
	}
}
