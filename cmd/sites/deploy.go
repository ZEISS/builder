package sites

import (
	"os"
	"path/filepath"

	"github.com/zeiss/builder/internal/adapters/client"
	"github.com/zeiss/builder/internal/adapters/db"
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/controllers"
	"github.com/zeiss/builder/internal/ui/models/sites"
	"github.com/zeiss/builder/pkg/apis"

	tea "charm.land/bubbletea/v2"
	"github.com/katallaxie/pkg/filex"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
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

	if err := db.RunMigrations(conn); err != nil {
		return err
	}

	// accountStore := db.New(conn)
	// accountController := controllers.NewAccountController(accountStore)

	// account := &models.Account{}
	// err = accountController.GetCurrent(cmd.Context(), account)
	// if err != nil {
	// 	return err
	// }

	// bearer, err := securityprovider.NewSecurityProviderBearerToken(cast.Value(account.IDToken))
	// if err != nil {
	// 	return err
	// }

	// c, err := apis.NewClient(config.DefaultConfig.URL, apis.WithRequestEditorFn(bearer.Intercept))
	// if err != nil {
	// 	return err
	// }

	c, err := apis.NewClient(config.DefaultConfig.URL)
	if err != nil {
		return err
	}

	sitesRepo := client.New(c)
	sitesController := controllers.NewSitesController(sitesRepo)

	// clear all the stdout output
	os.Stdout.WriteString("\x1b[2J\x1b[3J\x1b[H")

	_, err = tea.NewProgram(sites.NewDeploy(cmd.Context(), config.DefaultConfig, sitesController), tea.WithContext(cmd.Context())).Run()
	if err != nil {
		return err
	}

	return nil
}
