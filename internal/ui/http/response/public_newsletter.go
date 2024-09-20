package response

import (
	"time"

	"github.com/javor454/newsletter-assignment/internal/domain"
)

type PublicNewsletter struct {
	PublicID    string  `json:"public_id" example:"90c0a606-4429-44cc-9531-6f9cd038620a"`
	Name        string  `json:"name" example:"Newsletter name"`
	Description *string `json:"description,omitempty" example:"Some descriptive description"`
	CreatedAt   string  `json:"created_at" example:"2024-09-20T23:16:32Z"`
}

func CreatePublicNewsletterResponseFromEntity(n *domain.Newsletter) *PublicNewsletter {
	return &PublicNewsletter{
		PublicID:    n.PublicID().String(),
		Name:        n.Name(),
		Description: n.Description(),
		CreatedAt:   n.CreatedAt().Format(time.RFC3339),
	}
}
