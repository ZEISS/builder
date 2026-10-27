package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zeiss/builder/internal/config"
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
	flags := config.DefaultFlags

	spec, err := flags.Spec()
	if err != nil {
		return err
	}

	if err := spec.Validate(); err != nil {
		return err
	}

	return nil
}
