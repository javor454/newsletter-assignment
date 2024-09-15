package handler

import (
	"context"

	"github.com/javor454/newsletter-assignment/internal/application/dto"
	"github.com/javor454/newsletter-assignment/internal/domain"
)

type UserService interface {
	RegisterUser(ctx context.Context, user *domain.User) error
}

type RegisterUserHandler struct {
	us UserService
}

func NewRegisterUserHandler(us UserService) *RegisterUserHandler {
	return &RegisterUserHandler{
		us: us,
	}
}

func (r *RegisterUserHandler) Handle(ctx context.Context, email string, password string) (*dto.Token, error) {
	emailVo, err := domain.NewEmail(email)
	if err != nil {
		return nil, err
	}
	pass, err := domain.NewPassword(password)
	if err != nil {
		return nil, err
	}

	user := domain.NewUser(emailVo, pass)
	if err := r.us.RegisterUser(ctx, user); err != nil {
		return nil, err
	}

	// TODO: create token

	return &dto.Token{}, nil
}
