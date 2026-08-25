package auth

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/zeiss/builder/internal/config"
)

var AuthResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the config, credentials, and state of the builder",
	RunE:  runReset,
}

func runReset(_ *cobra.Command, _ []string) error {
	path, err := config.ExpandConfigPath(config.DefaultConfig.Store)
	if err != nil {
		return err
	}
	path = filepath.Dir(path)

	err = os.RemoveAll(path)
	if err != nil {
		return err
	}

	return nil
}
