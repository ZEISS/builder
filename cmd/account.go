package cmd

import (
	"github.com/zeiss/builder/cmd/account"

	"github.com/spf13/cobra"
)

func init() {
	AccountCmd.AddCommand(account.LoginCmd)
	AccountCmd.AddCommand(account.SwitchCmd)
	AccountCmd.AddCommand(account.RefreshCmd)
	AccountCmd.AddCommand(account.TokenCmd)
}

// AccountCmd is the command to manage accounts.
var AccountCmd = &cobra.Command{
	Use:   "account",
	Short: "Manage accounts",
}
