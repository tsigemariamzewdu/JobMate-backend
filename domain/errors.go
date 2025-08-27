package domain

import "errors"


var (
	// Token-related errors
	ErrInvalidInput            = errors.New("invalid input")
	ErrTokenVerificationFailed = errors.New("token verification failed")
	ErrTokenGenerationFailed   = errors.New("token generation failed")
	ErrTokenUsed               = errors.New("refresh token has already been used")

	// User and Database related errors
	ErrUserNotFound            = errors.New("user not found")
	ErrDatabaseOperationFailed = errors.New("database operation failed")
)
