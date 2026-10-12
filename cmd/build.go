package cmd

import (
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"
	"github.com/zeiss/builder/internal/ui/models/build"
)

// BuildCmd is the command to build the app.
var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the app",
	Long: `
This command gets you to the build mode. Allowing you to build apps.
`,
	RunE: runBuildCmd,
}

func runBuildCmd(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()

	// Clear the terminal screen
	os.Stdout.WriteString("\x1b[2J\x1b[3J\x1b[H")

	model := build.New(ctx)

	program := tea.NewProgram(model)
	_, err := program.Run()
	if err != nil {
		return err
	}

	return nil
}
