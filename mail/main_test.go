package mail

import (
	"testing"

	"github.com/inkclip/backend/config"
)

func newMailClient(t *testing.T) Client {
	// config, err := config.LoadConfig("..")
	config := config.Config{
		MailHostname: "localhost",
		MailPort:     1025,
		MailUsername: "",
		MailPassword: "",
	}

	client := NewMailClient(config)

	return client
}