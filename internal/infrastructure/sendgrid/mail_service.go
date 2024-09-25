package sendgrid

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/javor454/newsletter-assignment/app/config"
	"github.com/javor454/newsletter-assignment/app/logger"
	"github.com/javor454/newsletter-assignment/app/sendgrid"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

const (
	SubscribedTemplateName = "subscribed"
)

type MailService struct {
	lg        logger.Logger
	conf      *config.AppConfig
	client    *sendgrid.Client
	templates map[string]*template.Template
}

func NewMailService(lg logger.Logger, conf *config.AppConfig, client *sendgrid.Client) *MailService {
	m := &MailService{
		lg:        lg,
		conf:      conf,
		client:    client,
		templates: make(map[string]*template.Template),
	}

	if err := m.loadTemplates(conf.SendGridTemplateDir); err != nil {
		panic("failed to load templates: " + err.Error())
	}
	lg.Info("[EMAIL] Service initialized, templates loaded")

	return m
}

func (m *MailService) SendSubscribed(recipient, newsletterPublicID, token string) error {
	tmpl, ok := m.templates[SubscribedTemplateName]
	if !ok {
		return fmt.Errorf("template \"%s\" not loaded", SubscribedTemplateName)
	}

	link := m.createUnsubscribeLink(newsletterPublicID, token)
	m.lg.Debugf("[EMAIL] Unsubscribe link: %s", link)

	var body bytes.Buffer
	if err := tmpl.Execute(&body, map[string]any{
		"Recipient": recipient,
		"Link":      link,
	}); err != nil {
		return fmt.Errorf("template \"%s\" execute error: %w", SubscribedTemplateName, err)
	}

	from := mail.NewEmail("Jiri", "javornicky.jiri@gmail.com")
	to := mail.NewEmail("Recipient", recipient)
	message := mail.NewSingleEmail(from, "Subscribed to newsletter", to, body.String(), body.String())

	response, err := m.client.Send(message)
	if err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}

	m.lg.Debugf("[EMAIL] Send success. Status Code: %d, Body: %s", response.StatusCode, response.Body)

	return nil
}

func (m *MailService) createUnsubscribeLink(newsletterPublicID string, token string) string {
	return fmt.Sprintf(
		"%s:%d/api/v1/unsubscribe?newsletter_public_id=%s&token=%s",
		m.conf.Host,
		m.conf.HttpPort,
		newsletterPublicID,
		token,
	)
}

func (m *MailService) loadTemplates(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to load templates from %s: %w", path, err)
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".html" {
			return nil
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read template file %s: %w", path, err)
		}

		name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		tmpl, err := template.New(name).Parse(string(b))
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", name, err)
		}

		m.templates[name] = tmpl

		return nil
	})
}
