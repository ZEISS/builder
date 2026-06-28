package main

//go:generate go tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config config.client.yml api.yml
//go:generate go tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config config.models.yml api.yml

import (
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/gofiber/fiber/v3"
	"github.com/spf13/cobra"
	"github.com/zeiss/builder/server/adapters/handlers"
	"github.com/zeiss/builder/server/adapters/store"
	"github.com/zeiss/builder/server/controllers"
)

// Options for the CLI.
type Options struct {
	Port int `help:"Port to listen on" short:"p" default:"8888"`
}

func main() {
	var api huma.API

	// Create a CLI app which takes a port option.
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		app := fiber.New()

		siteController := controllers.NewSitesController(store.NewFS(""))
		sitesHandler := handlers.NewSitesHandler(siteController)

		apiConfig := huma.DefaultConfig("Builder API", "1.0.0")
		apiConfig.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
			"openid": {
				Type:             "openIdConnect",
				In:               "header",
				BearerFormat:     "JWT",
				OpenIDConnectURL: "",
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

		api = humafiber.New(app, apiConfig)
		sitesHandler.Register(api)

		// Tell the CLI how to start your server.
		hooks.OnStart(func() {
			fmt.Printf("Starting server on port %d...\n", options.Port)
			app.Listen(fmt.Sprintf(":%d", options.Port))
		})
	})

	// Add a command to print the OpenAPI spec.
	cli.Root().AddCommand(&cobra.Command{
		Use:   "openapi",
		Short: "Print the OpenAPI spec",
		Run: func(cmd *cobra.Command, args []string) {
			// Use downgrade to return OpenAPI 3.0.3 YAML since oapi-codegen doesn't
			// support OpenAPI 3.1 fully yet. Use `.YAML()` instead for 3.1.
			b, _ := api.OpenAPI().DowngradeYAML()
			fmt.Println(string(b))
		},
	})

	// Run the CLI. When passed no commands, it starts the server.
	cli.Run()
}
