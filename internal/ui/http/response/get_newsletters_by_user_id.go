package response

import (
	"time"

	"github.com/javor454/newsletter-assignment/internal/domain"
)

type GetNewslettersByUserIDResponse struct {
	ID          string  `json:"id"`
	PublicID    string  `json:"public_id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	CreatedAt   string  `json:"created_at"`
}

func CreateGetNewslettersByUserIDResponseFromEntity(n *domain.Newsletter) GetNewslettersByUserIDResponse {
	return GetNewslettersByUserIDResponse{
		ID:          n.ID().String(),
		PublicID:    n.PublicID().String(),
		Name:        n.Name(),
		Description: n.Description(),
		CreatedAt:   n.CreatedAt().Format(time.RFC3339),
	}
}
