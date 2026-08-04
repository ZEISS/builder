package configs

import (
	"os"
)

var DefaultConfig = New()

// Flags contains the command line flags.
type Flags struct {
	// Addr is the address to listen on.
	Addr string `envconfig:"BUILDER_ADDR" default:":3000"`
	// Domain is the domain name for the builder.
	Domain string `envconfig:"BUILDER_DOMAIN" default:""`
	// OIDCIssuer is the OIDC issuer URL.
	OIDCIssuer string `envconfig:"BUILDER_OIDC_ISSUER" default:""`
	// OIDCAudience is the OIDC audience.
	OIDCAudience string `envconfig:"BUILDER_OIDC_AUDIENCE" default:""`
	// DexFlags contains the flags for the Dex authentication provider.
	DexFlags DexFlags
	// FilesFlags contains the flags for the files directory.
	FilesFlags FilesFlags
	// PostgresFlags contains the flags for the PostgreSQL database.
	PostgresFlags PostgresFlags
}

// FilesFlags contains the flags for the files directory.
type FilesFlags struct {
	// Path is the path to the files directory.
	Path string `envconfig:"BUILDER_FILES_PATH" default:""`
}

// PostgresFlags contains the flags for the PostgreSQL database.
type PostgresFlags struct {
	// Enabled is a flag to enable PostgreSQL.
	Enabled bool `envconfig:"BUILDER_POSTGRES_ENABLED" default:"true"`
	// Host is the host of the PostgreSQL database.
	Host string `envconfig:"BUILDER_POSTGRES_HOST" default:""`
	// Port is the port of the PostgreSQL database.
	Port int `envconfig:"BUILDER_POSTGRES_PORT" default:""`
	// Database is the name of the PostgreSQL database.
	Database string `envconfig:"BUILDER_POSTGRES_DATABASE" default:""`
	// User is the username for the PostgreSQL database.
	User string `envconfig:"BUILDER_POSTGRES_USER" default:""`
	// Password is the password for the PostgreSQL database.
	Password string `envconfig:"BUILDER_POSTGRES_PASSWORD" default:""`
}

// DexFlags contains the flags for the Dex authentication provider.
type DexFlags struct {
	// Enabled is a flag to enable the Dex authentication provider.
	Enabled bool `envconfig:"BUILDER_DEX_ENABLED" default:"true"`
	// URL is the URL of the Dex authentication provider.
	URL string `envconfig:"BUILDER_DEX_URL" default:""`
	// ClientID is the client ID for the Dex authentication provider.
	ClientID string `envconfig:"BUILDER_DEX_CLIENT_ID" default:""`
	// ClientSecret is the client secret for the Dex authentication provider.
	ClientSecret string `envconfig:"BUILDER_DEX_CLIENT_SECRET" default:""`
	// CallbackURL is the callback URL for the Dex authentication provider.
	CallbackURL string `envconfig:"BUILDER_DEX_CALLBACK_URL" default:""`
	// LoginURL is the login URL for the Dex authentication provider.
	LoginURL string `envconfig:"BUILDER_DEX_LOGIN_URL" default:""`
}

// NewFlags returns a new instance of Flags.
func NewFlags() *Flags {
	return &Flags{
		PostgresFlags: PostgresFlags{
			Enabled:  true,
			Host:     "localhost",
			Port:     5432,
			Database: "default",
		},
		FilesFlags: FilesFlags{
			Path: "files",
		},
	}
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
