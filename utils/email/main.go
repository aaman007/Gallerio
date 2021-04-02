package email

import (
	"context"
	"fmt"
	"gopkg.in/mailgun/mailgun-go.v4"
	"log"
	"net/url"
	"time"
)

var (
	baseResetURL = "http://localhost:8000/reset"
	
	welcomeSubject = "Welcome to Gallerio"
	welcomeText    = "Greeting. Its a pleasure to have you here. Cheers"
	welcomeHtml    = `
	Hello<br/>
	<br/>
	We are pleased to have you here. Visit our site: <a href="gallerio.com">Gallerio</a>"<br/>
	Have a great day`
	
	resetPasswordSubject = "Reset your password"
	resetPasswordText    = `
	It appears you have requested to reset your password.
	Use the following link to reset your password
	%s
	You can also use the code below
	%s
	If you didn't requested this, then ignore this message`
	resetPasswordHtml = `
	It appears you have requested to reset your password.<br/>
	Use the following link to reset your password<br/>
	<a href="%s">"%s</a><br/>
	You can also use the code below<br/>
	%s<br/>
	If you didn't requested this, then ignore this message<br/>`
)

type ClientConfig func(*Client)

func WithMailgun(domain, apiKey string) ClientConfig {
	return func(client *Client) {
		client.mg = mailgun.NewMailgun(domain, apiKey)
	}
}

func WithSender(name, email string) ClientConfig {
	return func(client *Client) {
		client.from = buildEmail(name, email)
	}
}

func NewClient(opts ...ClientConfig) Client {
	client := Client{
		from: "support@gallerio.com",
	}
	for _, opt := range opts {
		opt(&client)
	}
	return client
}

type Client struct {
	from string
	mg   mailgun.Mailgun
}

func (c *Client) Welcome(name, email string) error {
	message := c.mg.NewMessage(c.from, welcomeSubject, welcomeText, buildEmail(name, email))
	message.SetHtml(welcomeHtml)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	
	_, _, err := c.mg.Send(ctx, message)
	log.Println(err)
	return err
}

func (c *Client) ResetPassword(email, token string) error {
	v := url.Values{}
	v.Set("token", token)
	resetUrl := baseResetURL + "?" + v.Encode()
	resetPasswordText = fmt.Sprintf(resetPasswordText, resetUrl, token)
	message := c.mg.NewMessage(c.from, resetPasswordSubject, resetPasswordText, email)
	resetPasswordHtml = fmt.Sprintf(resetPasswordHtml, resetUrl, resetUrl, token)
	message.SetHtml(resetPasswordHtml)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	
	_, _, err := c.mg.Send(ctx, message)
	return err
}

func buildEmail(name, email string) string {
	if name == "" {
		return email
	}
	return fmt.Sprintf("%s <%s>", name, email)
}
