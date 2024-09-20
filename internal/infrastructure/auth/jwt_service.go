package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/javor454/newsletter-assignment/internal/application/dto"
	"github.com/javor454/newsletter-assignment/internal/domain"
)

type JwtService struct {
	secret string
}

func NewJwtService(secret string) *JwtService {
	return &JwtService{
		secret: secret,
	}
}

func (j *JwtService) Generate(user *domain.User) (*dto.Token, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = user.ID().String()
	claims["iss"] = "newsletter-assignment" // TODO (nice2have): url better? from env probably
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(1 * time.Hour).Unix()

	tokenStr, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return dto.NewToken(tokenStr), nil
}

func (j *JwtService) ParseToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(j.secret), nil
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

	userID, exists := claims["sub"].(string)
	if !exists || userID == "" {
		return "", fmt.Errorf("user id missing")
	}

	return userID, nil
}
