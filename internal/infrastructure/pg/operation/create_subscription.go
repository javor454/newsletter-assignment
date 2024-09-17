package operation

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/javor454/newsletter-assignment/internal/application"
)

type CreateSubscription struct {
	pgConn *sql.DB
}

type CreateSubscriptionParams struct {
	ID              string
	SubscriberEmail string
	NewsletterID    string
}

func NewCreateSubscription(pgConn *sql.DB) *CreateSubscription {
	return &CreateSubscription{pgConn: pgConn}
}

func (o *CreateSubscription) Execute(ctx context.Context, p *CreateSubscriptionParams) error {
	const (
		duplicitSubscriptionConstraint = "subscriptions_subscriber_email_newsletter_id_key"
		query                          = `
			INSERT INTO subscriptions (id, subscriber_email, newsletter_id)
			VALUES ($1, $2, $3);
		`
	)
	// TODO: dont allow duplicit
	_, err := o.pgConn.ExecContext(ctx, query, p.ID, p.SubscriberEmail, p.NewsletterID)
	if err != nil {
		if strings.Contains(err.Error(), duplicitSubscriptionConstraint) {
			return application.AlreadySubscibedToNewsletterError
		}

		return fmt.Errorf("failed to create subscription: %s", err.Error())
	}

	return nil
}
