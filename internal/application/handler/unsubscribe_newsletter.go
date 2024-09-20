package handler

import (
	"context"

	"github.com/javor454/newsletter-assignment/internal/domain"
)

type UnsubscribeNewsletterRepository interface {
	Unsubscribe(ctx context.Context, email *domain.Email, newsletterPublicID *domain.ID) error
}

type UnsubscribeNewsletterHandler struct {
	unsubscribeNewsletter UnsubscribeNewsletterRepository
}

func NewUnsubscribeNewsletterHandler(
	unr UnsubscribeNewsletterRepository,
) *UnsubscribeNewsletterHandler {
	return &UnsubscribeNewsletterHandler{unsubscribeNewsletter: unr}
}

func (r *UnsubscribeNewsletterHandler) Handle(ctx context.Context, newsletterPublicID, email string) error {
	emailVo, err := domain.NewEmail(email)
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

	return nil
}
