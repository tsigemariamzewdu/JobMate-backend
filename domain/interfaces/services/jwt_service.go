package interfaces

import (
	"time"
)

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

	GenerateVerificationToken(userID string) (string, error)
}
