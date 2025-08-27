package domain

import "errors"

var (
	// User-related errors
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidInput = errors.New("invalid input")

	// Token-related errors
	ErrTokenVerificationFailed = errors.New("token verification failed")
	ErrTokenGenerationFailed   = errors.New("token generation failed")
	ErrTokenUsed               = errors.New("refresh token has already been used")

	// Database-related errors
	ErrDatabaseOperationFailed = errors.New("database operation failed")
)
