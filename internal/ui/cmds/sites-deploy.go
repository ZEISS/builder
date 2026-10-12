package cmds

import (
	"context"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/zeiss/builder/internal/ports"
	"github.com/zeiss/builder/pkg/utils"
)

// SitesDeployScanSuccess is sent when a site deployment scan is successful.
type SitesDeployScanSuccessMsg struct {
	Total int
	Files []string
}

// SitesDeployFileSuccessMsg is sent when a site deployment is successful.
type SitesDeployFileSuccessMsg struct{ File string }

// ScanDir scans the directory for files to deploy.
func ScanDir(path string, ignore []string) tea.Cmd {
	return func() tea.Msg {
		files := utils.ScanDir(path, ignore)
		return SitesDeployScanSuccessMsg{Total: len(files), Files: files}
	}
}

// DeploySite deploys a site by scanning the directory and uploading files.
func DeploySite(ctx context.Context, id, path string, files []string, filesCtrl ports.FilesController) tea.Cmd {
	cmds := make([]tea.Cmd, 0, len(files))
	for _, file := range files {
		file = strings.TrimPrefix(filepath.Clean(strings.TrimPrefix(file, path)), "/")
		cmds = append(cmds, func() tea.Msg {
			err := filesCtrl.UploadFile(ctx, id, path, file)
			if err != nil {
				return ErrorMsg{Err: err}
			}

			return SitesDeployFileSuccessMsg{File: file}
		})
	}

	return tea.Sequence(cmds...)
}
