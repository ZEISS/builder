package auth

import (
	"github.com/spf13/cobra"
	"github.com/zeiss/builder/internal/config"
)

func init() {
	AuthCmd.AddCommand(AuthLoginCmd)
	AuthCmd.AddCommand(AuthSwitchCmd)
	AuthCmd.AddCommand(AuthTokenCmd)

	AuthCmd.PersistentFlags().StringVar(&config.DefaultConfig.Flags.AuthFlags.ClientID, "client-id", "", "OIDC client id")
	AuthCmd.PersistentFlags().StringVar(&config.DefaultConfig.Flags.AuthFlags.ClientURL, "client-url", "", "OIDC client url")
}

var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate the builder (default: oidc)",
}
