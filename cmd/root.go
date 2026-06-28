package cmd

import (
	"fmt"

	"github.com/zeiss/builder/cmd/auth"
	"github.com/zeiss/builder/cmd/sites"
	"github.com/zeiss/builder/internal/config"

	"github.com/spf13/cobra"
)

var cfg = config.New()

const (
	versionFmt = "%s (%s %s)"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func init() {
	RootCmd.AddCommand(InitCmd)
	RootCmd.AddCommand(TaskCmd)
	RootCmd.AddCommand(sites.SitesCmd)
	RootCmd.AddCommand(auth.AuthCmd)

	RootCmd.PersistentFlags().StringVarP(&config.DefaultConfig.URL, "url", "u", config.DefaultConfig.URL, "URL")
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
	Use:     "builder",
	Short:   "builder",
	Version: fmt.Sprintf(versionFmt, version, commit, date),
}
