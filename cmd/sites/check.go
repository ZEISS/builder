package sites

import (
	"os"

	"github.com/zeiss/builder/internal/ui/models/sites"

	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"
)

var CheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Checks the status of a deployed site",
	RunE:  runCheck,
}

func runCheck(cmd *cobra.Command, _ []string) error {
	// clear all the stdout output
	os.Stdout.WriteString("\x1b[2J\x1b[3J\x1b[H")

	siteCheck := sites.NewCheckSite(cmd.Context())
	_, err := tea.NewProgram(siteCheck, tea.WithContext(cmd.Context())).Run()
	if err != nil {
		return err
	}

	return nil
}
