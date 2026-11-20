package owners

import (
	"fmt"
	"os"

	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/x/term"
	"github.com/zeiss/builder/internal/adapters/client"
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/oidc"
	"github.com/zeiss/builder/internal/services"
	"github.com/zeiss/builder/internal/ui/console"
	"github.com/zeiss/builder/pkg/apis"

	"github.com/oapi-codegen/oapi-codegen/v2/pkg/securityprovider"
	"github.com/spf13/cobra"
	"github.com/zeiss/pkg/cast"
)

type RemoveCmdFlags struct {
	// Email of the owners to add
	Email string `json:"email" yaml:"email" example:"user@example.com"`
}

var removeCmdFlags = RemoveCmdFlags{}

var RemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes an owner from a deployed site",
	RunE:  runRemove,
}

func runRemove(cmd *cobra.Command, _ []string) error {
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

	spec, err := flags.Spec()
	if err != nil {
		return err
	}

	if err := spec.Validate(); err != nil {
		return err
	}

	srv, err := services.SetupServices(ctx, path)
	if err != nil {
		return err
	}

	account := &models.Account{Current: true}
	if err := srv.GetAccount(ctx, account); err != nil {
		return err
	}

	oidcProvider, err := oidc.SetupOidcProvider(ctx, flags.URL, config.DefaultClientID)
	if err != nil {
		return err
	}

	if err := oidcProvider.Refresh(ctx, account.ToToken()); err != nil {
		return err
	}

	bearer, err := securityprovider.NewSecurityProviderBearerToken(cast.Value(account.IDToken))
	if err != nil {
		return err
	}

	api, err := apis.NewClientWithResponses(wellknownConfig.ApiURL, apis.WithRequestEditorFn(bearer.Intercept))
	if err != nil {
		return err
	}

	c := client.New(api)
	if err := c.RemoveOwner(ctx, spec.Sites.Name, removeCmdFlags.Email); err != nil {
		return err
	}

	tw, _, _ := term.GetSize(os.Stdout.Fd())
	body := console.Render(tw, console.Success(fmt.Sprintf("Removed %s successfully from %s", removeCmdFlags.Email, spec.Sites.Name)))
	fmt.Fprintln(colorprofile.NewWriter(os.Stderr, os.Environ()), body)

	return nil
}
