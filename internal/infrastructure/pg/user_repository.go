package pg

import (
	"context"
	"fmt"
	"time"

	"github.com/javor454/newsletter-assignment/internal/domain"
	"github.com/javor454/newsletter-assignment/internal/infrastructure/crypt"
	"github.com/javor454/newsletter-assignment/internal/infrastructure/pg/operation"
)

type UserRepository struct {
	createUser     *operation.CreateUser
	getUserByEmail *operation.GetUserByEmail
}

func NewUserRepository(createUser *operation.CreateUser, getUserByEmail *operation.GetUserByEmail) *UserRepository {
	return &UserRepository{
		createUser:     createUser,
		getUserByEmail: getUserByEmail,
	}
}

func (u *UserRepository) Register(ctx context.Context, user *domain.User) error {
	bcryptHash, err := crypt.NewBcryptHashFromPassword(user.Password())
	if err != nil {
		return err
	}

	p := &operation.CreateUserParams{
		ID:           user.Id().String(),
		Email:        user.Email().String(),
		PasswordHash: bcryptHash.String(),
	}

	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	if err := u.createUser.Execute(ctx, p); err != nil {
		return err
	}

	return nil
}

func (u *UserRepository) GetByEmailAndPassword(ctx context.Context, email *domain.Email, pass *domain.Password) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	res, err := u.getUserByEmail.Execute(ctx, &operation.GetUserByEmailParams{Email: email.String()})
	if err != nil {
		return nil, err
	}

	id := domain.CreateIDFromExisting(res.ID)
	bcryptHash := crypt.CreateHashFromExisting(res.PasswordHash)
	if !bcryptHash.IsEqual(pass) {
		// TODO: custom error
		return nil, fmt.Errorf("invalid password")
	}

	return domain.CreateUserFromExisting(id, email, pass), nil
}
