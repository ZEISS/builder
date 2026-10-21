package cmd

import (
	"fmt"
	"os"

	"github.com/zeiss/builder/internal/apps"
	"github.com/zeiss/pkg/filex"

	"github.com/spf13/cobra"
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
	ctx := cmd.Context()
	app, err := apps.SetupApp(ctx)
	if err != nil {
		return err
	}

	path := app.Path()
	ok, err := filex.FileExists(path)
	if err != nil {
		return err
	}

	if ok && !app.Force() {
		return fmt.Errorf("Use --force to reset the configuration")
	}

	err = os.RemoveAll(app.Path())
	if err != nil {
		return err
	}

	return nil
}
