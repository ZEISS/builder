package cmd

import (
	tea "charm.land/bubbletea/v2"
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/ui/models"

	"github.com/spf13/cobra"
)

// InitCmd is the command to initialize a new config.
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new config",
	RunE:  runInit,
}

func runInit(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	app := models.NewInit(ctx, config.DefaultConfig)

	_, err := tea.NewProgram(app, tea.WithContext(ctx)).Run()
	if err != nil {
		return err
	}

	return nil
}
