package sites

import (
	"errors"
	"os"

	"github.com/oapi-codegen/oapi-codegen/v2/pkg/securityprovider"
	"github.com/zeiss/builder/internal/adapters/client"
	"github.com/zeiss/builder/internal/adapters/db"
	"github.com/zeiss/builder/internal/adapters/oidc"
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/controllers"
	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/ui/models/sites"
	"github.com/zeiss/builder/pkg/apis"
	"github.com/zeiss/builder/server/middlewares/discovery"

	tea "charm.land/bubbletea/v2"
	"github.com/glebarez/sqlite"
	"github.com/spf13/cobra"
	"github.com/zeiss/pkg/cast"
	"gorm.io/gorm"
)

var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys a site",
	RunE:  runDeploy,
}

func runDeploy(cmd *cobra.Command, _ []string) error {
	err := config.DefaultConfig.LoadSpec()
	if err != nil {
		return err
	}

	path, err := config.ExpandConfigPath(config.DefaultConfig.Store)
	if err != nil {
		return err
	}

	conn, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return err
	}

	if err := db.RunMigrations(conn); err != nil {
		return err
	}

	discv := discovery.NewClient()
	wellknownConfig, err := discv.Discover(cmd.Context(), config.DefaultConfig.Flags.URL)
	if err != nil {
		return err
	}

	oidcProvider := oidc.New(config.DefaultConfig.Flags.URL, config.DefaultClientID, oidc.WithWellKnownConfig(wellknownConfig))

	store := db.New(conn)
	accountController := controllers.NewAccountController(config.DefaultConfig, store)
	authCtrl := controllers.NewDeviceAuthController(oidcProvider, store)

	account := &models.Account{}
	err = accountController.GetCurrent(cmd.Context(), account)
	if err != nil {
		return errors.New("no account found")
	}

	err = authCtrl.Refresh(cmd.Context(), account)
	if err != nil {
		return err
	}

	bearer, err := securityprovider.NewSecurityProviderBearerToken(cast.Value(account.IDToken))
	if err != nil {
		return err
	}

	c, err := apis.NewClientWithResponses(wellknownConfig.ApiURL, apis.WithRequestEditorFn(bearer.Intercept))
	if err != nil {
		return err
	}

	client := client.New(c)
	sitesController := controllers.NewSitesController(client)
	filesController := controllers.NewFilesController(client)

	// clear all the stdout output
	os.Stdout.WriteString("\x1b[2J\x1b[3J\x1b[H")

	model := sites.NewDeploy(cmd.Context(), config.DefaultConfig, sitesController, filesController)
	_, err = tea.NewProgram(model, tea.WithContext(cmd.Context())).Run()
	if err != nil {
		return err
	}

	return nil
}
