package sendgrid

import (
	"fmt"
	"log"

	"github.com/javor454/newsletter-assignment/app/sendgrid"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type MailService struct {
	client               *sendgrid.Client
	subscribedTemplateID string
}

func NewMailService(client *sendgrid.Client, subscribedTemplateID string) *MailService {
	return &MailService{
		client:               client,
		subscribedTemplateID: subscribedTemplateID,
	}
}

func (m *MailService) SendSubscribed(recipient string) error {
	from := mail.NewEmail("Jiri", "javornicky.jiri@gmail.com")
	email := mail.NewEmail("Recipient", recipient)
	message := mail.NewV3MailInit(from, "Subscribed to newsletter", email)

	message.SetTemplateID(m.subscribedTemplateID)

	response, err := m.client.Send(message)
	if err != nil {
		return fmt.Errorf("error sending template email: %v", err)
	}

	log.Printf("Template email sent successfully. Status Code: %d, Body: %s", response.StatusCode, response.Body)
	return nil
}
