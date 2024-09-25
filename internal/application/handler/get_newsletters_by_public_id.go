package handler

import (
	"context"

	"github.com/javor454/newsletter-assignment/internal/domain"
)

type GetNewsletterByPublicID interface {
	GetByPublicID(ctx context.Context, publicID *domain.ID) (*domain.Newsletter, error)
}

type GetNewsletterByPublicIDHandler struct {
	getNewsletterByPublicID GetNewsletterByPublicID
}

func NewGetNewsletterByPublicIDHandler(gnbpi GetNewsletterByPublicID) *GetNewsletterByPublicIDHandler {
	return &GetNewsletterByPublicIDHandler{getNewsletterByPublicID: gnbpi}
}

func (g *GetNewsletterByPublicIDHandler) Handle(ctx context.Context, publicID string) (*domain.Newsletter, error) {
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
