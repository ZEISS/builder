package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/zeiss/builder/server/adapters/handlers"
	"github.com/zeiss/builder/server/adapters/store"
	"github.com/zeiss/builder/server/configs"
	"github.com/zeiss/builder/server/controllers"
	"github.com/zeiss/builder/server/middlewares/auth/oidc"
	"github.com/zeiss/builder/server/middlewares/static"
	"gorm.io/gorm"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v3"
	logger "github.com/gofiber/fiber/v3/middleware/logger"
	requestid "github.com/gofiber/fiber/v3/middleware/requestid"
	goth "github.com/katallaxie/fiber-goth/v3"
	gorm_adapter "github.com/katallaxie/fiber-goth/v3/adapters/gorm"
	"github.com/katallaxie/fiber-goth/v3/providers"
	"github.com/katallaxie/fiber-goth/v3/providers/dex"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
	"github.com/zeiss/pkg/server"
	"gorm.io/driver/postgres"
)

var cfg *configs.Config

func init() {
	cfg = configs.New()

	err := envconfig.Process("", cfg.Flags)
	if err != nil {
		log.Fatal(err)
	}

	Root.PersistentFlags().StringVar(&cfg.Flags.Addr, "addr", cfg.Flags.Addr, "addr")
	Root.PersistentFlags().StringVar(&cfg.Flags.OIDCIssuer, "oidc-issuer", cfg.Flags.OIDCIssuer, "OIDC Issuer")
	Root.PersistentFlags().StringVar(&cfg.Flags.OIDCAudience, "oidc-audience", cfg.Flags.OIDCAudience, "OIDC Audience")
	Root.PersistentFlags().StringVar(&cfg.Flags.Root, "root", cfg.Flags.Root, "root")
	Root.PersistentFlags().StringVar(&cfg.Flags.Domain, "domain", cfg.Flags.Domain, "domain")

	Root.PersistentFlags().StringVar(&cfg.Flags.DexClientID, "dex-client-id", cfg.Flags.DexClientID, "Dex Client ID")
	Root.PersistentFlags().StringVar(&cfg.Flags.DexClientSecret, "dex-client-secret", cfg.Flags.DexClientSecret, "Dex Client Secret")
	Root.PersistentFlags().StringVar(&cfg.Flags.DexCallbackURL, "dex-callback-url", cfg.Flags.DexCallbackURL, "Dex Callback URL")
	Root.PersistentFlags().StringVar(&cfg.Flags.DexLoginURL, "dex-login-url", cfg.Flags.DexLoginURL, "Dex Login URL")

	Root.PersistentFlags().StringVar(&cfg.Flags.DatabaseHost, "db-host", cfg.Flags.DatabaseHost, "DB Host")
	Root.PersistentFlags().StringVar(&cfg.Flags.DatabaseUser, "db-username", cfg.Flags.DatabaseUser, "DB Username")
	Root.PersistentFlags().StringVar(&cfg.Flags.DatabasePass, "db-password", cfg.Flags.DatabasePass, "DB Password")
	Root.PersistentFlags().StringVar(&cfg.Flags.Database, "db-database", cfg.Flags.Database, "DB Database")
	Root.PersistentFlags().IntVar(&cfg.Flags.DatabasePort, "db-port", cfg.Flags.DatabasePort, "DB Port")

	Root.SilenceUsage = true
}

var Root = &cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		srv := NewWebSrv(cfg)

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
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", cfg.Flags.DatabaseHost, cfg.Flags.DatabaseUser, cfg.Flags.DatabasePass, cfg.Flags.Database, cfg.Flags.DatabasePort)
		conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return err
		}

		if err := gorm_adapter.RunMigrations(conn); err != nil {
			return err
		}

		providers.RegisterProvider(dex.New(cfg.Flags.DexClientID, cfg.Flags.DexClientSecret, cfg.Flags.OIDCIssuer, cfg.Flags.DexCallbackURL))

		ga := gorm_adapter.New(conn)

		fs := store.NewFS(s.cfg.Flags.Root)
		sitesCtrl := controllers.NewSitesController(fs)
		sitesHandler := handlers.NewSitesHandler(sitesCtrl)

		c := fiber.Config{}

		app := fiber.New(c)
		app.Use(requestid.New())
		app.Use(logger.New())

		gothConfig := goth.Config{
			Adapter:        ga,
			Secret:         goth.GenerateKey(),
			CookieHTTPOnly: true,
			LoginURL:       cfg.Flags.DexLoginURL,
			CookieDomain:   cfg.Flags.Domain,
		}

		root := app.Domain(cfg.Flags.Domain)
		root.Get("/session", goth.NewSessionHandler(gothConfig))
		root.Get("/login/:provider", goth.NewBeginAuthHandler(gothConfig))
		root.Get("/auth/:provider/callback", goth.NewCompleteAuthHandler(gothConfig))
		root.Get("/logout", goth.NewLogoutHandler(gothConfig))

		sites := app.Domain(":site." + cfg.Flags.Domain)
		config := static.Config{
			Root: http.Dir(cfg.Flags.Root),
		}
		sites.Use(goth.Session(gothConfig))
		sites.Use(goth.Protect(gothConfig))
		sites.Use(static.New(config))

		api := app.Group("/api")
		v1 := api.Group("/v1")

		apiConfig := huma.DefaultConfig("Builder API", "1.0.0")
		apiConfig.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
			"openid": {
				Type:             "openIdConnect",
				In:               "header",
				BearerFormat:     "JWT",
				OpenIDConnectURL: cfg.Flags.OIDCIssuer,
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
		spec.UseMiddleware(oidc.NewAuthMiddleware(spec, cfg.Flags.OIDCIssuer, cfg.Flags.OIDCAudience))
		sitesHandler.Register(spec)

		err = app.Listen(s.cfg.Flags.Addr)
		if err != nil {
			return err
		}

		return nil
	}
}
