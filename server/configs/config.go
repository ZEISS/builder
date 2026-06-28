package configs

import (
	"os"
)

// Flags contains the command line flags.
type Flags struct {
	Addr            string `envconfig:"BUILDER_ADDR" default:":3000"`
	Root            string `envconfig:"BUILDER_ROOT" default:""`
	Domain          string `envconfig:"BUILDER_DOMAIN" default:""`
	Database        string `envconfig:"BUILDER_DATABASE" default:"example"`
	DatabaseUser    string `envconfig:"BUILDER_DATABASE_USER" default:"example"`
	DatabasePass    string `envconfig:"BUILDER_DATABASE_PASS" default:"example"`
	DatabaseHost    string `envconfig:"BUILDER_DATABASE_HOST" default:"localhost"`
	DatabasePort    int    `envconfig:"BUILDER_DATABASE_PORT" default:"5432"`
	DexClientID     string `envconfig:"BUILDER_DEX_CLIENT_ID" default:""`
	DexClientSecret string `envconfig:"BUILDER_DEX_CLIENT_SECRET" default:""`
	DexCallbackURL  string `envconfig:"BUILDER_DEX_CALLBACK_URL" default:""`
	DexLoginURL     string `envconfig:"BUILDER_DEX_LOGIN_URL" default:""`
	OIDCIssuer      string `envconfig:"BUILDER_OIDC_ISSUER" default:""`
	OIDCAudience    string `envconfig:"BUILDER_OIDC_AUDIENCE" default:""`
}

// NewFlags returns a new instance of Flags.
func NewFlags() *Flags {
	return &Flags{}
}

// New returns a new instance of Config.
func New() *Config {
	return &Config{
		Flags: NewFlags(),
	}
}

// Config contains the configuration.
type Config struct {
	Flags *Flags
}

// Cwd returns the current working directory.
func (c *Config) Cwd() (string, error) {
	return os.Getwd()
}
