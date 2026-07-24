package configs

import (
	"os"
)

// Flags contains the command line flags.
type Flags struct {
	Addr            string `envconfig:"BUILDER_ADDR" default:":3000"`
	Domain          string `envconfig:"BUILDER_DOMAIN" default:""`
	DexClientID     string `envconfig:"BUILDER_DEX_CLIENT_ID" default:""`
	DexClientSecret string `envconfig:"BUILDER_DEX_CLIENT_SECRET" default:""`
	DexCallbackURL  string `envconfig:"BUILDER_DEX_CALLBACK_URL" default:""`
	DexLoginURL     string `envconfig:"BUILDER_DEX_LOGIN_URL" default:""`
	OIDCIssuer      string `envconfig:"BUILDER_OIDC_ISSUER" default:""`
	OIDCAudience    string `envconfig:"BUILDER_OIDC_AUDIENCE" default:""`
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
	Enabled bool `envconfig:"BUILDER_SQLITE_ENABLED" default:""`
	// Path is the path to the SQLite database.
	Path string `envconfig:"BUILDER_SQLITE_PATH" default:""`
	// Database is the name of the SQLite database.
	Database string `envconfig:"BUILDER_SQLITE_DATABASE" default:""`
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
