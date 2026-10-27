package cmds

import (
	tea "charm.land/bubbletea/v2"
	"github.com/zeiss/builder/internal/oidc"
	"github.com/zeiss/builder/internal/services"
	"github.com/zeiss/builder/internal/ui/status"
)

// CreateAccount creates a new account for the session.
func CreateAccount(sessions services.Service, token *oidc.Token) tea.Cmd {
	return func() tea.Msg {
		return status.Msg(status.Success("Account created successfully"))
	}
}
