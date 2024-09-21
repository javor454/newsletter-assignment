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
	createNewsletter                  *operation.CreateNewsletter
	getNewslettersByUserID            *operation.GetNewslettersByUserID
	getNewslettersBySubscriptionEmail *operation.GetNewslettersBySubscriptionEmail
	getNewsletterByPublicID           *operation.GetNewslettersByPublicID
}

func NewNewsletterRepository(
	cn *operation.CreateNewsletter,
	gn *operation.GetNewslettersByUserID,
	gns *operation.GetNewslettersBySubscriptionEmail,
	gnbpi *operation.GetNewslettersByPublicID,
) *NewsletterRepository {
	return &NewsletterRepository{
		createNewsletter:                  cn,
		getNewslettersByUserID:            gn,
		getNewslettersBySubscriptionEmail: gns,
		getNewsletterByPublicID:           gnbpi,
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

func (u *NewsletterRepository) GetBySubscriptionEmail(ctx context.Context, email *domain.Email, pageSize, pageNumber int) ([]*domain.Newsletter, *dto.Pagination, error) {
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	rows, pagination, err := u.getNewslettersBySubscriptionEmail.Execute(ctx, &operation.GetNewslettersBySubscriptionEmailParams{
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
			return nil, nil, fmt.Errorf("invalid uuid format in db %w", err)
		}
		publicID, err := domain.CreateIDFromExisting(row.PublicID)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid uuid format in db %w", err)
		}
		newsletters = append(newsletters, domain.CreateNewsletterFromExisting(id, publicID, row.Name, row.Description, row.CreatedAt))
	}

	return newsletters, pagination, nil
}

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
			return nil, nil, fmt.Errorf("invalid uuid format in db %w", err)
		}
		publicID, err := domain.CreateIDFromExisting(row.PublicID)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid uuid format in db %w", err)
		}
		newsletters = append(newsletters, domain.CreateNewsletterFromExisting(id, publicID, row.Name, row.Description, row.CreatedAt))
	}

	return newsletters, pagination, nil
}

func (u *NewsletterRepository) GetByPublicID(ctx context.Context, publicID *domain.ID) (*domain.Newsletter, error) {
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	row, err := u.getNewsletterByPublicID.Execute(ctx, &operation.GetNewslettersByPublicIDParams{
		PublicID: publicID.String(),
	})
	if err != nil {
		return nil, err
	}

	id, err := domain.CreateIDFromExisting(row.ID)
	if err != nil {
		return nil, err
	}
	publicID, err = domain.CreateIDFromExisting(row.PublicID)
	if err != nil {
		return nil, err
	}

	return domain.CreateNewsletterFromExisting(
		id,
		publicID,
		row.Name,
		row.Description,
		row.CreatedAt,
	), nil

}
