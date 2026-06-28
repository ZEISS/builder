package specs

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/katallaxie/pkg/filex"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

//go:embed example.yml
var Files embed.FS

const (
	// Default values for the spec.
	DefaultVersion = 1
	// DefaultName is the default filename for the spec file. It can be changed in the builder configuration.
	DefaultName = ".builder.yml"
	// DefaultFolder is the default folder where the spec file is stored.
	DefaultFolder = ".builder"
)

// Spec defines the structure for a specification document.
// This could be used in various contexts such as API specifications,
// software design documents, etc.
type Spec struct {
	Version      int     `json:"version" yaml:"version"`
	Name         string  `json:"name,omitempty" yaml:"name"`
	Description  string  `json:"description,omitempty" yaml:"description,omitempty"`
	Rules        Rules   `json:"rules" yaml:"rules"`
	Constitution string  `json:"constitution" yaml:"constitution"`
	Models       []Model `json:"models,omitempty" yaml:"models,omitempty"`
	Deploy       *Deploy `json:"domain,omitempty" yaml:"domain,omitempty"`
	Tasks        Tasks   `json:"tasks" yaml:"tasks"`
	Root         string  `json:"root,omitempty" yaml:"root,omitempty"`
}

// Deploy defines the deployment configuration for the specification.
type Deploy struct {
	// Path is the path to the deployment configuration file.
	Path string `json:"path" yaml:"path"`
	// Site is the name of the site to deploy to.
	Site string `json:"site" yaml:"site"`
	// Ignore is a list of files to ignore during deployment.
	Ignore []string `json:"ignore" yaml:"ignore"`
}

// Model specifies a model.
type Model struct {
	RequestOptions map[string]any `json:"requestOptions,omitempty" yaml:"requestOptions,omitempty"`
	ID             string         `json:"id" yaml:"id"`
	Name           string         `json:"name" yaml:"name"`
	URL            string         `json:"url" yaml:"url"`
	Provider       Provider       `json:"provider" yaml:"provider"`
	Roles          []string       `json:"roles" yaml:"roles"`
}

// Provider is a string that specifies the provider for a model.
// For example, openai, azure, etc.
type Provider string

// Tasks defines the tasks for the specification.
type Tasks map[string]Task

// Rules defines the rules for the specification.
type Rules map[string][]Rule

// Rule defines a rule for the specification.
type Rule string

// Task defines a specifiction.
type Task struct {
	// ID is the ID of the task.
	ID string `json:"id" yaml:"id"`
	// Name is the name of the task.
	Name string `json:"name" yaml:"name"`
	// Description is the description of the task.
	Description string `json:"description" yaml:"description"`
	// Context is the context for the task.
	Context Context `json:"context" yaml:"context"`
}

// Documents is the documents for the specification.
type Document struct {
	// ID identifies a document in the specification.
	ID string `json:"id" yaml:"id"`
	// Generates the document from the task's output.
	Generates string `json:"generates" yaml:"generates"`
	// Description is the description of the document.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Instruction is the instruction for generating the document.
	Instruction string `json:"instruction" yaml:"instruction"`
	// Template is the basis for the document.
	Template string `json:"template" yaml:"template"`
}

// Documents is the documents for the specification.
type Documents []Document

// Context is the context for a task.
type Context struct {
	// Requires is the list of requirements for the task.
	Requires []string `json:"requires,omitempty" yaml:"requires,omitempty"`
	// Instruction is the instruction for the task.
	Instruction string `json:"instruction" yaml:"instruction"`
	// TrackBy is the list of fields to track by.
	TrackBy []string `json:"trackBy,omitempty" yaml:"trackBy,omitempty"`
	// Documents is the list of documents for the task.
	Documents Documents `json:"documents" yaml:"documents"`
}

// New returns a new instance of Spec with default values.
func New() *Spec {
	return &Spec{
		Version: DefaultVersion,
	}
}

// Default return a default instance of the spec.
func Default() *Spec {
	return New()
}

// Example returns a default instance of an example.
func Example() (*Spec, error) {
	s := Default()
	data, err := Files.ReadFile("example.yml")
	if err != nil {
		return nil, err
	}

	err = s.UnmarshalYAML(data)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// UnmarshalYAML unmarshals YAML data into a Spec struct.
func (s *Spec) UnmarshalYAML(data []byte) error {
	spec := struct {
		Rules        Rules   `yaml:"rules"`
		Tasks        Tasks   `yaml:"tasks"`
		Constitution string  `yaml:"constitution"`
		Description  string  `yaml:"description"`
		Name         string  `yaml:"name"`
		Models       []Model `yaml:"models"`
		Version      int     `yaml:"version"`
		Deploy       *Deploy `yaml:"deploy"`
		Root         string  `yaml:"root"`
	}{
		Version: DefaultVersion,
		Root:    DefaultFolder,
	}

	if err := yaml.Unmarshal(data, &spec); err != nil {
		return errors.WithStack(err)
	}

	if err := yaml.Unmarshal(data, &spec); err != nil {
		return errors.WithStack(err)
	}

	s.Version = spec.Version
	s.Name = spec.Name
	s.Rules = spec.Rules
	s.Constitution = spec.Constitution
	s.Description = spec.Description
	s.Tasks = spec.Tasks
	s.Deploy = spec.Deploy
	s.Root = spec.Root

	return nil
}

// Write is the write function for the spec.
func Write(s *Spec, file string, force bool) error {
	ok, _ := filex.FileExists(filepath.Clean(file))
	if ok && !force {
		return fmt.Errorf("%s already exists, use --force to overwrite", file)
	}

	f, err := os.Create(filepath.Clean(file))
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	err = yaml.NewEncoder(f).Encode(s)
	if err != nil {
		return err
	}

	return nil
}
