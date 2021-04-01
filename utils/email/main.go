package email

import (
	"context"
	"fmt"
	"gopkg.in/mailgun/mailgun-go.v4"
	"log"
	"time"
)

const (
	welcomeSubject = "Welcome to Gallerio"
	welcomeText = "Greeting. Its a pleasure to have you here. Cheers"
	welcomeHtml = `
	Hello<br/>
	<br/>
	We are pleased to have you here. Visit our site: <a href="gallerio.com">Gallerio</a>"<br/>
	Have a great day`
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
	mg mailgun.Mailgun
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

func buildEmail(name, email string) string {
	if name == "" {
		return email
	}
	return fmt.Sprintf("%s <%s>", name, email)
}