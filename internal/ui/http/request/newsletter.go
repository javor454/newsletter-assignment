package request

type CreateNewsletterRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description,omitempty"`
}

type SubscribeToNewsletter struct {
	Email string `json:"email" binding:"required"`
}
