package controllers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tsigemariamzewdu/JobMate-backend/domain"
)

// AuthController handles HTTP requests related to authentication.
type AuthController struct {
	AuthUsecase domain.IAuthUsecase
}

// NewAuthController creates a new instance of AuthController with its dependencies.
func NewAuthController(authUsecase domain.IAuthUsecase) *AuthController {
	return &AuthController{
		AuthUsecase: authUsecase,
	}
}

// user register controller
func (ac *AuthController) Register(c *gin.Context) {

	ctx := c.Request.Context()

	var input domain.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	// context timeout handling
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// create user via usecase

	user, err := ac.AuthUsecase.Register(ctx, &input, nil)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidEmailFormat):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		case errors.Is(err, domain.ErrEmailAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		case errors.Is(err, domain.ErrPasswordHashingFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		case errors.Is(err, domain.ErrTokenGenerationFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate verification token"})
		case errors.Is(err, domain.ErrUserCreationFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User could not be created"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong", "details": err.Error()})
		}
		return
	}

	// prepare safe response - omit sensitive info
	response := gin.H{
		"message":  "User registered successfully",
		"user_id":  user.UserID,
		"email":    user.Email,
		"provider": user.Provider,
	}

	c.JSON(http.StatusCreated, response)
}

// user login controller
func (ac *AuthController) Login(c *gin.Context) {

	ctx := c.Request.Context()

	var loginUser domain.User

	// bind and validate input
	if err := c.ShouldBindJSON(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// perform login
	result, err := ac.AuthUsecase.Login(ctx, &loginUser)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidEmailFormat):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		case errors.Is(err, domain.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		case errors.Is(err, domain.ErrOAuthUserCannotLoginWithPassword):
			c.JSON(http.StatusBadRequest, gin.H{"error": "This account uses OAuth login only"})
		case errors.Is(err, domain.ErrEmailNotVerified):
			c.JSON(http.StatusForbidden, gin.H{"error": "Please verify your email address first"})
		case errors.Is(err, domain.ErrTokenGenerationFailed):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Login failed",
				"details": err.Error(),
			})
		}
		return
	}

	// prepare sanitized user response
	safeUser := gin.H{
		"user_id":   result.User.UserID,
		"email":     result.User.Email,
		"firstName": result.User.FirstName,
		"lastName":  result.User.LastName,
		"provider":  result.User.Provider,
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "auth_token",
		Value:    result.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(result.ExpiresIn.Seconds()),
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    safeUser,
	})
}

// Logout is an HTTP handler that handles user logout and clears session-related cookies.
func (au *AuthController) Logout(c *gin.Context) {
	// get the user id
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to fetch userID"})
		return
	}

	err := au.AuthUsecase.Logout(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to log user out", "details": err.Error()})
		return
	}

	// delete the authentication cookie after logout
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
	c.Status(http.StatusOK)
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
		Secure:   true, 
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int((60 * 24 * time.Hour).Seconds()),
	})

	// Return a success response to the client.
	c.JSON(http.StatusOK, gin.H{
		"message":    "Token refreshed successfully",
		"expires_in": int(expiresIn.Seconds()),
	})
}