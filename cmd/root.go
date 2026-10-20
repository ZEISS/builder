package cmd

import (
	"fmt"
	"os"

	"github.com/zeiss/builder/cmd/sites"
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/ui/defaults"
	exitbanner "github.com/zeiss/builder/internal/ui/exit"
	ui "github.com/zeiss/builder/internal/ui/models"
	"github.com/zeiss/builder/internal/versions"
	"github.com/zeiss/builder/internal/workspaces"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/colorprofile"
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/x/term"
	"github.com/spf13/cobra"
)

func init() {
	// Commands
	RootCmd.AddCommand(InitCmd)
	RootCmd.AddCommand(ResetCmd)
	RootCmd.AddCommand(CheckCmd)
	RootCmd.AddCommand(sites.SitesCmd)
	RootCmd.AddCommand(AccountCmd)

	// State
	RootCmd.PersistentFlags().StringVarP(&config.DefaultFlags.URL, "url", "u", config.DefaultFlags.URL, "URL")
	RootCmd.PersistentFlags().StringVarP(&config.DefaultFlags.Path, "path", "p", config.DefaultFlags.Path, "path")
	RootCmd.PersistentFlags().StringVarP(&config.DefaultFlags.Session, "session", "s", config.DefaultFlags.Session, "session")
	RootCmd.PersistentFlags().BoolVar(&config.DefaultFlags.Continue, "continue", config.DefaultFlags.Continue, "continue")
	RootCmd.PersistentFlags().StringVarP(&config.DefaultFlags.Config, "config", "c", config.DefaultFlags.Config, "config file")
	RootCmd.PersistentFlags().BoolVarP(&config.DefaultFlags.Verbose, "verbose", "v", config.DefaultFlags.Verbose, "verbose output")
	RootCmd.PersistentFlags().BoolVarP(&config.DefaultFlags.Force, "force", "f", config.DefaultFlags.Force, "force init")
	RootCmd.PersistentFlags().BoolVar(&config.DefaultFlags.Debug, "debug", config.DefaultFlags.Debug, "debug mode")
	RootCmd.PersistentFlags().BoolVarP(&config.DefaultFlags.Yolo, "yolo", "y", config.DefaultFlags.Yolo, "yolo mode")

	// Runtime
	RootCmd.PersistentFlags().StringSliceVar(&config.DefaultFlags.Plugins, "plugin", config.DefaultFlags.Plugins, "plugin")
	RootCmd.PersistentFlags().StringSliceVar(&config.DefaultFlags.Vars, "var", config.DefaultFlags.Vars, "variables")

	// Output
	RootCmd.SilenceErrors = true
	RootCmd.SilenceUsage = true
}

var RootCmd = &cobra.Command{
	Use:   "sparkles",
	Short: "sparkles",
	RunE:  runRootCmd,
}

func runRootCmd(cmd *cobra.Command, args []string) error {
	workspace, err := workspaces.SetupWithProgressBar(cmd.Context(), config.DefaultFlags)
	if err != nil {
		return err
	}

	inputFilter := ui.NewFilter()
	var env uv.Environ = os.Environ()

	defaultUI := defaults.DefaultUI(workspace)
	rootUI := ui.NewUI(defaultUI, config.DefaultFlags.Session, config.DefaultFlags.Continue)

	program := tea.NewProgram(
		rootUI,
		tea.WithEnvironment(env),
		tea.WithContext(cmd.Context()),
		tea.WithFilter(inputFilter.Filter),
	)

	if _, err := program.Run(); err != nil {
		return err
	}

	if err := workspace.Cleanup(); err != nil {
		return err
	}

	session, ok := workspace.Session()
	if !ok {
		return nil
	}

	printExitBanner(session)

	return nil
}

func printExitBanner(session *models.Session) {
	tw, _, _ := term.GetSize(os.Stdout.Fd())
	body := exitbanner.Render(session, versions.VersionString, tw)
	if body == "" {
		return
	}

	fmt.Fprintln(colorprofile.NewWriter(os.Stderr, os.Environ()), body)
}
