package email

import (
	"fmt"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

const (
	welcomeSubject = "Welcome to LensLockedBR.com!"
)

const welcomeText = `Hi there!

Welcome to LensLockedBR.com! We really hope you enjoy using
our application!

Best regards,

LH
`

const welcomeHTML = `Hi there!<br/><br/>

<p>Welcome to LensLockedBR.com! We really hope you enjoy using
our application!</p>

<p>Best regards,</p>

<p>LH</p>
`
type Client struct {
	from string
	mg mailgun.Mailgun
}

func (c *Client) Welcome(toName, toEmail string) error {
	message := mailgun.NewMessage(c.from, welcomeSubject, 
                                      welcomeText,
                                      buildEmail(toName, toEmail))
	message.SetHtml(welcomeHTML)

	_, _, err := c.mg.Send(message)
	return err
}

type ClientConfig func(*Client)

func NewClient(opts ...ClientConfig) *Client {
	client := Client {
		from: "support@lenslockedbr.com",
	}

	for _, opt := range opts {
		opt(&client)
	}

	return &client
}

func WithMailgun(domain, apiKey, publicKey string) ClientConfig {
	return func(c *Client) {
		mg := mailgun.NewMailgun(domain, apiKey, publicKey)
		c.mg = mg
	}
}

func WithSender(name, email string) ClientConfig {
	return func(c *Client) {
		c.from = buildEmail(name, email)
	}
}

/////////////////////////////////////////////////////////////////////
//
// Helper Methods
//
/////////////////////////////////////////////////////////////////////

func buildEmail(name, email string) string {
	if name == "" {
		return email
	}

	return fmt.Sprintf("%s <%s>", name, email)
}


