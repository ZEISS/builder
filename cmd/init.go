package cmd

import (
	"github.com/zeiss/builder/pkg/specs"

	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new config",
	RunE:  runInit,
}

func runInit(_ *cobra.Command, _ []string) error {
	example, err := specs.Example()
	if err != nil {
		return err
	}

	if err := specs.Write(example, cfg.File, cfg.Flags.Force); err != nil {
		return err
	}

	return nil
}
