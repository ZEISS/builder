package cmd

import (
	"context"
	"fmt"

	"github.com/zeiss/builder/server/adapters/database"
	"github.com/zeiss/builder/server/adapters/files"
	"github.com/zeiss/builder/server/adapters/handlers"
	"github.com/zeiss/builder/server/configs"
	"github.com/zeiss/builder/server/controllers"
	"github.com/zeiss/builder/server/middlewares/auth/oidc"
	"github.com/zeiss/builder/server/middlewares/healthz"
	sites_middleware "github.com/zeiss/builder/server/middlewares/sites"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/gofiber/storage/azureblob/v2"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
	goth "github.com/zeiss/fiber-goth/v3"
	gorm_adapter "github.com/zeiss/fiber-goth/v3/adapters/gorm"
	"github.com/zeiss/fiber-goth/v3/providers"
	"github.com/zeiss/fiber-goth/v3/providers/dex"
	"github.com/zeiss/pkg/dbx/pg"
	"github.com/zeiss/pkg/server"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	versionFmt = "%s (%s %s)"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func init() {
	cobra.OnInitialize(initConfig)

	Root.PersistentFlags().StringVar(&configs.DefaultConfig.Flags.Addr, "addr", configs.DefaultConfig.Flags.Addr, "addr")
	Root.PersistentFlags().StringVar(&configs.DefaultConfig.Flags.OIDCIssuer, "oidc-issuer", configs.DefaultConfig.Flags.OIDCIssuer, "OIDC Issuer")
	Root.PersistentFlags().StringVar(&configs.DefaultConfig.Flags.OIDCAudience, "oidc-audience", configs.DefaultConfig.Flags.OIDCAudience, "OIDC Audience")
	Root.PersistentFlags().StringVar(&configs.DefaultConfig.Flags.Domain, "domain", configs.DefaultConfig.Flags.Domain, "domain")

	Root.PersistentFlags().StringVar(&configs.DefaultConfig.Flags.DexFlags.CallbackURL, "dex-callback-url", configs.DefaultConfig.Flags.DexFlags.CallbackURL, "Dex Callback URL")
	Root.PersistentFlags().StringVar(&configs.DefaultConfig.Flags.DexFlags.ClientID, "dex-client-id", configs.DefaultConfig.Flags.DexFlags.ClientID, "Dex Client ID")
	Root.PersistentFlags().StringVar(&configs.DefaultConfig.Flags.DexFlags.ClientSecret, "dex-client-secret", configs.DefaultConfig.Flags.DexFlags.ClientSecret, "Dex Client Secret")
	Root.PersistentFlags().StringVar(&configs.DefaultConfig.Flags.DexFlags.LoginURL, "dex-login-url", configs.DefaultConfig.Flags.DexFlags.LoginURL, "Dex Login URL")

	Root.PersistentFlags().BoolVar(&configs.DefaultConfig.Flags.PostgresFlags.Enabled, "postgres", configs.DefaultConfig.Flags.PostgresFlags.Enabled, "PostgreSQL Enabled")
	Root.PersistentFlags().StringVar(&configs.DefaultConfig.Flags.PostgresFlags.Host, "postgres-host", configs.DefaultConfig.Flags.PostgresFlags.Host, "PostgreSQL Host")
	Root.PersistentFlags().IntVar(&configs.DefaultConfig.Flags.PostgresFlags.Port, "postgres-port", configs.DefaultConfig.Flags.PostgresFlags.Port, "PostgreSQL Port")
	Root.PersistentFlags().StringVar(&configs.DefaultConfig.Flags.PostgresFlags.Database, "postgres-database", configs.DefaultConfig.Flags.PostgresFlags.Database, "PostgreSQL Database")
	Root.PersistentFlags().StringVar(&configs.DefaultConfig.Flags.PostgresFlags.User, "postgres-user", configs.DefaultConfig.Flags.PostgresFlags.User, "PostgreSQL User")
	Root.PersistentFlags().StringVar(&configs.DefaultConfig.Flags.PostgresFlags.Password, "postgres-password", configs.DefaultConfig.Flags.PostgresFlags.Password, "PostgreSQL Password")
	Root.PersistentFlags().StringVar(&configs.DefaultConfig.Flags.PostgresFlags.SSLMode, "postgres-ssl-mode", configs.DefaultConfig.Flags.PostgresFlags.SSLMode, "PostgreSQL SSL Mode")

	Root.SilenceUsage = true
}

func initConfig() {
	_ = godotenv.Load() // ignore error

	err := envconfig.Process("", configs.DefaultConfig.Flags)
	cobra.CheckErr(err)
}

var Root = &cobra.Command{
	Short:   "Server is the backend to the Builder CLI",
	Version: fmt.Sprintf(versionFmt, version, commit, date),
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
		dsn := pg.NewConfig()
		dsn.Database = s.cfg.Flags.PostgresFlags.Database
		dsn.Host = s.cfg.Flags.PostgresFlags.Host
		dsn.Port = s.cfg.Flags.PostgresFlags.Port
		dsn.User = s.cfg.Flags.PostgresFlags.User
		dsn.Password = s.cfg.Flags.PostgresFlags.Password
		dsn.SslMode = s.cfg.Flags.PostgresFlags.SSLMode

		conn, err := gorm.Open(postgres.Open(dsn.FormatDSN()), &gorm.Config{
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
		providers.RegisterProvider(dex.New(s.cfg.Flags.DexFlags.ClientID, s.cfg.Flags.DexFlags.ClientSecret, s.cfg.Flags.OIDCIssuer, s.cfg.Flags.DexFlags.CallbackURL))

		// This is one and only option that we have to use right now.
		storage := azureblob.New(
			azureblob.Config{
				Account:   s.cfg.Flags.AzureBlobFlags.AccountName,
				Container: s.cfg.Flags.AzureBlobFlags.ContainerName,
				Endpoint:  s.cfg.Flags.AzureBlobFlags.Endpoint,
				Credentials: azureblob.Credentials{
					Account: s.cfg.Flags.AzureBlobFlags.CredentialsAccount,
					Key:     s.cfg.Flags.AzureBlobFlags.CredentialsKey,
				},
			},
		)

		files := files.New(storage)
		sitesCtrl := controllers.NewSitesController(db)
		filesCtrl := controllers.NewFilesController(files)
		sitesHandler := handlers.NewSitesHandler(sitesCtrl, filesCtrl)

		c := fiber.Config{}
		app := fiber.New(c)
		app.Use(requestid.New())
		app.Use(logger.New())

		app.Get("/healthz", healthz.New())

		ga := gorm_adapter.New(conn)

		gothConfig := goth.Config{
			Adapter:        ga,
			Secret:         goth.GenerateKey(),
			CookieHTTPOnly: true,
			LoginURL:       s.cfg.Flags.DexFlags.LoginURL,
			CookieDomain:   s.cfg.Flags.Domain,
		}

		root := app.Domain(s.cfg.Flags.Domain)
		root.Use(goth.Session(gothConfig))

		root.Get("/session", goth.NewSessionHandler(gothConfig))
		root.Get("/login/:provider", goth.NewBeginAuthHandler(gothConfig))
		root.Get("/auth/:provider/callback", goth.NewCompleteAuthHandler(gothConfig))
		root.Get("/logout", goth.NewLogoutHandler(gothConfig))

		sites := app.Domain(":site." + s.cfg.Flags.Domain)
		config := sites_middleware.Config{
			Storage: storage,
		}

		sites.Use(goth.Session(gothConfig))
		sites.Use(goth.Protect(gothConfig))
		sites.Use(sites_middleware.New(config))

		api := root.Group("/api")
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
		spec.UseMiddleware(oidc.NewAuthMiddleware(spec, configs.DefaultConfig.Flags.OIDCIssuer, configs.DefaultConfig.Flags.OIDCAudience))
		sitesHandler.Register(spec)

		err = app.Listen(s.cfg.Flags.Addr)
		if err != nil {
			return err
		}

		return nil
	}
}
