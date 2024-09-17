package handler

import (
	"context"

	"github.com/javor454/newsletter-assignment/internal/domain"
)

type CreateNewsletter interface {
	Create(ctx context.Context, userID *domain.ID, newsletter *domain.Newsletter) error
}

type CreateNewsletterHandler struct {
	createNewsletter CreateNewsletter
}

func NewCreateNewsletterHandler(cn CreateNewsletter) *CreateNewsletterHandler {
	return &CreateNewsletterHandler{createNewsletter: cn}
}

func (r *CreateNewsletterHandler) Handle(ctx context.Context, userID, name string, description *string) error {
	id := domain.CreateIDFromExisting(userID)
	newsletter := domain.NewNewsletter(name, description)

	if err := r.createNewsletter.Create(ctx, id, newsletter); err != nil {
		return err
	}

	return nil
}
