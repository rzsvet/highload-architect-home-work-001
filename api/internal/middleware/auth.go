package middleware

import (
	"api/internal/service"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		claims, err := userService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Безопасное извлечение user_id
		userID, err := extractUserID(claims)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("email", claims["email"])
		c.Next()
	}
}

// extractUserID безопасно извлекает user_id из claims
func extractUserID(claims map[string]interface{}) (int, error) {
	userIDValue, exists := claims["user_id"]
	if !exists {
		return 0, errors.New("user_id not found in token")
	}

	switch v := userIDValue.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil // Конвертируем float64 в int
	case int64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, errors.New("invalid user_id type in token")
	}
}
