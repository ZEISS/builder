package cmds

import (
	"context"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
	"github.com/glebarez/sqlite"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/securityprovider"
	"github.com/zeiss/builder/internal/adapters/db"
	"github.com/zeiss/builder/internal/config"
	"github.com/zeiss/builder/internal/controllers"
	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/pkg/apis"
	"github.com/zeiss/builder/server/middlewares/discovery"
	"github.com/zeiss/pkg/cast"
	"github.com/zeiss/pkg/filex"
	"gorm.io/gorm"
)

// InitMsg is sent when the user initializes the application.
type InitMsg struct {
	Client *apis.ClientWithResponses
	Config config.Config
}

// Init is the command to initialize the application.
func Init(ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		err := config.DefaultConfig.LoadSpec()
		if err != nil {
			return ErrorMsg{Err: err}
		}

		path, err := filex.ExpandHomeFolder(config.DefaultConfig.Store)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		err = filex.MkdirAll(filepath.Dir(path), 0o777)
		if err != nil {
			return ErrorMsg{Err: err}
		}

		conn, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
		if err != nil {
			return ErrorMsg{Err: err}
		}

		if err := db.RunMigrations(conn); err != nil {
			return ErrorMsg{Err: err}
		}

		discv := discovery.NewClient()
		wellknownConfig, err := discv.Discover(ctx, config.DefaultConfig.Flags.URL)
		if err != nil {
			return err
		}

		accountStore := db.New(conn)
		accountController := controllers.NewAccountController(config.DefaultConfig, accountStore)

		account := &models.Account{}
		err = accountController.GetCurrent(ctx, account)
		if err != nil {
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

		return InitMsg{api, config.DefaultConfig}
	}
}
