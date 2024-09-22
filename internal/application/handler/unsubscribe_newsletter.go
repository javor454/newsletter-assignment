package handler

import (
	"context"

	"github.com/javor454/newsletter-assignment/internal/application"
	"github.com/javor454/newsletter-assignment/internal/domain"
)

type UnsubscribeNewsletterRepository interface {
	Unsubscribe(ctx context.Context, email *domain.Email, newsletterPublicID *domain.ID) error
}

type TokenParser interface {
	ParseToken(tokenStr string) (string, error)
}

type SubscriptionCache interface {
	RemoveSubscribedNewsletter(ctx context.Context, email *domain.Email, newsletterPublicID *domain.ID) error
}

type UnsubscribeNewsletterHandler struct {
	unsubscribeNewsletter UnsubscribeNewsletterRepository
	tokenParser           TokenParser
	subscriptionCache     SubscriptionCache
}

func NewUnsubscribeNewsletterHandler(
	unr UnsubscribeNewsletterRepository,
	tp TokenParser,
	sc SubscriptionCache,
) *UnsubscribeNewsletterHandler {
	return &UnsubscribeNewsletterHandler{unsubscribeNewsletter: unr, tokenParser: tp, subscriptionCache: sc}
}

func (r *UnsubscribeNewsletterHandler) Handle(ctx context.Context, newsletterPublicID, token string) error {
	parsed, err := r.tokenParser.ParseToken(token)
	if err != nil {
		return application.InvalidTokenError
	}
	emailVo, err := domain.NewEmail(parsed)
	if err != nil {
		return err
	}
	pubID, err := domain.CreateIDFromExisting(newsletterPublicID)
	if err != nil {
		return err
	}

	if err := r.unsubscribeNewsletter.Unsubscribe(ctx, emailVo, pubID); err != nil {
		return err
	}

	if err := r.subscriptionCache.RemoveSubscribedNewsletter(ctx, emailVo, pubID); err != nil {
		return err
	}

	return nil
}
