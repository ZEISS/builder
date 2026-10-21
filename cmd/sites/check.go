package sites

import (
	"fmt"
	"net/http"

	"github.com/zeiss/builder/internal/apps"
	"github.com/zeiss/builder/pkg/apis"

	"github.com/oapi-codegen/oapi-codegen/v2/pkg/securityprovider"
	"github.com/spf13/cobra"
	"github.com/zeiss/pkg/cast"
)

var CheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Checks the status of a deployed site",
	RunE:  runCheck,
}

func runCheck(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()

	app, err := apps.SetupApp(ctx)
	if err != nil {
		return err
	}

	if err := app.Initialize(ctx); err != nil {
		return err
	}

	spec := app.Spec()
	wellknownConfig := app.WellKnownConfig()

	account := app.Account()
	bearer, err := securityprovider.NewSecurityProviderBearerToken(cast.Value(account.IDToken))
	if err != nil {
		return err
	}

	client, err := apis.NewClientWithResponses(wellknownConfig.ApiURL, apis.WithRequestEditorFn(bearer.Intercept))
	if err != nil {
		return err
	}

	params := &apis.GetSiteParams{Name: spec.Sites.Name}
	resp, err := client.GetSiteWithResponse(ctx, params)
	if err != nil {
		return err
	}

	if resp.StatusCode() > http.StatusOK {
		return fmt.Errorf("could check the site: %s", spec.Sites.Name)
	}

	fmt.Printf("site exists: %s", spec.Sites.Name)

	return nil
}
