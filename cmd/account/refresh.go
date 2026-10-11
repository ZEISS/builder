package account

import (
	"path/filepath"

	"github.com/zeiss/builder/internal/adapters/db"
	"github.com/zeiss/builder/internal/adapters/oidc"
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/controllers"
	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/server/middlewares/discovery"

	"github.com/glebarez/sqlite"
	"github.com/spf13/cobra"
	"github.com/zeiss/pkg/filex"
	"gorm.io/gorm"
)

// RefreshCmd is the command for refreshing the account.
var RefreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh the account",
	RunE:  runRefreshCmd,
}

func runRefreshCmd(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

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

	discv := discovery.NewClient()
	wellknownConfig, err := discv.Discover(ctx, config.DefaultConfig.Flags.URL)
	if err != nil {
		return err
	}

	oidcProvider := oidc.New(config.DefaultConfig.Flags.URL, config.DefaultClientID, oidc.WithWellKnownConfig(wellknownConfig))

	store := db.New(conn)
	accountCtrl := controllers.NewAccountController(config.DefaultConfig, store)
	authCtrl := controllers.NewDeviceAuthController(oidcProvider, store)

	current := &models.Account{}
	err = accountCtrl.GetCurrent(ctx, current)
	if err != nil {
		return err
	}

	err = authCtrl.Refresh(ctx, current)
	if err != nil {
		return err
	}

	err = accountCtrl.Update(ctx, current)
	if err != nil {
		return err
	}

	return nil
}
