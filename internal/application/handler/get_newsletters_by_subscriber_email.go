package handler

import (
	"context"

	"github.com/javor454/newsletter-assignment/internal/application/dto"
	"github.com/javor454/newsletter-assignment/internal/domain"
)

type GetNewslettersBySubscriberEmail interface {
	GetBySubscriberEmail(ctx context.Context, email *domain.Email, pageSize, pageNumber int) ([]*domain.Newsletter, *dto.Pagination, error)
}

type GetNewslettersBySubscriberEmailHandler struct {
	getNewslettersByUserID GetNewslettersBySubscriberEmail
}

func NewGetNewslettersBySubscriberEmailHandler(gnbui GetNewslettersBySubscriberEmail) *GetNewslettersBySubscriberEmailHandler {
	return &GetNewslettersBySubscriberEmailHandler{getNewslettersByUserID: gnbui}
}

func (g *GetNewslettersBySubscriberEmailHandler) Handle(
	ctx context.Context,
	email string,
	pageSize, pageNumber int,
) ([]*domain.Newsletter, *dto.Pagination, error) {
	emailVo, err := domain.NewEmail(email)
	if err != nil {
		return nil, nil, err
	}

	newsletters, pagination, err := g.getNewslettersByUserID.GetBySubscriberEmail(ctx, emailVo, pageSize, pageNumber)
	if err != nil {
		return nil, nil, err
	}

	return newsletters, pagination, nil
}
