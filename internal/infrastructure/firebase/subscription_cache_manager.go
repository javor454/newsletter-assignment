package firebase

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/javor454/newsletter-assignment/app/firebase"
)

type SubscriptionCacheManager struct {
	client *firebase.Client
}

func NewSubscriptionCacheManager(client *firebase.Client) *SubscriptionCacheManager {
	return &SubscriptionCacheManager{client: client}
}

func (s *SubscriptionCacheManager) CacheSubscription(ctx context.Context, email, newsletterID string) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	encodedEmail := base64.URLEncoding.EncodeToString([]byte(email))

	var records map[string]bool
	if err := s.client.NewRef("subscriptions").Child(encodedEmail).Get(ctx, &records); err != nil {
		return fmt.Errorf("could not get subscription records: %w", err)
	}

	if records == nil {
		records = make(map[string]bool)
	}

	records[newsletterID] = true
	if err := s.client.NewRef("subscriptions").Child(encodedEmail).Set(ctx, &records); err != nil {
		return fmt.Errorf("could not set subscription records: %w", err)
	}

	return nil
}
