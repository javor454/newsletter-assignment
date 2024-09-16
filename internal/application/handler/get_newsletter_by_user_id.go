package handler

import (
	"context"

	"github.com/javor454/newsletter-assignment/internal/domain"
)

type GetNewslettersByUserID interface {
	GetByUserID(ctx context.Context, userID *domain.ID, pageSize, pageNumber int) ([]*domain.Newsletter, error)
}

type GetNewslettersByUserIDHandler struct {
	getNewslettersByUserID GetNewslettersByUserID
}

func NewGetNewslettersByUserIDHandler(gnbui GetNewslettersByUserID) *GetNewslettersByUserIDHandler {
	return &GetNewslettersByUserIDHandler{getNewslettersByUserID: gnbui}
}

func (g *GetNewslettersByUserIDHandler) Handle(ctx context.Context, userID string, pageSize, pageNumber int) ([]*domain.Newsletter, error) {
	id := domain.CreateIDFromExisting(userID)
	newsletters, err := g.getNewslettersByUserID.GetByUserID(ctx, id, pageSize, pageNumber)
	if err != nil {
		return nil, err
	}

	return newsletters, nil
}
