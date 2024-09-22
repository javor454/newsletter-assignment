package handler

import (
	"context"

	"github.com/javor454/newsletter-assignment/internal/domain"
)

type TokenGenerator interface {
	GenerateSubscriptionToken(email *domain.Email) (string, error)
}

type SubscribeToNewsletterRepository interface {
	Subscribe(ctx context.Context, subscription *domain.Subscription) error
}

type SubscribeToNewsletterHandler struct {
	tokenGenerator        TokenGenerator
	subscribeToNewsletter SubscribeToNewsletterRepository
}

func NewSubscribeToNewsletterHandler(
	tg TokenGenerator,
	stn SubscribeToNewsletterRepository,
) *SubscribeToNewsletterHandler {
	return &SubscribeToNewsletterHandler{tokenGenerator: tg, subscribeToNewsletter: stn}
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

	token, err := r.tokenGenerator.GenerateSubscriptionToken(emailVo)
	if err != nil {
		return err
	}

	subscription := domain.NewSubscription(pubID, emailVo, token)

	if err := r.subscribeToNewsletter.Subscribe(ctx, subscription); err != nil {
		return err
	}

	return nil
}
