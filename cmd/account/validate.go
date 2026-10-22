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

	"github.com/spf13/cobra"
)

// ValidateCmd is the command for validating the token.
var ValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Check if the token is valid",
	Long: `
	  Check if the token is valid and not expired.
	`,
	RunE: runValidateCmd,
}

func runValidateCmd(cmd *cobra.Command, args []string) error {
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
	body := console.Render(tw, console.Success(fmt.Sprintf("Token is valid until: %s", account.ExpiresAt.Local())))
	fmt.Fprintln(colorprofile.NewWriter(os.Stderr, os.Environ()), body)

	return nil
}
