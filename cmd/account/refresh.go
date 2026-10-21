package account

import (
	"fmt"

	"github.com/zeiss/builder/internal/adapters/oidc"
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/services"

	"github.com/spf13/cobra"
)

// RefreshCmd is the command for refreshing the account.
var RefreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh the account",
	RunE:  runRefreshCmd,
}

func runRefreshCmd(cmd *cobra.Command, args []string) error {
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

	account := &models.Account{Current: true}
	if err := srv.GetAccount(ctx, account); err != nil {
		return err
	}

	if err := oidcProvider.Refresh(ctx, account); err != nil {
		return err
	}

	err = srv.UpdateAccount(ctx, account)
	if err != nil {
		return err
	}

	fmt.Println("Refreshed account successfully.")

	return nil
}
