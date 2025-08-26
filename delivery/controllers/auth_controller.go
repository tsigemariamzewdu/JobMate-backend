package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tsigemariamzewdu/JobMate-backend/domain"
)

type AuthController struct {
	AuthUsecase domain.IAuthUsecase
}

// RefreshToken is an HTTP handler that handles the token refreshing endpoint.
func (au *AuthController) RefreshToken(c *gin.Context) {
	// Retrieve the refresh token from the HTTP-only cookie.
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not found in cookie"})
		return
	}

	// Call the use case with the refresh token.
	newAccessToken, newRefreshToken, expiresIn, err := au.AuthUsecase.RefreshToken(c, refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token or session expired"})
		return
	}

	// Set the new access token in an HTTP-only cookie.
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    *newAccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // Should be true in production
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(expiresIn.Seconds()),
	})

	// Set the new refresh token in a separate HTTP-only cookie.
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    *newRefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // Should be true in production
		SameSite: http.SameSiteLaxMode,
		MaxAge: int((60 * 24 * time.Hour).Seconds()),
	})

	// Return a success response to the client.
	c.JSON(http.StatusOK, gin.H{
		"message": "Token refreshed successfully",
		"expires_in": int(expiresIn.Seconds()),
	})
}