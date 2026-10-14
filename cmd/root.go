package cmd

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/colorprofile"
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/x/term"
	"github.com/zeiss/builder/cmd/sites"
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/ui/defaults"
	exitbanner "github.com/zeiss/builder/internal/ui/exit"
	ui "github.com/zeiss/builder/internal/ui/models"
	"github.com/zeiss/builder/internal/versions"
	"github.com/zeiss/builder/internal/workspaces"

	"github.com/spf13/cobra"
)

// Flags holds the command line flags.
var Flags = config.NewDefaultFlags()

func init() {
	// Commands
	RootCmd.AddCommand(InitCmd)
	RootCmd.AddCommand(ResetCmd)
	RootCmd.AddCommand(CheckCmd)
	RootCmd.AddCommand(sites.SitesCmd)
	RootCmd.AddCommand(AccountCmd)

	// State
	RootCmd.PersistentFlags().StringVarP(&Flags.URL, "url", "u", Flags.URL, "URL")
	RootCmd.PersistentFlags().StringVarP(&Flags.Path, "path", "p", Flags.Path, "path")
	RootCmd.PersistentFlags().StringVarP(&Flags.Session, "session", "s", Flags.Session, "session")
	RootCmd.PersistentFlags().BoolVar(&Flags.Continue, "continue", Flags.Continue, "continue")
	RootCmd.PersistentFlags().StringVarP(&Flags.Config, "config", "c", Flags.Config, "config file")
	RootCmd.PersistentFlags().BoolVarP(&Flags.Verbose, "verbose", "v", Flags.Verbose, "verbose output")
	RootCmd.PersistentFlags().BoolVarP(&Flags.Force, "force", "f", Flags.Force, "force init")
	RootCmd.PersistentFlags().BoolVar(&Flags.Debug, "debug", Flags.Debug, "debug mode")
	RootCmd.PersistentFlags().BoolVarP(&Flags.Yolo, "yolo", "y", Flags.Yolo, "yolo mode")

	// Runtime
	RootCmd.PersistentFlags().StringSliceVar(&Flags.Plugins, "plugin", Flags.Plugins, "plugin")
	RootCmd.PersistentFlags().StringSliceVar(&Flags.Vars, "var", Flags.Vars, "variables")

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
	workspace, err := workspaces.SetupWithProgressBar(Flags)
	if err != nil {
		return err
	}

	inputFilter := ui.NewFilter()
	var env uv.Environ = os.Environ()
	defaultUI := defaults.DefaultUI(workspace)
	rootUI := ui.NewUI(defaultUI, Flags.Session, Flags.Continue)

	program := tea.NewProgram(
		rootUI,
		tea.WithEnvironment(env),
		tea.WithContext(cmd.Context()),
		tea.WithFilter(inputFilter.Filter),
	)

	if _, err := program.Run(); err != nil {
		return err
	}

	workspace.Cleanup()

	printExitBanner()

	return nil
}

func printExitBanner() {
	tw, _, _ := term.GetSize(os.Stdout.Fd())
	body := exitbanner.Render(versions.VersionString, tw)
	if body == "" {
		return
	}

	fmt.Fprintln(colorprofile.NewWriter(os.Stderr, os.Environ()), body)
}
