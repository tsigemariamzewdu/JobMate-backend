package domain

import (
	"context"
	"time"
)

// IAuthUsecase defines the business logic for authentication operations.
type IAuthUsecase interface {
	// RefreshToken validates the provided refresh token, invalidates it,
	// and issues a new access token and a new refresh token.
	// It returns the new access token, new refresh token, access token duration, and an error.
	RefreshToken(ctx context.Context, refreshToken string) (*string, *string, time.Duration, error)
}

// IJWTService defines the contract for our JWT service.
type IJWTService interface {
	// GenerateAccessToken creates a short-lived access token with user claims.
	GenerateAccessToken(userID string, preferredLanguage string) (string, time.Duration, error)

	// GenerateRefreshToken creates a long-lived refresh token.
	GenerateRefreshToken(userID string) (string, error)

	// ValidateAccessToken validates the access token's signature and expiration,
	// and returns the user ID and preferred language from its claims.
	ValidateAccessToken(tokenString string) (userID string, preferredLanguage string, err error)

	// ValidateRefreshToken validates the refresh token's signature.
	ValidateRefreshToken(tokenString string) (userID string, err error)
}