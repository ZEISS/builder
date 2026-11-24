package account

import (
	"fmt"
	"os"

	"github.com/zeiss/builder/internal/adapters/client"
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/services"
	"github.com/zeiss/builder/internal/ui/console"

	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/x/term"
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

	path, err := config.AppDirPath()
	if err != nil {
		return err
	}

	srv, err := services.SetupServices(ctx, path)
	if err != nil {
		return err
	}

	_, err = client.SetupClient(ctx, srv, flags.URL)
	if err != nil {
		return err
	}

	account := &models.Account{Current: true}
	if err := srv.GetAccount(ctx, account); err != nil {
		return err
	}

	tw, _, _ := term.GetSize(os.Stdout.Fd())
	body := console.Render(tw, console.Success(fmt.Sprintf("Refreshed account %s successfully", account.Email)))
	fmt.Fprintln(colorprofile.NewWriter(os.Stderr, os.Environ()), body)

	return nil
}
