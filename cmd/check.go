package cmd

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/pkg/specs"
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
	absPath, err := config.DefaultFlags.Abs()
	if err != nil {
		return err
	}

	spec, err := specs.Read(absPath)
	if err != nil {
		return err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(spec)
	if err != nil {
		return err
	}

	return nil
}
