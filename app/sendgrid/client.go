package sendgrid

import (
	"github.com/javor454/newsletter-assignment/app/config"
	"github.com/javor454/newsletter-assignment/app/logger"
	"github.com/sendgrid/sendgrid-go"
)

type Client struct {
	*sendgrid.Client
}

func NewClient(lg logger.Logger, cf *config.AppConfig) *Client {
	lg.Debug("[SENDGRID] Creating client...")
	c := &Client{sendgrid.NewSendClient(cf.SendGridApiKey)}
	lg.Info("[SENDGRID] Created")

	return c
}
