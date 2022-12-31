package mail

import (
	"testing"

	"github.com/inkclip/backend/util"
	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	client := newMailClient(t)
	arg := SendContent{
		Recipient: util.RandomEmail(),
		Subject:   util.RandomString(10),
		Body:      util.RandomString(10),
	}
	err := client.Send(arg)
	require.NoError(t, err)
}

func TestVertifyMailContent(t *testing.T) {
	client := newMailClient(t)
	recipient := util.RandomEmail()
	token := util.RandomString(10)
	arg := client.VertifyMailContent(recipient, token)

	err := client.Send(arg)
	require.NoError(t, err)
}
