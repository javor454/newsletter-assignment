package operation

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/javor454/newsletter-assignment/internal/application"
)

type CreateSubscriptionParams struct {
	ID              string
	SubscriberEmail string
	NewsletterID    string
}

// TODO: get newsletter ID can be probably merged with this
func CreateSubscriptionTx(ctx context.Context, tx *sql.Tx, p *CreateSubscriptionParams) error {
	const (
		duplicitSubscriptionConstraint = "subscriptions_subscriber_email_newsletter_id_key"
		query                          = `
			INSERT INTO subscriptions (id, subscriber_email, newsletter_id)
			VALUES ($1, $2, $3);
		`
	)
	_, err := tx.ExecContext(ctx, query, p.ID, p.SubscriberEmail, p.NewsletterID)
	if err != nil {
		// TODO: vyzkouset duplicitni subscribe
		if strings.Contains(err.Error(), duplicitSubscriptionConstraint) {
			return application.AlreadySubscibedToNewsletterError
		}

		return fmt.Errorf("failed to create subscription: %w", err)
	}

	return nil
}
