package cmd

import (
	"fmt"

	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/pkg/specs"
	"github.com/zeiss/pkg/filex"

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
	absPath, err := config.DefaultFlags.Abs()
	if err != nil {
		return err
	}

	ok, _ := filex.FileExists(absPath)
	if ok && !config.DefaultFlags.Force {
		return fmt.Errorf("use --force to overwrite the config file: %s", absPath)
	}

	spec, err := specs.Example()
	if err != nil {
		return err
	}

	err = specs.Write(spec, absPath, config.DefaultFlags.Force)
	if err != nil {
		return err
	}

	return nil
}
