package operation

import (
	"context"
	"database/sql"
	"fmt"
)

type CreateSubscriptionParams struct {
	ID                string
	SubscriberEmail   string
	NewsletterID      string
	SubscriptionToken string
}

// TODO: get newsletter ID can be probably merged with this
func CreateOrUpdateSubscriptionTx(ctx context.Context, tx *sql.Tx, p *CreateSubscriptionParams) error {
	const query = `
			INSERT INTO subscriptions (id, subscriber_email, newsletter_id, token)
        	VALUES ($1, $2, $3, $4)
        	ON CONFLICT (subscriber_email, newsletter_id)
        	DO UPDATE SET disabled_at = NULL;
		`

	_, err := tx.ExecContext(ctx, query, p.ID, p.SubscriberEmail, p.NewsletterID, p.SubscriptionToken)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	return nil
}
