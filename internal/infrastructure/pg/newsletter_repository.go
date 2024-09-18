package pg

import (
	"context"
	"fmt"
	"time"

	"github.com/javor454/newsletter-assignment/internal/application/dto"
	"github.com/javor454/newsletter-assignment/internal/domain"
	"github.com/javor454/newsletter-assignment/internal/infrastructure/pg/operation"
)

type NewsletterRepository struct {
	createNewsletter                *operation.CreateNewsletter
	getNewslettersByUserID          *operation.GetNewslettersByUserID
	getNewslettersBySubscriberEmail *operation.GetNewslettersBySubscriberEmail
}

func NewNewsletterRepository(
	cn *operation.CreateNewsletter,
	gn *operation.GetNewslettersByUserID,
	gns *operation.GetNewslettersBySubscriberEmail,
) *NewsletterRepository {
	return &NewsletterRepository{
		createNewsletter:                cn,
		getNewslettersByUserID:          gn,
		getNewslettersBySubscriberEmail: gns,
	}
}

func (u *NewsletterRepository) Create(ctx context.Context, userID *domain.ID, newsletter *domain.Newsletter) error {
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	if err := u.createNewsletter.Execute(ctx, &operation.CreateNewsletterParams{
		UserID:      userID.String(),
		ID:          newsletter.ID().String(),
		PublicID:    newsletter.PublicID().String(),
		Name:        newsletter.Name(),
		Description: newsletter.Description(),
		CreatedAt:   newsletter.CreatedAt(),
	}); err != nil {
		return err
	}

	return nil
}

func (u *NewsletterRepository) GetBySubscriberEmail(ctx context.Context, email *domain.Email, pageSize, pageNumber int) ([]*domain.Newsletter, *dto.Pagination, error) {
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	rows, pagination, err := u.getNewslettersBySubscriberEmail.Execute(ctx, &operation.GetNewslettersBySubscriberEmailParams{
		Email:      email.String(),
		PageSize:   pageSize,
		PageNumber: pageNumber,
	})
	if err != nil {
		return nil, nil, err
	}

	newsletters := make([]*domain.Newsletter, 0, len(rows))
	for _, row := range rows {
		id, err := domain.CreateIDFromExisting(row.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid uuid format in db %s", err.Error())
		}
		publicID, err := domain.CreateIDFromExisting(row.PublicID)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid uuid format in db %s", err.Error())
		}
		newsletters = append(newsletters, domain.CreateNewsletterFromExisting(id, publicID, row.Name, row.Description, row.CreatedAt))
	}

	return newsletters, pagination, nil
}

// TODO: test for nonexistent user id

func (u *NewsletterRepository) GetByUserID(ctx context.Context, userID *domain.ID, pageSize, pageNumber int) ([]*domain.Newsletter, *dto.Pagination, error) {
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond) // TODO: scale with pageSize?
	defer cancel()

	rows, pagination, err := u.getNewslettersByUserID.Execute(ctx, &operation.GetNewslettersByUserIDParams{
		UserID:     userID.String(),
		PageSize:   pageSize,
		PageNumber: pageNumber,
	})
	if err != nil {
		return nil, nil, err
	}

	newsletters := make([]*domain.Newsletter, 0, len(rows))
	for _, row := range rows {
		id, err := domain.CreateIDFromExisting(row.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid uuid format in db %s", err.Error())
		}
		publicID, err := domain.CreateIDFromExisting(row.PublicID)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid uuid format in db %s", err.Error())
		}
		newsletters = append(newsletters, domain.CreateNewsletterFromExisting(id, publicID, row.Name, row.Description, row.CreatedAt))
	}

	return newsletters, pagination, nil
}
