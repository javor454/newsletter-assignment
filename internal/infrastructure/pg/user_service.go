package pg

import (
	"context"

	"github.com/javor454/newsletter-assignment/internal/domain"
	"github.com/javor454/newsletter-assignment/internal/infrastructure/pg/operation"
)

type UserService struct {
	createUser *operation.CreateUser
}

func NewUserService(createUser *operation.CreateUser) *UserService {
	return &UserService{
		createUser: createUser,
	}
}

func (u *UserService) RegisterUser(ctx context.Context, user *domain.User) error {
	p := &operation.CreateUserParams{
		Email:        user.Email().String(),
		PasswordHash: user.Password().String(),
	}
	return nil
}
