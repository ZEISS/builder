package sites

import (
	"github.com/spf13/cobra"
	"github.com/zeiss/builder/cmd/sites/owners"
)

func init() {
	SitesCmd.AddCommand(CheckCmd)
	SitesCmd.AddCommand(CreateCmd)
	SitesCmd.AddCommand(DeleteCmd)
	SitesCmd.AddCommand(DeployCmd)
	SitesCmd.AddCommand(owners.RootCmd)
}

var SitesCmd = &cobra.Command{
	Use:   "sites",
	Short: "Manages sites",
}
