package cmds

import (
	"context"
	"net/http"

	tea "charm.land/bubbletea/v2"
	"github.com/zeiss/builder/pkg/apis"
)

// SitesCheckExistsMsg is sent when a site exists.
type SitesCheckExistsMsg struct{}

// SitesCheckDoesNotExistMsg is sent when a site does not exist.
type SitesCheckDoesNotExistMsg struct{}

// SitesCheckExists checks if a site exists and returns a SiteExistsMsg if it does.
func SitesCheckExists(ctx context.Context, client *apis.ClientWithResponses, name string) tea.Cmd {
	return func() tea.Msg {
		params := &apis.GetSiteParams{Name: name}

		resp, err := client.GetSiteWithResponse(ctx, params)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		if resp.StatusCode() > http.StatusOK {
			return SitesCheckDoesNotExistMsg{}
		}

		return SitesCheckExistsMsg{}
	}
}
