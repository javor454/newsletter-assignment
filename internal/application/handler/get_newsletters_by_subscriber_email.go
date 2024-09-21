package handler

import (
	"context"

	"github.com/javor454/newsletter-assignment/internal/application/dto"
	"github.com/javor454/newsletter-assignment/internal/domain"
)

type GetNewslettersBySubscriptionEmail interface {
	GetBySubscriptionEmail(ctx context.Context, email *domain.Email, pageSize, pageNumber int) ([]*domain.Newsletter, *dto.Pagination, error)
}

type GetNewslettersBySubscriptionEmailHandler struct {
	getNewslettersBySubscriptionEmail GetNewslettersBySubscriptionEmail
}

func NewGetNewslettersBySubscriptionEmailHandler(gnbui GetNewslettersBySubscriptionEmail) *GetNewslettersBySubscriptionEmailHandler {
	return &GetNewslettersBySubscriptionEmailHandler{getNewslettersBySubscriptionEmail: gnbui}
}

func (g *GetNewslettersBySubscriptionEmailHandler) Handle(
	ctx context.Context,
	email string,
	pageSize, pageNumber int,
) ([]*domain.Newsletter, *dto.Pagination, error) {
	emailVo, err := domain.NewEmail(email)
	if err != nil {
		return nil, nil, err
	}

	newsletters, pagination, err := g.getNewslettersBySubscriptionEmail.GetBySubscriptionEmail(ctx, emailVo, pageSize, pageNumber)
	if err != nil {
		return nil, nil, err
	}

	return newsletters, pagination, nil
}
