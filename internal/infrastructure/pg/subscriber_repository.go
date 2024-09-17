package pg

import (
	"context"
	"time"

	"github.com/javor454/newsletter-assignment/internal/domain"
	"github.com/javor454/newsletter-assignment/internal/infrastructure/pg/operation"
)

type SubscriberRepository struct {
	getNewsletterByPublicID *operation.GetNewsletterIDByPublicID
	createSubscription      *operation.CreateSubscription
}

func NewSubscriberRepository(
	gn *operation.GetNewsletterIDByPublicID,
	cs *operation.CreateSubscription,
) *SubscriberRepository {
	return &SubscriberRepository{createSubscription: cs, getNewsletterByPublicID: gn}
}

func (s *SubscriberRepository) Subscribe(ctx context.Context, subscription *domain.Subscription) error {
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	idRow, err := s.getNewsletterByPublicID.Execute(
		ctx,
		&operation.GetNewsletterIDByPublicIDParams{PublicID: subscription.NewsletterPublicID().String()},
	)
	if err != nil {
		return err
	}

	if err := s.createSubscription.Execute(ctx, &operation.CreateSubscriptionParams{
		ID:              subscription.ID().String(),
		SubscriberEmail: subscription.Email().String(),
		NewsletterID:    idRow.ID,
	}); err != nil {
		return err
	}

	return nil
}
