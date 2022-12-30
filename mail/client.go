package mail

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"

	"github.com/inkclip/backend/config"
)

type MailClient struct {
	config config.Config
}

func NewMailClient(config config.Config) (*MailClient, error) {
	server := &MailClient{
		config: config,
	}

	return server, nil
}

type SendContent struct {
	Recipient string
	Subject   string
	Body      string
}

func (client *MailClient) vertifyMailContent(recipient string, token string) SendContent {
	return SendContent{
		Recipient: recipient,
		Subject:   "Verify your email address",
		Body:      fmt.Sprintf("Please click the following link to verify your email address: https://%s/verify?token=%s&email=%s", client.config.FrontURL, token, recipient),
	}
}

func (client *MailClient) send(content *SendContent) error {
	from := "noreply@inkclip.app"
	recipients := []string{content.Recipient}
	subject := content.Subject
	body := content.Body
	var auth smtp.Auth
	if client.config.MailUsername != "" && client.config.MailPassword != "" {
		auth = smtp.CRAMMD5Auth(client.config.MailUsername, client.config.MailPassword)
	} else {
		auth = nil
	}

	msg := []byte(strings.ReplaceAll(fmt.Sprintf("To: %s\nSubject: %s\n\n%s", strings.Join(recipients, ","), subject, body), "\n", "\r\n"))
	err := smtp.SendMail(fmt.Sprintf("%s:%d", client.config.MailHostname, client.config.MailPort), auth, from, recipients, msg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	return nil
}
