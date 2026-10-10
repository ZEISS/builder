package account

import (
	"log"
	"path/filepath"

	"github.com/zeiss/builder/internal/adapters/db"
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/controllers"
	"github.com/zeiss/builder/internal/models"

	"github.com/glebarez/sqlite"
	"github.com/spf13/cobra"
	"github.com/zeiss/pkg/cast"
	"github.com/zeiss/pkg/filex"
	"gorm.io/gorm"
)

// TokenCmd is the command for refreshing the account.
var TokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Token of the current account",
	RunE:  runTokenCmd,
}

func runTokenCmd(cmd *cobra.Command, args []string) error {
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
	accountCtrl := controllers.NewAccountController(config.DefaultConfig, store)

	current := &models.Account{}
	err = accountCtrl.GetCurrent(ctx, current)
	if err != nil {
		return err
	}

	log.Println(cast.Value(current.AccessToken))

	return nil
}
