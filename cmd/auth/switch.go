package auth

import (
	"os"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
	"github.com/katallaxie/pkg/filex"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
	"github.com/zeiss/builder/internal/adapters/db"
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/controllers"
	"github.com/zeiss/builder/internal/ui/models/auth"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var AuthSwitchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch to a different authentication provider",
	RunE:  runAuthSwitch,
}

func runAuthSwitch(cmd *cobra.Command, _ []string) error {
	err := envconfig.Process("", &config.DefaultConfig.Flags.AuthFlags)
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

	store := db.New(conn)
	accountCtrl := controllers.NewAccountController(store)

	// clear all the stdout output
	os.Stdout.WriteString("\x1b[2J\x1b[3J\x1b[H")

	_, err = tea.NewProgram(auth.NewSwitch(cmd.Context(), accountCtrl), tea.WithContext(cmd.Context())).Run()
	if err != nil {
		return err
	}

	return nil
}
