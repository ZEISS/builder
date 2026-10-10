package sites

import (
	"errors"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/securityprovider"
	"github.com/zeiss/builder/internal/adapters/db"
	"github.com/zeiss/builder/internal/adapters/oidc"
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/controllers"
	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/ui/models/sites"
	"github.com/zeiss/builder/pkg/apis"
	"github.com/zeiss/builder/server/middlewares/discovery"
	"github.com/zeiss/pkg/cast"
	"github.com/zeiss/pkg/filex"
	"gorm.io/gorm"

	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"
)

var CheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Checks the status of a deployed site",
	RunE:  runCheck,
}

func runCheck(cmd *cobra.Command, _ []string) error {
	err := config.DefaultConfig.LoadSpec()
	if err != nil {
		return err
	}

	path, err := filex.ExpandHomeFolder(config.DefaultConfig.Store)
	if err != nil {
		return err
	}

	err = filex.MkdirAll(filepath.Dir(path), 0o777)
	if err != nil {
		return err
	}

	conn, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return err
	}

	// run migrations, automatically
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

	client, err := apis.NewClientWithResponses(wellknownConfig.ApiURL, apis.WithRequestEditorFn(bearer.Intercept))
	if err != nil {
		return err
	}

	siteCheck := sites.NewCheckSite(cmd.Context(), client, config.DefaultConfig.Spec.Sites.Name)
	_, err = tea.NewProgram(siteCheck, tea.WithContext(cmd.Context())).Run()
	if err != nil {
		return err
	}

	return nil
}
