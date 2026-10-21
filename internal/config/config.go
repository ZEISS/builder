package config

import (
	"cmp"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/zeiss/builder/pkg/specs"

	"github.com/zeiss/pkg/utilx"
)

// DefaultFlags holds the default command line flags.
var DefaultFlags = NewFlags()

// DefaultConfig holds the configuration for the builder.
var DefaultConfig = NewConfig()

const appName = "sparkles"

// DefaultDataDir is the default data directory for the application.
const DefaultDataDir = "sparkles.db"

// DefaultClientID is the default client ID using the PKCE flow.
const DefaultClientID = "builder-cli"

// DefaultSpecFileName is the default spec file name.
const DefaultSpecFileName = ".builder.yml"

// DefaultPath is the default path to store workspace data.
const DefaultPath = ".sparkles"

// UnknownCwd is the unknown current working directory.
const UnknownCwd = ""

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
	return Flags{
		Config: DefaultSpecFileName,
		Path:   DefaultPath,
	}
}

// Config contains the configuration.
type Config struct {
	// Stderr is the stderr file.
	Stderr *os.File
	// Stdin is the stdin file.
	Stdin *os.File
	// Stdout is the stdout file.
	Stdout *os.File
	// Task contains the flags for a task.
	Task TaskFlags
}

// TaskFlags contains the flags for a task.
type TaskFlags struct {
	// Name is the name of the task to execute.
	Name string
}

// New returns a new config.
func NewConfig() Config {
	return Config{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// Spec is a helper to load the spec from the config file.
func (f *Flags) Spec() (*specs.Spec, error) {
	path, err := f.Abs()
	if err != nil {
		return nil, err
	}

	s, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	spec := &specs.Spec{}
	if err := spec.UnmarshalYAML(s); err != nil {
		return nil, err
	}

	return spec, nil
}

// Abs returns the absolute path of the config file.
func (f *Flags) Abs() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	path := filepath.Join(cwd, f.Config)
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return filepath.Clean(abs), nil
}

// AppDirPath returns the path to the main data directory for the application.
func AppDirPath() (string, error) {
	if xdgDataHome := os.Getenv("XDG_DATA_HOME"); utilx.NotEmpty(xdgDataHome) {
		return filepath.Join(xdgDataHome, appName), nil
	}

	if runtime.GOOS == "windows" {
		localAppData := cmp.Or(
			os.Getenv("LOCALAPPDATA"),
			filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local"),
		)
		return filepath.Join(localAppData, appName), nil
	}

	home, err := HomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".local", "share", appName), nil
}

// DataDir returns the path to the data directory for the application.
func DataDir() (string, error) {
	path, err := AppDirPath()
	if err != nil {
		return "", err
	}

	dataDir := filepath.Clean(filepath.Join(path, DefaultDataDir))

	return dataDir, nil
}

// HomeDir returns the home directory.
func HomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return usr.HomeDir, err
}

// Cwd returns the current working directory.
func Cwd() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Clean(wd), nil
}
