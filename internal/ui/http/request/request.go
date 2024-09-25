package request

type CreateNewsletterRequest struct {
	Name        string  `json:"name" binding:"required" example:"Tiktok News 420"`
	Description *string `json:"description,omitempty" example:"Amazing news from the TikTok world. You would not believe number 4."`
}

type SubscribeToNewsletter struct {
	Email string `json:"email" binding:"required" example:"test@test.com"`
}

type UserRequest struct {
	Email    string `json:"email" binding:"required" example:"test@test.com"`
	Password string `json:"password" binding:"required" example:"Pa$$W0rD"`
}
