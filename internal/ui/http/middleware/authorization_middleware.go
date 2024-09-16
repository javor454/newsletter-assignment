package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/javor454/newsletter-assignment/app/logger"
)

const UserIDKey = "user_id"

type DecodeToken interface {
	Handle(string) (string, error)
}

type AuthMiddleware struct {
	decodeToken DecodeToken
	lg          logger.Logger
}

func NewAuthMiddleware(dt DecodeToken, lg logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{decodeToken: dt, lg: lg}
}

func (a *AuthMiddleware) Handle(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})

		return
	}

	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})

		return
	}

	userID, err := a.decodeToken.Handle(bearerToken[1])
	if err != nil {
		a.lg.WithError(err).Error("Error decoding token")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})

		return
	}

	c.Set(UserIDKey, userID)
	c.Next()
}
