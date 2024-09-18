package handler

import (
	"context"

	"github.com/javor454/newsletter-assignment/internal/domain"
)

type SubscribeToNewsletterRepository interface {
	Subscribe(ctx context.Context, subscription *domain.Subscription) error
}

type SubscribeToNewsletterCache interface {
	CacheSubscription(ctx context.Context, subscription *domain.Subscription) error
}

type SubscribeToNewsletterHandler struct {
	subscribeToNewsletter SubscribeToNewsletterRepository
	subscriptionCache     SubscribeToNewsletterCache
}

func NewSubscribeToNewsletterHandler(
	stn SubscribeToNewsletterRepository,
	sc SubscribeToNewsletterCache,
) *SubscribeToNewsletterHandler {
	return &SubscribeToNewsletterHandler{subscribeToNewsletter: stn, subscriptionCache: sc}
}

func (r *SubscribeToNewsletterHandler) Handle(ctx context.Context, newsletterPublicID, email string) error {
	emailVo, err := domain.NewEmail(email)
	if err != nil {
		return err
	}
	pubID, err := domain.CreateIDFromExisting(newsletterPublicID)
	if err != nil {
		return err
	}

	subscription := domain.NewSubscription(pubID, emailVo)
	if err := r.subscribeToNewsletter.Subscribe(ctx, subscription); err != nil {
		return err
	}

	if err := r.subscriptionCache.CacheSubscription(ctx, subscription); err != nil {
		return err
	}

	return nil
}
