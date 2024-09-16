package pg

import (
	"context"
	"time"

	"github.com/javor454/newsletter-assignment/internal/domain"
	"github.com/javor454/newsletter-assignment/internal/infrastructure/pg/operation"
)

type NewsletterRepository struct {
	createNewsletter       *operation.CreateNewsletter
	getNewslettersByUserID *operation.GetNewslettersByUserID
}

func NewNewsletterRepository(cn *operation.CreateNewsletter, gn *operation.GetNewslettersByUserID) *NewsletterRepository {
	return &NewsletterRepository{
		createNewsletter:       cn,
		getNewslettersByUserID: gn,
	}
}

func (u *NewsletterRepository) Create(ctx context.Context, userID *domain.ID, newsletter *domain.Newsletter) error {
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	if err := u.createNewsletter.Execute(ctx, &operation.CreateNewsletterParams{
		UserID:       userID.String(),
		NewsletterID: newsletter.Id().String(),
		Name:         newsletter.Name(),
		Description:  newsletter.Description(),
	}); err != nil {
		return err
	}

	return nil
}

// TODO: test for nonexistent user id
func (u *NewsletterRepository) GetByUserID(ctx context.Context, userID *domain.ID, pageSize, pageNumber int) ([]*domain.Newsletter, error) {
	rows, err := u.getNewslettersByUserID.Execute(ctx, &operation.GetNewslettersByUserIDParams{
		UserID:     userID.String(),
		PageSize:   pageSize,
		PageNumber: pageNumber,
	})
	if err != nil {
		return nil, err
	}

	newsletters := make([]*domain.Newsletter, 0, len(rows))
	for _, row := range rows {
		id := domain.CreateIDFromExisting(row.ID)
		newsletters = append(newsletters, domain.CreateNewsletterFromExisting(id, row.Name, row.Description))
	}

	return newsletters, nil
}
