package account

import (
	"fmt"
	"strings"

	"github.com/zeiss/builder/internal/adapters/oidc"
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/services"

	"github.com/spf13/cobra"
)

// LoginCmd is the command for logging in to builder.
var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to builder",
	RunE:  runLoginCmd,
}

func runLoginCmd(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	flags := config.DefaultFlags

	wellknownConfig, err := config.DiscoverWellKnownConfig(ctx, flags.URL)
	if err != nil {
		return err
	}

	path, err := config.AppDirPath()
	if err != nil {
		return err
	}

	srv, err := services.SetupServices(ctx, path)
	if err != nil {
		return err
	}

	oidcProvider := oidc.New(config.DefaultFlags.URL, config.DefaultClientID, oidc.WithWellKnownConfig(wellknownConfig))
	deviceAuth, err := oidcProvider.Begin(ctx)
	if err != nil {
		return err
	}

	var s strings.Builder

	fmt.Fprintf(&s, "Open the following URL to verify your device: %s\n\n", deviceAuth.VerificationURI)
	fmt.Fprintf(&s, "User code: %s\n", deviceAuth.UserCode)
	fmt.Println(s.String())

	account, err := oidcProvider.Finish(ctx, deviceAuth)
	if err != nil {
		return err
	}
	account.Current = true

	err = srv.CreateAccount(ctx, account)
	if err != nil {
		return err
	}

	fmt.Println("Logged in successfully.")

	return nil
}
