package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"service-user/internal/app/errs"
	"service-user/internal/app/utils"
)

// AuthMiddleware проверяет JWT токен
func AuthMiddleware(jwtManager *utils.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Извлекаем токен из cookies
		tokenString, err := extractTokenFromCookie(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Парсим и валидируем токен
		claims, err := jwtManager.DecodeJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Извлекаем user_id
		userID, err := extractUserID(claims)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token payload"})
			c.Abort()
			return
		}

		// Передаем user_id в контекст запроса
		c.Set("user_id", userID)
		c.Next()
	}
}

// extractTokenFromCookie извлекает access_token из cookie
func extractTokenFromCookie(c *gin.Context) (string, error) {
	cookie, err := c.Cookie("access_token")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(cookie), nil
}

// extractUserID извлекает user_id из claims
func extractUserID(claims jwt.Claims) (uuid.UUID, error) {
	mapClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, errs.ErrTokenInvalid
	}

	sub, ok := mapClaims["sub"].(string)

	if !ok {
		return uuid.Nil, errs.ErrTokenInvalid
	}

	return uuid.Parse(sub)
}
