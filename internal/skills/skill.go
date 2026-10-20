package skills

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"

	"go.yaml.in/yaml/v2"
)

const (
	defaultSkillsRootPath   = "skills"
	skillFileName           = "SKILL.md"
	yamlFrontmatterDelim    = "---"
	yamlFrontmatterSections = 3
	toolNameRead            = "read"
	toolNameWrite           = "write"
	toolNameEdit            = "edit"
	toolNameBash            = "bash"
)

var (
	// ErrMissingYAMLFrontmatter is returned when the yaml frontmatter is missing.
	ErrMissingYAMLFrontmatter = fmt.Errorf("agents: missing yaml frontmatter")
	// ErrInvalidSkillFormat is returned when the skill format is invalid.
	ErrInvalidSkillFormat = fmt.Errorf("agents: invalid skill format")
)

// Skill represents a skill parsed from a SKILL.md file.
type Skill interface {
	// RootPath returns the root path of the skill.
	RootPath() string
	// Name returns the name of the skill.
	Name() string
	// Description returns the description of the skill.
	Description() string
	// License returns the license of the skill.
	License() string
	// Compatibility returns the compatibility of the skill.
	Compatibility() string
	// Metadata returns the metadata of the skill.
	Metadata() map[string]any
	// AllowedTools returns the allowed tools of the skill.
	AllowedTools() string
	// Body returns the body of the skill.
	Body() string

	yaml.Unmarshaler
}

// SkillImpl represents a skill parsed from a SKILL.md file.
type SkillImpl struct {
	rootPath      string
	name          string
	description   string
	license       string
	compatibility string
	metadata      map[string]any
	allowedTools  string
	body          string
}

// UnmarshalYAML unmarshals the skill from YAML.
func (s *SkillImpl) UnmarshalYAML(unmarshal func(interface{}) error) error {
	skill := struct {
		Name          string         `yaml:"name"`
		Description   string         `yaml:"description"`
		License       string         `yaml:"license"`
		Compatibility string         `yaml:"compatibility"`
		Metadata      map[string]any `yaml:"metadata"`
		AllowedTools  string         `yaml:"allowed-tools"`
		Body          string         `yaml:"body"`
	}{}

	if err := unmarshal(&skill); err != nil {
		return err
	}

	s.name = skill.Name
	s.description = skill.Description
	s.license = skill.License
	s.compatibility = skill.Compatibility
	s.metadata = skill.Metadata
	s.allowedTools = skill.AllowedTools
	s.body = skill.Body

	return nil
}

// RootPath returns the root path of the skill.
func (s *SkillImpl) RootPath() string {
	return s.rootPath
}

// Name returns the name of the skill.
func (s *SkillImpl) Name() string {
	return s.name
}

// Description returns the description of the skill.
func (s *SkillImpl) Description() string {
	return s.description
}

// License returns the license of the skill.
func (s *SkillImpl) License() string {
	return s.license
}

// Compatibility returns the compatibility of the skill.
func (s *SkillImpl) Compatibility() string {
	return s.compatibility
}

// Metadata returns the metadata of the skill.
func (s *SkillImpl) Metadata() map[string]any {
	return s.metadata
}

// AllowedTools returns the allowed tools of the skill.
func (s *SkillImpl) AllowedTools() string {
	return s.allowedTools
}

// Body returns the body of the skill.
func (s *SkillImpl) Body() string {
	return s.body
}

// Parser respresents a parser that can parse SKILL.md files.
type Parser interface {
	// Parse parses a SKILL.md file and returns a Skill.
	Parse(r io.Reader) (Skill, error)
}

var _ Parser = (*ParserImpl)(nil)

// ParserImpl is an implementation of the Parser interface.
type ParserImpl struct{}

// NewParser returns a new Parser.
func NewParser() *ParserImpl {
	return &ParserImpl{}
}

// Parse parses a skill definition from an io.Reader.
func (p *ParserImpl) Parse(r io.Reader) (Skill, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	content := string(data)

	content = strings.TrimSpace(content)
	if !strings.HasPrefix(content, yamlFrontmatterDelim) {
		return nil, ErrMissingYAMLFrontmatter
	}

	parts := strings.SplitN(content, yamlFrontmatterDelim, yamlFrontmatterSections)
	if len(parts) < yamlFrontmatterSections {
		return nil, ErrInvalidSkillFormat
	}

	yamlPart := parts[1]
	body := strings.TrimSpace(parts[2])

	skill := &SkillImpl{}
	if unmarshalErr := yaml.Unmarshal([]byte(yamlPart), &skill); unmarshalErr != nil {
		return nil, unmarshalErr
	}

	skill.body = body

	return skill, nil
}

// LoadSkills loads a skill by seeking skills files in the given path.
func LoadSkills(path string) ([]Skill, error) {
	skills := make([]Skill, 0)
	parser := NewParser()

	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if !strings.Contains(d.Name(), skillFileName) {
			return nil
		}

		skill, errr := parser.Parse(strings.NewReader(d.Name()))
		if errr != nil {
			return errr
		}
		skills = append(skills, skill)

		return nil
	})

	return skills, err
}
