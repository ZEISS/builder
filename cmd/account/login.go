package account

import (
	"path/filepath"

	"github.com/zeiss/builder/internal/adapters/db"
	"github.com/zeiss/builder/internal/adapters/oidc"
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/controllers"
	"github.com/zeiss/builder/internal/ui/models/account"

	tea "charm.land/bubbletea/v2"
	"github.com/glebarez/sqlite"
	"github.com/spf13/cobra"
	"github.com/zeiss/pkg/filex"
	"gorm.io/gorm"
)

// LoginCmd is the command for logging in to builder.
var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to builder",
	RunE:  runLoginCmd,
}

func runLoginCmd(cmd *cobra.Command, args []string) error {
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

	store := db.New(conn)
	oidcProvider := oidc.New(config.DefaultConfig.URL, config.DefaultConfig.Flags.AuthFlags.ClientID)

	accountCtrl := controllers.NewAccountController(config.DefaultConfig, store)
	authCtrl := controllers.NewDeviceAuthController(oidcProvider, store)

	app := account.New(ctx, authCtrl, accountCtrl)
	program := tea.NewProgram(app, tea.WithContext(cmd.Context()))

	_, err = program.Run()
	if err != nil {
		return err
	}

	return nil
}
