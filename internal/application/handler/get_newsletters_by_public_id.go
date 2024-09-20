package handler

import (
	"context"

	"github.com/javor454/newsletter-assignment/internal/domain"
)

type GetNewsletterByPublicID interface {
	GetByPublicID(ctx context.Context, publicID *domain.ID) (*domain.Newsletter, error)
}

type GetNewslettersByPublicIDHandler struct {
	getNewsletterByPublicID GetNewsletterByPublicID
}

func NewGetNewslettersByPublicIDHandler(gnbpi GetNewsletterByPublicID) *GetNewslettersByPublicIDHandler {
	return &GetNewslettersByPublicIDHandler{getNewsletterByPublicID: gnbpi}
}

func (g *GetNewslettersByPublicIDHandler) Handle(ctx context.Context, publicID string) (*domain.Newsletter, error) {
	id, err := domain.CreateIDFromExisting(publicID)
	if err != nil {
		return nil, err
	}

	newsletters, err := g.getNewsletterByPublicID.GetByPublicID(ctx, id)
	if err != nil {
		return nil, err
	}

	return newsletters, nil
}
