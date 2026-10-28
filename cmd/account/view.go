package account

import (
	"fmt"
	"os"

	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/x/term"
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/services"
	"github.com/zeiss/builder/internal/ui/console"
	"github.com/zeiss/pkg/cast"

	"github.com/spf13/cobra"
)

// Show is a command that shows the current account.
var ShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the current account",
	Long: `
	  Show the current account and the token.
	`,
	RunE: runShowCmd,
}

func runShowCmd(cmd *cobra.Command, args []string) error {
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

	tw, _, _ := term.GetSize(os.Stdout.Fd())
	body := console.Render(tw,
		console.Success(fmt.Sprintf("Your are logged in as: %s", account.Email)),
		console.Success(fmt.Sprintf("That's your access token: \n\n%s", cast.Value(account.AccessToken))),
	)
	fmt.Fprintln(colorprofile.NewWriter(os.Stderr, os.Environ()), body)

	return nil
}
