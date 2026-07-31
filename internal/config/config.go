package config

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/zeiss/builder/pkg/specs"
	"github.com/zeiss/pkg/filex"
	// Sqlite driver based on CGO.
)

// DefaultClientID is the default client ID using the PKCE flow.
const DefaultClientID = "builder-cli"

const (
	// UnknownCwd is the unknown current working directory.
	UnknownCwd = ""
)

// DefaultConfig is the default configuration.
var DefaultConfig = New()

// Flags contains the command line flags.
type Flags struct {
	// Plugins contains the plugins to use.
	Plugins []string
	// Vars contains the variables to use.
	Vars []string
	// Dry indicates whether to print dry run messages.
	Dry bool
	// Force indicates whether to force overwrite.
	Force bool
	// Root is the root directory of the project.
	Root bool
	// Verbose indicates whether to print verbose messages.
	Verbose bool
	// Version indicates whether to print version.
	Version bool
	// TaskFlags contains the flags for a task.
	TaskFlags TaskFlags
	// AuthFlags contains the flags for the authentication.
	AuthFlags AuthFlags
}

// NewFlags returns a new flags.
func NewFlags() Flags {
	return Flags{}
}

// Config contains the configuration.
type Config struct {
	URL      string
	Flags    Flags
	Store    string
	Stdin    *os.File
	Stdout   *os.File
	Stderr   *os.File
	Spec     *specs.Spec
	File     string
	Path     string
	Plugin   string
	FileMode os.FileMode
	Verbose  bool
	Task     TaskFlags
}

// TaskFlags contains the flags for a task.
type TaskFlags struct {
	// Name is the name of the task to execute.
	Name string
}

// AuthFlags contains the flags for the authentication.
type AuthFlags struct {
	// ClientID is the client ID for the OIDC provider.
	ClientID string `envconfig:"CLIENT_ID" default:"builder-cli"`
	// ClientURL is the URL of the OIDC provider.
	ClientURL string `envconfig:"CLIENT_URL" default:"http://builder.internal:5556/dex"`
}

// New returns a new config.
func New() Config {
	return Config{
		File:   ".builder.yml",
		Store:  "~/.builder/session.db",
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Flags:  Flags{},
		Spec:   specs.Default(),
	}
}

// Vars returns the variables.
func (c *Config) Vars() []string {
	return c.Flags.Vars
}

// InitDefaultConfig initializes the default configuration.
func (c *Config) InitDefaultConfig() error {
	folder, err := filex.ExpandHomeFolder(c.File)
	if err != nil {
		return err
	}

	c.File = folder

	return nil
}

// HomeDir returns the home directory.
func (c *Config) HomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return usr.HomeDir, err
}

// Cwd return the current working directory.
// Which is the path in which the builder spec is located.
func (c *Config) Cwd() (string, error) {
	abs, err := filepath.Abs(filepath.Clean(c.File))
	if err != nil {
		return UnknownCwd, err
	}

	return filepath.Dir(abs), nil
}

// LoadSpec is a helper to load the spec from the config file.
func (c *Config) LoadSpec() error {
	f, err := os.ReadFile(filepath.Clean(c.File))
	if err != nil {
		return err
	}

	return c.Spec.UnmarshalYAML(f)
}
