package handler

import (
	"context"
	"fmt"

	"github.com/javor454/newsletter-assignment/internal/application/dto"
	"github.com/javor454/newsletter-assignment/internal/domain"
)

type LoginUserHandler struct {
}

func NewLoginUserHandler() *LoginUserHandler {
	return &LoginUserHandler{}
}

func (r *LoginUserHandler) Handle(ctx context.Context, email string, password string) (*dto.Token, error) {
	emailVo, err := domain.NewEmail(email)
	if err != nil {
		return nil, err
	}
	pass, err := domain.NewPassword(password)
	if err != nil {
		return nil, err
	}

	// TODO: retrieve email from db and compare
	if emailVo.String() != "email from db" {
		return nil, fmt.Errorf("invalid email")
	}

	// TODO: retrieve hash from db
	if !pass.IsEqual("hash from db") {
		return nil, fmt.Errorf("invalid password")
	}

	// TODO: create token
	return &dto.Token{}, nil
}
