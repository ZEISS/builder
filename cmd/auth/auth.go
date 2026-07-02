package auth

import (
	"github.com/spf13/cobra"
	"github.com/zeiss/builder/internal/config"
)

func init() {
	AuthCmd.AddCommand(AuthLoginCmd)
	AuthCmd.AddCommand(AuthSwitchCmd)
	AuthCmd.AddCommand(AuthTokenCmd)

	AuthCmd.PersistentFlags().BoolVar(&config.DefaultConfig.Flags.AuthFlags.Dex, "dex", true, "Enable the Dex as provider")
	AuthCmd.PersistentFlags().StringVar(&config.DefaultConfig.Flags.AuthFlags.DexClientID, "dex-client-id", "", "Dex client id")
	AuthCmd.PersistentFlags().StringVar(&config.DefaultConfig.Flags.AuthFlags.DexClientSecret, "dex-client-secret", "", "Dex client secret")
	AuthCmd.PersistentFlags().StringVar(&config.DefaultConfig.Flags.AuthFlags.DexClientURL, "dex-client-url", "http://127.0.0.1:5556/dex", "Dex client url")
}

var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate the builder (default: dex)",
}
