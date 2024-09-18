package firebase

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"firebase.google.com/go/v4/db"
	"github.com/javor454/newsletter-assignment/internal/domain"
)

type SubscriptionRepository struct {
	client *db.Client
}

func NewSubscriptionRepository(client *db.Client) *SubscriptionRepository {
	return &SubscriptionRepository{client: client}
}

func (s *SubscriptionRepository) CacheSubscription(ctx context.Context, subscription *domain.Subscription) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	encodedEmail := base64.URLEncoding.EncodeToString([]byte(subscription.Email().String()))

	ref := s.client.NewRef("subscriptions").Child(encodedEmail)

	// upserts newsletterID
	err := ref.Update(ctx, map[string]interface{}{
		subscription.NewsletterPublicID().String(): true,
	})
	if err != nil {
		return fmt.Errorf("failed to upsert subscription: %s", err.Error())
	}

	return nil
}
