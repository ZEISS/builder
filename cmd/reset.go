package cmd

import (
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/controllers"
	"github.com/zeiss/builder/internal/ui/models"

	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"
)

// ResetCmd is the command to reset the config.
var ResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the configuration",
	Long: `
This command deletes the 'builder.db' file and resets the config to its default values.
`,
	RunE: runResetCmd,
}

func runResetCmd(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()

	accountCtrl := controllers.NewAccountController(config.DefaultConfig, nil)
	app := models.NewReset(ctx, accountCtrl)

	_, err := tea.NewProgram(app, tea.WithContext(ctx)).Run()
	if err != nil {
		return err
	}

	return err
}
