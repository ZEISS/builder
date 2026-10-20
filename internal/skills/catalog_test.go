package skills_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeiss/builder/internal/skills"
)

func TestNewCatalog(t *testing.T) {
	catalog := skills.NewCatalog()
	require.NotNil(t, catalog)
}
