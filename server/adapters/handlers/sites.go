package handlers

import (
	"github.com/zeiss/builder/server/ports"
)

type sitesHandler struct {
	sitesCtrl       ports.Sites
	filesCtrl       ports.Files
	permissionsCtrl ports.Permissions
}

// NewSitesHandler creates a new SitesHandler with the given Sites controller.
func NewSitesHandler(sitesCtrl ports.Sites, filesCtrl ports.Files, permissionsCtrl ports.Permissions) *sitesHandler {
	return &sitesHandler{
		sitesCtrl:       sitesCtrl,
		filesCtrl:       filesCtrl,
		permissionsCtrl: permissionsCtrl,
	}
}
