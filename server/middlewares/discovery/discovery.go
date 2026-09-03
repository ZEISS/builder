package discovery

import (
	"github.com/gofiber/fiber/v3"
	"github.com/zeiss/pkg/utilx"
)

// WellKnownConfig holds the well-known configuration.
type WellKnownConfig struct {
	// OidcIssuer is the URL to the OIDC issuer.
	//
	// Default: "http://builder.internal:5556/dex"
	OidcIssuer string `json:"oidc_issuer"`
	// ApiURL is the URL to the API.
	//
	// Default: "http://builder.internal:8080/api/v1"
	ApiURL string `json:"api_url"`
}

// WellKnownConfiguratonURL returns the URL of the well-known configuration.
const WellKnownConfigurationURL = "/.well-known/builder-configuration"

// DefaultWellKnownConfig returns the default well-known configuration.
func DefaultWellKnownConfig() WellKnownConfig {
	return WellKnownConfig{
		OidcIssuer: "http://builder.internal:5556/dex",
		ApiURL:     "http://builder.internal:8080/api/v1",
	}
}

// DefaultWellKnownFunc returns the default well-known configuration function.
func DefaultWellKnownFunc() WellKnownConfig {
	return DefaultWellKnownConfig()
}

// Config holds the Discovery configuration.
type Config struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c fiber.Ctx) bool
	// WellKnownFunc returns the well-known configuration.
	//
	// Optional. Default: nil
	WellKnownFunc func() WellKnownConfig
}

// DefaultConfig returns the default Discovery configuration.
func DefaultConfig() Config {
	return Config{
		WellKnownFunc: DefaultWellKnownFunc,
	}
}

// New creates a discovery middleware handler.
func New(config ...Config) fiber.Handler {
	// Set default config
	cfg := DefaultConfig()

	// Override config if provided
	if len(config) > 0 {
		cfg = config[0]

		if utilx.IsNil(cfg.WellKnownFunc) {
			cfg.WellKnownFunc = DefaultWellKnownFunc
		}
	}

	// Return new handler
	return func(c fiber.Ctx) error {
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		return c.JSON(cfg.WellKnownFunc())
	}
}
