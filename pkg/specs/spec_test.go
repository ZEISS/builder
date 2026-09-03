package specs_test

import (
	"testing"

	"github.com/zeiss/builder/pkg/specs"

	"github.com/go-openapi/testify/v2/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	s := specs.New()
	require.NotNil(t, s)
	require.Equal(t, specs.DefaultVersion, s.Version)
}

func TestDefault(t *testing.T) {
	t.Parallel()

	s := specs.Default()
	require.NotNil(t, s)
	require.Equal(t, specs.DefaultVersion, s.Version)
}

func TestExample(t *testing.T) {
	t.Parallel()

	s, err := specs.Example()
	require.NoError(t, err)
	require.NotNil(t, s)
	require.Equal(t, 1, s.Version)
	require.Equal(t, "example", s.Name)
}
