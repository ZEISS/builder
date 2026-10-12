package cmd

import (
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/ui/models"

	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"
)

// CheckCmd is the command to check the configuration.
var CheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check the configuration",
	Long: `
This command checks the configuration and displays any errors.
`,
	RunE: runCheckCmd,
}

func runCheckCmd(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()

	specPath, err := config.DefaultConfig.Abs()
	if err != nil {
		return err
	}

	app := models.NewCheck(ctx, specPath)
	_, err = tea.NewProgram(app, tea.WithContext(ctx)).Run()
	if err != nil {
		return err
	}

	return nil
}
