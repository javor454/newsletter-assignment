package helper

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type UserRow struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
}

func GetUserByEmail(email string, pgConn *sql.DB) (*UserRow, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const query = `
		SELECT id, email, password_hash, created_at
		FROM users
		WHERE email = $1;
	`
	var res UserRow
	if err := pgConn.QueryRowContext(ctx, query, email).Scan(
		&res.ID,
		&res.Email,
		&res.PasswordHash,
		&res.CreatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &res, nil
}

func CreateUser(id, email, passwordHash string, pgConn *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const query = `
		INSERT INTO users(id, email, password_hash)
		VALUES ($1, $2, $3);
	`

	_, err := pgConn.ExecContext(ctx, query, id, email, passwordHash)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func RemoveUsersByUserID(ids []string, pgConn *sql.DB) error {
	if len(ids) == 0 {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const query = "DELETE FROM users WHERE id = ANY($1);"
	_, err := pgConn.ExecContext(ctx, query, pq.Array(ids))
	if err != nil {
		return fmt.Errorf("failed to remove users: %w", err)
	}

	return nil
}

type NewsletterRow struct {
	ID          string    `json:"id"`
	PublicID    string    `json:"public_id"`
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func GetNewslettersByUserID(userID string, pgConn *sql.DB) ([]*NewsletterRow, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const query = "SELECT id, public_id, user_id, name, description, created_at FROM newsletters WHERE user_id = $1;"

	rows, err := pgConn.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get newsletters: %w", err)
	}

	newsletters := make([]*NewsletterRow, 0, 10)
	for rows.Next() {
		var row NewsletterRow
		if err := rows.Scan(
			&row.ID,
			&row.PublicID,
			&row.UserID,
			&row.Name,
			&row.Description,
			&row.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan newsletters: %w", err)
		}

		newsletters = append(newsletters, &row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get newsletters by user id operation failed: %w", err)
	}

	return newsletters, nil
}

func CreateNewsletter(newsletterID, publicID, userID, name, description string, pgConn *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const query = `
		INSERT INTO newsletters(id, public_id, user_id, name, description)
		VALUES ($1, $2, $3, $4, $5);
	`

	_, err := pgConn.ExecContext(ctx, query, newsletterID, publicID, userID, name, description)
	if err != nil {
		return fmt.Errorf("failed to create newsletter: %w", err)
	}

	return nil
}

func RemoveNewsletterByID(ids []string, pgConn *sql.DB) error {
	if len(ids) == 0 {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const query = "DELETE FROM newsletters WHERE id = ANY($1);"
	_, err := pgConn.ExecContext(ctx, query, pq.Array(ids))
	if err != nil {
		return fmt.Errorf("failed to remove newsletters: %w", err)
	}

	return nil
}
