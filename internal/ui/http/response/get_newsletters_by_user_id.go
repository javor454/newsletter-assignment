package response

import "github.com/javor454/newsletter-assignment/internal/domain"

type GetNewslettersByUserIDResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

func CreateGetNewslettersByUserIDResponseFromEntity(n *domain.Newsletter) GetNewslettersByUserIDResponse {
	return GetNewslettersByUserIDResponse{
		ID:          n.Id().String(),
		Name:        n.Name(),
		Description: n.Description(),
	}
}
