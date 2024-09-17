package response

import (
	"time"

	"github.com/javor454/newsletter-assignment/internal/domain"
)

type GetNewslettersBySubscriberEmailResponse struct {
	PublicID    string  `json:"public_id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	CreatedAt   string  `json:"created_at"`
}

func CreateGetNewslettersBySubscriberEmailResponseFromEntity(n *domain.Newsletter) GetNewslettersBySubscriberEmailResponse {
	return GetNewslettersBySubscriberEmailResponse{
		PublicID:    n.PublicID().String(),
		Name:        n.Name(),
		Description: n.Description(),
		CreatedAt:   n.CreatedAt().Format(time.RFC3339),
	}
}
