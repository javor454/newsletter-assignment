package sendgrid

import (
	"github.com/javor454/newsletter-assignment/app/config"
	"github.com/sendgrid/sendgrid-go"
)

type Client struct {
	*sendgrid.Client
}

func NewClient(cf *config.AppConfig) *Client {
	return &Client{sendgrid.NewSendClient(cf.SendGridApiKey)}
}
