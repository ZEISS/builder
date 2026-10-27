package account

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/x/term"
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/oidc"
	"github.com/zeiss/builder/internal/services"
	"github.com/zeiss/builder/internal/ui/console"

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

	path, err := config.AppDirPath()
	if err != nil {
		return err
	}

	srv, err := services.SetupServices(ctx, path)
	if err != nil {
		return err
	}

	oidcProvider, err := oidc.SetupOidcProvider(ctx, flags.URL, config.DefaultClientID)
	if err != nil {
		return err
	}

	deviceAuth, err := oidcProvider.Begin(ctx)
	if err != nil {
		return err
	}

	var s strings.Builder

	fmt.Fprintf(&s, "Open the following URL to verify your device: %s\n\n", deviceAuth.VerificationURI)
	fmt.Fprintf(&s, "User code: %s\n", deviceAuth.UserCode)
	fmt.Println(s.String())

	token, err := oidcProvider.Finish(ctx, deviceAuth)
	if err != nil {
		return err
	}

	account := models.FromToken(token)
	err = srv.CreateAccount(ctx, account)
	if err != nil {
		return err
	}

	tw, _, _ := term.GetSize(os.Stdout.Fd())
	body := console.Render(tw, console.Success(fmt.Sprintf("Logged in %s successfully", account.Email)))
	fmt.Fprintln(colorprofile.NewWriter(os.Stderr, os.Environ()), body)

	return nil
}
