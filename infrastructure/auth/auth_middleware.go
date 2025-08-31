package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	svc "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/services"
)

type AuthMiddleware struct {
	JWTService svc.IJWTService
}

func NewAuthMiddleware(jwtService svc.IJWTService) *AuthMiddleware {
	return &AuthMiddleware{JWTService: jwtService}
}

func (a *AuthMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authorization header missing",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid authorization header format",
			})
			return
		}

		token := parts[1]
		userID, lang, err := a.JWTService.ValidateAccessToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid or expired token",
			})
			return
		}

		c.Set("userID", userID)
		c.Set("preferredLanguage", lang)

		c.Next()
	}
}
