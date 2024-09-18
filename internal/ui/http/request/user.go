package request

type RegisterUserRequest struct {
	Email    string `json:"email" binding:"required" example:"test@test.com"`
	Password string `json:"password" binding:"required" example:"Pa$$W0rD"`
}
