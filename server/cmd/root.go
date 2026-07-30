package cmd

import (
	"context"

	"github.com/zeiss/builder/server/adapters/database"
	"github.com/zeiss/builder/server/adapters/files"
	"github.com/zeiss/builder/server/adapters/handlers"
	"github.com/zeiss/builder/server/configs"
	"github.com/zeiss/builder/server/controllers"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/joho/godotenv"
	goth "github.com/katallaxie/fiber-goth/v3"
	gorm_adapter "github.com/katallaxie/fiber-goth/v3/adapters/gorm"
	"github.com/katallaxie/fiber-goth/v3/providers"
	"github.com/katallaxie/fiber-goth/v3/providers/dex"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
	"github.com/zeiss/pkg/filex"
	"github.com/zeiss/pkg/server"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	cobra.OnInitialize(initConfig)

	Root.PersistentFlags().StringVar(&configs.DefaultConfig.Flags.Addr, "addr", configs.DefaultConfig.Flags.Addr, "addr")
	Root.PersistentFlags().StringVar(&configs.DefaultConfig.Flags.OIDCIssuer, "oidc-issuer", configs.DefaultConfig.Flags.OIDCIssuer, "OIDC Issuer")
	Root.PersistentFlags().StringVar(&configs.DefaultConfig.Flags.OIDCAudience, "oidc-audience", configs.DefaultConfig.Flags.OIDCAudience, "OIDC Audience")
	Root.PersistentFlags().StringVar(&configs.DefaultConfig.Flags.Domain, "domain", configs.DefaultConfig.Flags.Domain, "domain")

	// Configure the files path
	Root.PersistentFlags().StringVar(&configs.DefaultConfig.Flags.FilesFlags.Path, "files-path", configs.DefaultConfig.Flags.FilesFlags.Path, "Files Path")

	Root.PersistentFlags().StringVar(&configs.DefaultConfig.Flags.DexFlags.CallbackURL, "dex-callback-url", configs.DefaultConfig.Flags.DexFlags.CallbackURL, "Dex Callback URL")
	Root.PersistentFlags().StringVar(&configs.DefaultConfig.Flags.DexFlags.ClientID, "dex-client-id", configs.DefaultConfig.Flags.DexFlags.ClientID, "Dex Client ID")
	Root.PersistentFlags().StringVar(&configs.DefaultConfig.Flags.DexFlags.ClientSecret, "dex-client-secret", configs.DefaultConfig.Flags.DexFlags.ClientSecret, "Dex Client Secret")
	Root.PersistentFlags().StringVar(&configs.DefaultConfig.Flags.DexFlags.LoginURL, "dex-login-url", configs.DefaultConfig.Flags.DexFlags.LoginURL, "Dex Login URL")

	Root.PersistentFlags().BoolVar(&configs.DefaultConfig.Flags.SqliteFlags.Enabled, "sqlite", configs.DefaultConfig.Flags.SqliteFlags.Enabled, "SQLite Enabled")
	Root.PersistentFlags().StringVar(&configs.DefaultConfig.Flags.SqliteFlags.Database, "sqlite-database", configs.DefaultConfig.Flags.SqliteFlags.Database, "SQLite Database")
	Root.PersistentFlags().StringVar(&configs.DefaultConfig.Flags.SqliteFlags.Path, "sqlite-path", configs.DefaultConfig.Flags.SqliteFlags.Path, "SQLite Path")

	Root.SilenceUsage = true
}

func initConfig() {
	err := godotenv.Load()
	cobra.CheckErr(err)

	err = envconfig.Process("", configs.DefaultConfig.Flags)
	cobra.CheckErr(err)
}

var Root = &cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		srv := NewWebSrv(configs.DefaultConfig)

		s, _ := server.WithContext(cmd.Context())
		s.Listen(srv, false)

		return s.Wait()
	},
}

var _ server.Listener = (*WebSrv)(nil)

// WebSrv is the server that implements the Noop interface.
type WebSrv struct {
	cfg *configs.Config
}

// NewWebSrv returns a new instance of NoopSrv.
func NewWebSrv(cfg *configs.Config) *WebSrv {
	return &WebSrv{cfg}
}

// Host holds a Fiber app instance.
type Host struct {
	Fiber *fiber.App
}

// Start starts the server.
func (s *WebSrv) Start(ctx context.Context, ready server.ReadyFunc, run server.RunFunc) func() error {
	return func() error {
		// Create files folder if not exists
		err := filex.MkdirAll(s.cfg.Flags.FilesFlags.Path, 0o755)
		if err != nil {
			return err
		}

		conn, err := gorm.Open(sqlite.Open(s.cfg.Flags.SqliteFlags.Path), &gorm.Config{
			TranslateError: true,
		})
		if err != nil {
			return err
		}

		err = database.RunMigrations(conn)
		if err != nil {
			return err
		}

		db := database.NewDatabase(conn)
		fs := files.NewFiles(s.cfg)

		providers.RegisterProvider(dex.New(s.cfg.Flags.DexFlags.ClientID, s.cfg.Flags.DexFlags.ClientSecret, s.cfg.Flags.OIDCIssuer, s.cfg.Flags.DexFlags.CallbackURL))

		// fs := store.NewFS(s.cfg.Flags.FilesFlags.Path)
		sitesCtrl := controllers.NewSitesController(db)
		filesCtrl := controllers.NewFilesController(fs)
		sitesHandler := handlers.NewSitesHandler(sitesCtrl, filesCtrl)

		c := fiber.Config{}

		app := fiber.New(c)
		app.Use(requestid.New())
		app.Use(logger.New())

		ga := gorm_adapter.New(conn)

		gothConfig := goth.Config{
			Adapter:        ga,
			Secret:         goth.GenerateKey(),
			CookieHTTPOnly: true,
			LoginURL:       s.cfg.Flags.DexFlags.LoginURL,
			CookieDomain:   s.cfg.Flags.Domain,
		}

		root := app.Domain(s.cfg.Flags.Domain)
		root.Get("/session", goth.NewSessionHandler(gothConfig))
		root.Get("/login/:provider", goth.NewBeginAuthHandler(gothConfig))
		root.Get("/auth/:provider/callback", goth.NewCompleteAuthHandler(gothConfig))
		root.Get("/logout", goth.NewLogoutHandler(gothConfig))

		// sites := app.Domain(":site." + cfg.Flags.Domain)
		// config := static.Config{
		// 	Root: http.Dir(cfg.Flags.c),
		// }
		// sites.Use(goth.Session(gothConfig))
		// sites.Use(goth.Protect(gothConfig))
		// sites.Use(static.New(config))

		api := app.Group("/api")
		v1 := api.Group("/v1")

		apiConfig := huma.DefaultConfig("Builder API", "1.0.0")
		apiConfig.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
			"openid": {
				Type:             "openIdConnect",
				In:               "header",
				BearerFormat:     "JWT",
				OpenIDConnectURL: s.cfg.Flags.OIDCIssuer,
				Flows: &huma.OAuthFlows{
					AuthorizationCode: &huma.OAuthFlow{
						Scopes: map[string]string{
							"openid":         "",
							"profile":        "",
							"email":          "",
							"offline_access": "",
						},
					},
				},
			},
		}

		spec := humafiber.NewWithGroup(app, v1, apiConfig)
		// spec.UseMiddleware(oidc.NewAuthMiddleware(spec, cfg.Flags.OIDCIssuer, cfg.Flags.OIDCAudience))
		sitesHandler.Register(spec)

		err = app.Listen(s.cfg.Flags.Addr)
		if err != nil {
			return err
		}

		return nil
	}
}
