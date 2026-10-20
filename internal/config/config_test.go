package config_test

import (
	"testing"

	"github.com/zeiss/builder/internal/config"

	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	t.Parallel()

	cfg := config.NewConfig()

	require.NotNil(t, cfg)
}

func TestNewFlags(t *testing.T) {
	t.Parallel()

	flags := config.NewFlags()

	require.NotNil(t, flags)
}
