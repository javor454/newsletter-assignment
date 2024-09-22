package handler

import (
	"context"

	"github.com/javor454/newsletter-assignment/internal/domain"
)

type GetUser interface {
	GetByEmailAndPassword(ctx context.Context, email *domain.Email, pass *domain.Password) (*domain.User, error)
}

type GenerateToken interface {
	GenerateUserToken(user *domain.User) (string, error)
}

type LoginUserHandler struct {
	getUser       GetUser
	generateToken GenerateToken
}

func NewLoginUserHandler(ur GetUser, ts GenerateToken) *LoginUserHandler {
	return &LoginUserHandler{getUser: ur, generateToken: ts}
}

func (r *LoginUserHandler) Handle(ctx context.Context, email string, password string) (string, error) {
	emailVo, err := domain.NewEmail(email)
	if err != nil {
		return "", err
	}
	pass, err := domain.NewPassword(password)
	if err != nil {
		return "", err
	}

	user, err := r.getUser.GetByEmailAndPassword(ctx, emailVo, pass)
	if err != nil {
		return "", err
	}

	token, err := r.generateToken.GenerateUserToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}
