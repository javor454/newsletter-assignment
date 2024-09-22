package firebase

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/javor454/newsletter-assignment/app/firebase"
	"github.com/javor454/newsletter-assignment/internal/domain"
)

type SubscriptionCacheManager struct {
	client *firebase.Client
}

func NewSubscriptionCacheManager(client *firebase.Client) *SubscriptionCacheManager {
	return &SubscriptionCacheManager{client: client}
}

func (s *SubscriptionCacheManager) AddSubscribedNewsletter(ctx context.Context, email, newsletterPublicID string) error {
	encodedEmail := base64.URLEncoding.EncodeToString([]byte(email))

	getSubsCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond) // TODO: add exponential backoff or smth
	defer cancel()

	var records map[string]bool
	if err := s.client.NewRef("subscriptions").Child(encodedEmail).Get(getSubsCtx, &records); err != nil {
		return fmt.Errorf("could not get subscription records: %w", err)
	}

	if records == nil {
		records = make(map[string]bool)
	}

	setSubsCtx, cancel := context.WithTimeout(ctx, 500*time.Second) // TODO: add exponential backoff or smth
	defer cancel()

	records[newsletterPublicID] = true
	if err := s.client.NewRef("subscriptions").Child(encodedEmail).Set(setSubsCtx, &records); err != nil {
		return fmt.Errorf("could not set subscription records: %w", err)
	}

	return nil
}

func (s *SubscriptionCacheManager) RemoveSubscribedNewsletter(ctx context.Context, email *domain.Email, newsletterPublicID *domain.ID) error {
	encodedEmail := base64.URLEncoding.EncodeToString([]byte(email.String()))

	ctx, cancel := context.WithTimeout(ctx, 500*time.Second) // TODO: add exponential backoff or smth
	defer cancel()
	var records map[string]bool
	if err := s.client.NewRef("subscriptions").Child(encodedEmail).Get(ctx, &records); err != nil {
		return fmt.Errorf("could not get subscription records: %w", err)
	}

	delete(records, newsletterPublicID.String())

	ctx, cancel = context.WithTimeout(ctx, 500*time.Second) // TODO: add exponential backoff or smth
	defer cancel()

	if err := s.client.NewRef("subscriptions").Child(encodedEmail).Set(ctx, &records); err != nil {
		return fmt.Errorf("could not set subscription records: %w", err)
	}

	return nil
}
