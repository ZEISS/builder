package config

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/zeiss/builder/pkg/specs"
	"github.com/zeiss/pkg/filex"
)

// DefaultClientID is the default client ID using the PKCE flow.
const DefaultClientID = "builder-cli"

// DefaultSpecFileName is the default spec file name.
const DefaultSpecFileName = ".builder.yaml"

// DefaultPath is the default path to store workspace data.
const DefaultPath = ".sparkles"

// UnknownCwd is the unknown current working directory.
const UnknownCwd = ""

// DefaultConfig is the default configuration.
var DefaultConfig = New()

// Flags contains the command line flags.
type Flags struct {
	// URL is the URL of the builder server.
	URL string
	// Plugins contains the plugins to use.
	Plugins []string
	// Vars contains the variables to use.
	Vars []string
	// Force indicates whether to force overwrite.
	Force bool
	// Verbose indicates whether to print verbose messages.
	Verbose bool
	// Debug indicates whether to enable debug mode.
	Debug bool
	// Yolo indicates whether to enable yolo mode.
	Yolo bool
	// TaskFlags contains the flags for a task.
	TaskFlags TaskFlags
	// Config is the path to the config file.
	Config string
	// Path is the path to the store file.
	Path string
	// Session is the session ID to use.
	Session string
	// Continue indicates whether to continue a session.
	Continue bool
}

// NewFlags returns a new flags.
func NewFlags() Flags {
	return Flags{}
}

// NewDefaultFlags returns the default flags.
func NewDefaultFlags() Flags {
	return Flags{
		Config: DefaultSpecFileName,
		Path:   DefaultPath,
	}
}

// Config contains the configuration.
type Config struct {
	File     string
	FileMode os.FileMode
	Flags    Flags
	Stderr   *os.File
	Stdin    *os.File
	Stdout   *os.File
	Task     TaskFlags
	Verbose  bool
}

// TaskFlags contains the flags for a task.
type TaskFlags struct {
	// Name is the name of the task to execute.
	Name string
}

// New returns a new config.
func New() Config {
	return Config{
		File:   DefaultSpecFileName,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Flags:  NewDefaultFlags(),
	}
}

// Vars returns the variables.
func (c *Config) Vars() []string {
	return c.Flags.Vars
}

// HomeDir returns the home directory.
func (c *Config) HomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return usr.HomeDir, err
}

// StorePath returns the path to the store file.
func (c *Config) StorePath() (string, error) {
	path, err := filex.ExpandHomeFolder(c.Flags.Path)
	if err != nil {
		return "", err
	}

	abs, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return "", err
	}

	return abs, nil
}

// Cwd return the current working directory.
// Which is the path in which the builder spec is located.
func (c *Config) Cwd() (string, error) {
	abs, err := filepath.Abs(filepath.Clean(c.File))
	if err != nil {
		return UnknownCwd, err
	}

	return filepath.Dir(filepath.Clean(abs)), nil
}

// Spec is a helper to load the spec from the config file.
func (c *Config) Spec() (*specs.Spec, error) {
	path, err := c.Abs()
	if err != nil {
		return nil, err
	}

	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	spec := &specs.Spec{}
	if err := spec.UnmarshalYAML(f); err != nil {
		return nil, err
	}

	return spec, nil
}

// Path returns the expanded path of the config file.
func (c *Config) Path() (string, error) {
	path, err := filex.ExpandHomeFolder(c.Flags.Path)
	if err != nil {
		return "", err
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return abs, nil
}

// DeleteStore deletes the session store.
func (c *Config) DeleteStore() error {
	path, err := filex.ExpandHomeFolder(c.Flags.Path)
	if err != nil {
		return err
	}

	return os.RemoveAll(path)
}

// Abs returns the absolute path of the config file.
func (c *Config) Abs() (string, error) {
	path, err := filex.ExpandHomeFolder(c.File)
	if err != nil {
		return "", err
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return filepath.Clean(abs), nil
}

// ExpandConfigPath expands the config path and creates the store folder.
func ExpandConfigPath(path string) (string, error) {
	path, err := filex.ExpandHomeFolder(path)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(filepath.Dir(path), 0o777)
	if err != nil {
		return "", err
	}

	return path, nil
}
