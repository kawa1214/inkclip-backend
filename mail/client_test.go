package mail

import (
	"testing"

	"github.com/inkclip/backend/util"
	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	client := newMailClient(t)
	arg := &SendContent{
		Recipient: "test@test.com",
		Subject:   "Test",
		Body:      "Test",
	}
	err := client.send(arg)
	require.NoError(t, err)
}

func TestVertifyMailContent(t *testing.T) {
	client := newMailClient(t)
	recipient := util.RandomEmail()
	token := "token"
	arg := client.vertifyMailContent(recipient, token)

	err := client.send(&arg)
	require.NoError(t, err)
}
