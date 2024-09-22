package handler

import (
	"context"

	"github.com/javor454/newsletter-assignment/internal/domain"
)

type RegisterUser interface {
	Register(ctx context.Context, user *domain.User) error
}

type RegisterUserHandler struct {
	registerUser  RegisterUser
	generateToken GenerateToken
}

func NewRegisterUserHandler(us RegisterUser, ts GenerateToken) *RegisterUserHandler {
	return &RegisterUserHandler{
		registerUser:  us,
		generateToken: ts,
	}
}

func (r *RegisterUserHandler) Handle(ctx context.Context, email string, password string) (string, error) {
	emailVo, err := domain.NewEmail(email)
	if err != nil {
		return "", err
	}
	pass, err := domain.NewPassword(password)
	if err != nil {
		return "", err
	}

	user := domain.NewUser(emailVo, pass)
	if err := r.registerUser.Register(ctx, user); err != nil {
		return "", err
	}

	token, err := r.generateToken.GenerateUserToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}
