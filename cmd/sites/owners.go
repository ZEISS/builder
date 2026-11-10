package sites

import (
	"fmt"
	"net/http"
	"os"

	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/x/term"
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

var OwnersCmd = &cobra.Command{
	Use:   "owners",
	Short: "Lists the owners of a deployed site",
	RunE:  runOwners,
}

func runOwners(cmd *cobra.Command, _ []string) error {
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

	client, err := apis.NewClientWithResponses(wellknownConfig.ApiURL, apis.WithRequestEditorFn(bearer.Intercept))
	if err != nil {
		return err
	}

	resp, err := client.ListOwnersWithResponse(ctx, spec.Sites.Name)
	if err != nil {
		return err
	}

	if resp.StatusCode() > http.StatusOK {
		return fmt.Errorf("could check the site: %s", spec.Sites.Name)
	}

	tw, _, _ := term.GetSize(os.Stdout.Fd())

	for _, owner := range *resp.JSON200.Owners {
		body := console.Render(tw, console.Success(owner.Email))
		fmt.Fprintln(colorprofile.NewWriter(os.Stderr, os.Environ()), body)
	}

	return nil
}
