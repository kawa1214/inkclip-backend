package mail

import (
	"testing"

	"github.com/inkclip/backend/config"
	"github.com/stretchr/testify/require"
)

func newMailClient(t *testing.T) Client {
	config, err := config.LoadConfig("..")
	require.NoError(t, err)

	client := NewMailClient(config)

	return client
}
