package cmd

import (
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/ui/models"

	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"
)

// InitCmd is the command to initialize a new config.
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new config",
	Long: `
Initialize a new config (.builder.yml) in the current directory.
The config will be created with default values.
	`,
	RunE: runInit,
}

func runInit(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()

	absPath, err := config.DefaultConfig.Abs()
	if err != nil {
		return err
	}

	app := models.NewInit(ctx, absPath, config.DefaultConfig.Flags.Force)
	programm := tea.NewProgram(app, tea.WithContext(ctx))

	_, err = programm.Run()
	if err != nil {
		return err
	}

	return nil
}
