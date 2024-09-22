package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/javor454/newsletter-assignment/internal/domain"
)

type TokenManager struct {
	secret string
	host   string
}

func NewTokenManager(secret string, host string) *TokenManager {
	return &TokenManager{
		secret: secret,
		host:   host,
	}
}

func (t *TokenManager) GenerateUserToken(user *domain.User) (string, error) {
	return t.generateToken(user.ID().String(), 1*time.Hour)
}

func (t *TokenManager) GenerateSubscriptionToken(email *domain.Email) (string, error) {
	return t.generateToken(email.String(), 0)
}

func (t *TokenManager) generateToken(subject string, expiration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = subject
	claims["iss"] = t.host
	claims["iat"] = time.Now().Unix()
	if expiration > 0 {
		claims["exp"] = time.Now().Add(expiration).Unix()
	}

	tokenStr, err := token.SignedString([]byte(t.secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenStr, nil
}

func (t *TokenManager) ParseToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(t.secret), nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return "", fmt.Errorf("token expired")
		}
	}

	subject, exists := claims["sub"].(string)
	if !exists || subject == "" {
		return "", fmt.Errorf("subject missing")
	}

	return subject, nil
}
