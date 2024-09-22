package operation

import (
	"context"
	"database/sql"
	"fmt"
)

type UpdateDisableSubscription struct {
	pgConn *sql.DB
}

type UpdateDisableSubscriptionParams struct {
	Email              string
	NewsletterPublicID string
}

func NewUpdateDisableSubscription(pgConn *sql.DB) *UpdateDisableSubscription {
	return &UpdateDisableSubscription{
		pgConn: pgConn,
	}
}

// Execute TODO: Tests this properly
func (u *UpdateDisableSubscription) Execute(ctx context.Context, p *UpdateDisableSubscriptionParams) error {
	const query = `
		UPDATE subscriptions SET disabled_at = CURRENT_TIMESTAMP
	 	WHERE subscriber_email = $1 AND newsletter_id = (
	 	    SELECT id FROM newsletters WHERE public_id = $2
	 	);
	`

	_, err := u.pgConn.ExecContext(ctx, query, p.Email, p.NewsletterPublicID)
	if err != nil {
		return fmt.Errorf("failed to execute update unsubscribe : %w", err)
	}

	return nil
}
