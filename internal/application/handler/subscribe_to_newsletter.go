package handler

import (
	"context"

	"github.com/javor454/newsletter-assignment/internal/domain"
)

type SubscribeToNewsletter interface {
	Subscribe(ctx context.Context, subscription *domain.Subscription) error
}

type SubscribeToNewsletterHandler struct {
	subscribeToNewsletter SubscribeToNewsletter
}

func NewSubscribeToNewsletterHandler(stn SubscribeToNewsletter) *SubscribeToNewsletterHandler {
	return &SubscribeToNewsletterHandler{subscribeToNewsletter: stn}
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

	return nil
}
