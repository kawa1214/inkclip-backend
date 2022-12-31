package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig("..")

	require.NotEmpty(t, config)
	require.NoError(t, err)
}
