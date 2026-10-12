package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/x/term"
	"github.com/zeiss/builder/cmd/sites"
	"github.com/zeiss/builder/internal/config"
	exitbanner "github.com/zeiss/builder/internal/ui/exit"
	"github.com/zeiss/builder/internal/versions"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(InitCmd)
	RootCmd.AddCommand(ResetCmd)
	RootCmd.AddCommand(CheckCmd)
	RootCmd.AddCommand(sites.SitesCmd)
	RootCmd.AddCommand(AccountCmd)

	RootCmd.PersistentFlags().StringVarP(&config.DefaultConfig.Flags.URL, "url", "u", config.DefaultConfig.Flags.URL, "URL")
	RootCmd.PersistentFlags().StringVarP(&config.DefaultConfig.File, "config", "c", config.DefaultConfig.File, "config file")
	RootCmd.PersistentFlags().BoolVarP(&config.DefaultConfig.Flags.Verbose, "verbose", "v", config.DefaultConfig.Flags.Verbose, "verbose output")
	RootCmd.PersistentFlags().BoolVarP(&config.DefaultConfig.Flags.Dry, "dry", "d", config.DefaultConfig.Flags.Dry, "dry run")
	RootCmd.PersistentFlags().BoolVarP(&config.DefaultConfig.Flags.Root, "root", "r", config.DefaultConfig.Flags.Root, "run as root")
	RootCmd.PersistentFlags().BoolVarP(&config.DefaultConfig.Flags.Force, "force", "f", config.DefaultConfig.Flags.Force, "force init")
	RootCmd.PersistentFlags().StringSliceVarP(&config.DefaultConfig.Flags.Plugins, "plugin", "p", config.DefaultConfig.Flags.Plugins, "plugin")
	RootCmd.PersistentFlags().StringSliceVar(&config.DefaultConfig.Flags.Vars, "var", config.DefaultConfig.Flags.Vars, "variables")

	RootCmd.SilenceErrors = true
	RootCmd.SilenceUsage = true
}

var RootCmd = &cobra.Command{
	Use:   "builder",
	Short: "builder",
	RunE:  runRootCmd,
}

func runRootCmd(cmd *cobra.Command, args []string) error {
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
