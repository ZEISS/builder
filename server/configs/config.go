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
	// SqliteFlags contains the SQLite flags.
	SqliteFlags SqliteFlags
}

// FilesFlags contains the flags for the files directory.
type FilesFlags struct {
	// Path is the path to the files directory.
	Path string `envconfig:"BUILDER_FILES_PATH" default:""`
}

// SqliteFlags returns the path to the SQLite database.
type SqliteFlags struct {
	// Enabled is a flag to enable SQLite.
	Enabled bool `envconfig:"BUILDER_SQLITE_ENABLED" default:"true"`
	// Path is the path to the SQLite database.
	Path string `envconfig:"BUILDER_SQLITE_PATH" default:""`
	// Database is the name of the SQLite database.
	Database string `envconfig:"BUILDER_SQLITE_DATABASE" default:""`
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
		SqliteFlags: SqliteFlags{
			Enabled:  true,
			Path:     "builder.db",
			Database: "builder",
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
