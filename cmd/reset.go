package cmd

import (
	"fmt"
	"os"

	"github.com/zeiss/builder/internal/config"

	"github.com/spf13/cobra"
	"github.com/zeiss/pkg/filex"
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
	appDir, err := config.AppDirPath()
	if err != nil {
		return err
	}

	ok, _ := filex.FileExists(appDir)
	if ok && !config.DefaultFlags.Force {
		return fmt.Errorf("use --force to delete the app folder: %s", appDir)
	}

	err = os.RemoveAll(appDir)
	if err != nil {
		return err
	}

	return nil
}
