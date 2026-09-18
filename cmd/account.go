package cmd

import (
	"github.com/zeiss/builder/cmd/account"

	"github.com/spf13/cobra"
)

func init() {
	AccountCmd.AddCommand(account.LoginCmd)
}

// AccountCmd is the command to manage accounts.
var AccountCmd = &cobra.Command{
	Use:   "account",
	Short: "Manage accounts",
	RunE:  runAccount,
}

func runAccount(cmd *cobra.Command, args []string) error {
	return nil
}
