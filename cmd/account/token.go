package account

import (
	"log"

	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/services"

	"github.com/spf13/cobra"
	"github.com/zeiss/pkg/cast"
)

// TokenCmd is the command for refreshing the account.
var TokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Token of the current account",
	RunE:  runTokenCmd,
}

func runTokenCmd(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	path, err := config.AppDirPath()
	if err != nil {
		return err
	}

	srv, err := services.SetupServices(ctx, path)
	if err != nil {
		return err
	}

	account := &models.Account{Current: true}
	err = srv.GetAccount(ctx, account)
	if err != nil {
		return err
	}

	log.Println(cast.Value(account.IDToken))

	return nil
}
