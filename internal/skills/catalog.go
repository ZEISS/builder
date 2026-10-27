package skills

import "sync"

// Catalog represents a catalog of skills.
type Catalog interface {
	// Load loads the catalog from the configured roots.
	Load() error
	// Skills returns the list of skills in the catalog.
	Skills() []Skill
	// Skill returns the skill with the given name, or nil if it does not exist.
	Skill(name string) (Skill, bool)
	// Enable is enabling a skill in the catalog.
	Enable(name string)
	// Disable is disabling a skill in the catalog.
	Disable(name string)
	// Clear clears the catalog of skills.
	Clear()
}

var _ Catalog = (*CatalogImpl)(nil)

// CatalogImpl represents a catalog of skills.
type CatalogImpl struct {
	skills       sync.Map
	activeSkills sync.Map
	opts         *CatalogImplOpts
}

// CatalogOpt represents options for creating a Catalog.
type CatalogOpt func(*CatalogImplOpts)

// CatalogOpts is a slice of CatalogOpt functions.
type CatalogImplOpts struct {
	// Roots is the list of root directories to scan for skills.
	Roots []string
}

// DefaultCatalogOpts returns the default options for creating a Catalog.
func DefaultCatalogOpts() *CatalogImplOpts {
	return &CatalogImplOpts{}
}

// WithRoots returns a CatalogOpt that sets the roots of the catalog.
func WithRoots(roots ...string) CatalogOpt {
	return func(opts *CatalogImplOpts) {
		opts.Roots = roots
	}
}

// NewCatalog creates a new Catalog.
func NewCatalog(opts ...CatalogOpt) *CatalogImpl {
	options := DefaultCatalogOpts()

	for _, opt := range opts {
		opt(options)
	}

	catalog := new(CatalogImpl)
	catalog.opts = options

	return catalog
}

// Load loads the catalog from the configured roots.
func (c *CatalogImpl) Load() error {
	for _, root := range c.opts.Roots {
		skills, err := LoadSkills(root)
		if err != nil {
			return err
		}

		for _, skill := range skills {
			c.skills.Store(skill.Name, skill)
		}
	}

	return nil
}

// Skills returns the list of skills in the catalog.
func (c *CatalogImpl) Skills() []Skill {
	var skills []Skill

	c.skills.Range(func(_, skill any) bool {
		skills = append(skills, skill.(Skill))
		return true
	})

	return skills
}

// Enable is enabling a skill in the catalog.
func (c *CatalogImpl) Enable(name string) {
	s, ok := c.skills.Load(name)
	if !ok {
		return
	}
	c.activeSkills.Store(name, s)
}

// Disable is disabling a skill in the catalog.
// It removes the skill from the active skills map.
func (c *CatalogImpl) Disable(name string) {
	c.activeSkills.Delete(name)
}

// Skill returns the skill with the given name, or nil if it does not exist.
func (c *CatalogImpl) Skill(name string) (Skill, bool) {
	s, ok := c.skills.Load(name)
	return s.(Skill), ok
}

// Clear clears the catalog of skills.
func (c *CatalogImpl) Clear() {
	c.skills.Clear()
}
