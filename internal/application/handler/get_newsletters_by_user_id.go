package handler

import (
	"context"

	"github.com/javor454/newsletter-assignment/internal/application/dto"
	"github.com/javor454/newsletter-assignment/internal/domain"
)

type GetNewslettersByUserID interface {
	GetByUserID(ctx context.Context, userID *domain.ID, pageSize, pageNumber int) ([]*domain.Newsletter, *dto.Pagination, error)
}

type GetNewslettersByUserIDHandler struct {
	getNewslettersByUserID GetNewslettersByUserID
}

func NewGetNewslettersByUserIDHandler(gnbui GetNewslettersByUserID) *GetNewslettersByUserIDHandler {
	return &GetNewslettersByUserIDHandler{getNewslettersByUserID: gnbui}
}

func (g *GetNewslettersByUserIDHandler) Handle(ctx context.Context, userID string, pageSize, pageNumber int) ([]*domain.Newsletter, *dto.Pagination, error) {
	id := domain.CreateIDFromExisting(userID)
	newsletters, pagination, err := g.getNewslettersByUserID.GetByUserID(ctx, id, pageSize, pageNumber)
	if err != nil {
		return nil, nil, err
	}

	return newsletters, pagination, nil
}
