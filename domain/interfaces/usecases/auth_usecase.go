package interfaces

import (
	"context"
	"time"

	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"
)

type IAuthUsecase interface {
	// Register registers a new user using provided user data.
	// It should validate the input, hash the password, store the user in the database,
	// and send a verification email.
	Register(ctx context.Context, user, oauthUser *models.User) (*models.User, error)
	
	// Login authenticates a user using an identifier (username or email) and password.
	// It should verify credentials, check if the email is verified, and return accesstoken, refreshtoken, user data.
	Login(ctx context.Context, input *models.User) (*models.LoginResult, error)

	// Logout logs out a user by invalidating their session or deleting the stored refresh token.
	// This ensures the user can no longer refresh their access token.
	Logout(ctx context.Context, userID string,token string) error
	
	// RefreshToken validates the provided refresh token, invalidates it,
	// and issues a new access token and a new refresh token.
	// It returns the new access token, new refresh token, access token duration, and an error.
	RefreshToken(ctx context.Context, refreshToken string) ( *string, time.Duration, error)




	// OAuthLogin handles login/registration via an external OAuth2 provider.
	OAuthLogin(ctx context.Context, oauthUser *models.User) (*models.LoginResult, error)
}
