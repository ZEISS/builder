package sites

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/x/term"
	"github.com/zeiss/builder/internal/adapters/client"
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/services"
	"github.com/zeiss/builder/internal/ui/console"
	"github.com/zeiss/builder/pkg/apis"
	"github.com/zeiss/builder/pkg/utils"

	"github.com/oapi-codegen/oapi-codegen/v2/pkg/securityprovider"
	"github.com/spf13/cobra"
	"github.com/zeiss/pkg/cast"
)

var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys a site",
	RunE:  runDeploy,
}

func runDeploy(cmd *cobra.Command, _ []string) error {
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

	bearer, err := securityprovider.NewSecurityProviderBearerToken(cast.Value(account.IDToken))
	if err != nil {
		return err
	}

	api, err := apis.NewClientWithResponses(wellknownConfig.ApiURL, apis.WithRequestEditorFn(bearer.Intercept))
	if err != nil {
		return err
	}

	c := client.New(api)

	cwd, err := config.Cwd()
	if err != nil {
		return err
	}

	site, err := c.GetSite(ctx, spec.Sites.Name)
	if err != nil {
		return err
	}

	scanDir := filepath.Base(filepath.Join(cwd, spec.Sites.Path))
	files := utils.ScanDir(scanDir, spec.Sites.Ignore)

	for _, file := range files {
		file = strings.TrimPrefix(filepath.Clean(strings.TrimPrefix(file, spec.Sites.Path)), "/")
		err := c.UploadFile(ctx, site.ID, scanDir, file)
		if err != nil {
			return err
		}
	}

	tw, _, _ := term.GetSize(os.Stdout.Fd())
	body := console.Render(tw, console.Success(fmt.Sprintf("%d files uploaded successfully.", len(files))))
	fmt.Fprintln(colorprofile.NewWriter(os.Stderr, os.Environ()), body)

	return nil
}
