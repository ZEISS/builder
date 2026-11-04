package sites

import (
	"github.com/spf13/cobra"
)

func init() {
	SitesCmd.AddCommand(DeployCmd)
	SitesCmd.AddCommand(CheckCmd)
	SitesCmd.AddCommand(CreateCmd)
	SitesCmd.AddCommand(ClaimCmd)
}

var SitesCmd = &cobra.Command{
	Use:   "sites",
	Short: "Manages sites",
}
