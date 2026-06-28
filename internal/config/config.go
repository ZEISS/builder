package config

import (
	"os"
	"os/user"
	"path/filepath"
	"sync"

	"github.com/katallaxie/pkg/filex"
	"github.com/zeiss/builder/pkg/specs"
	// Sqlite driver based on CGO.
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
	URL    string
	Flags  *Flags
	Store  string
	Stdin  *os.File
	Stdout *os.File
	Stderr *os.File
	Spec   *specs.Spec
	File   string
	Path   string
	Plugin string
	sync.RWMutex
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
	Dex             bool   `envconfig:"DEX" default:"true"`
	DexClientID     string `envconfig:"DEX_CLIENT_ID"`
	DexClientSecret string `envconfig:"DEX_CLIENT_SECRET"`
	DexClientURL    string `envconfig:"DEX_CLIENT_URL"`
}

// New returns a new config.
func New() *Config {
	return &Config{
		File:   ".builder.yml",
		Store:  "~/.builder/session.db",
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Flags:  &Flags{},
		Spec:   specs.Default(),
	}
}

// Cwd returns the current working directory.
func (c *Config) Cwd() (string, error) {
	return os.Getwd()
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

// LoadSpec is a helper to load the spec from the config file.
func (c *Config) LoadSpec() error {
	f, err := os.ReadFile(filepath.Clean(c.File))
	if err != nil {
		return err
	}

	return c.Spec.UnmarshalYAML(f)
}
