package cmd

import (
	"fmt"

	"github.com/zeiss/builder/internal/apps"
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/pkg/specs"

	"github.com/spf13/cobra"
	"github.com/zeiss/pkg/filex"
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
	app, err := apps.SetupApp(cmd.Context())
	if err != nil {
		return err
	}

	specPath, err := app.SpecPath()
	if err != nil {
		return err
	}

	ok, _ := filex.FileExists(specPath)
	if ok && !config.DefaultFlags.Force {
		return fmt.Errorf("use --force to overwrite the config file: %s", specPath)
	}

	spec, err := specs.Example()
	if err != nil {
		return err
	}

	err = specs.Write(spec, specPath, config.DefaultFlags.Force)
	if err != nil {
		return err
	}

	return nil
}
