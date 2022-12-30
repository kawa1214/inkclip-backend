package mail

import (
	"log"
	"testing"

	"github.com/inkclip/backend/config"
	"github.com/stretchr/testify/require"
)

func newMailClient(t *testing.T) *MailClient {
	config, err := config.LoadConfig("..")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	client, err := NewMailClient(config)
	require.NoError(t, err)

	return client
}
