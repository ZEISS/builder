package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zeiss/builder/internal/apps"
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

	app, err := apps.SetupApp(ctx)
	if err != nil {
		return err
	}

	err = app.Validate()
	if err != nil {
		return err
	}

	return nil
}
