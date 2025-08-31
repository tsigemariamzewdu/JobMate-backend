package interfaces

import (
	"context"

	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"
)

type IAuthRepository interface {
	// CreateUser saves a new user to the database.
	CreateUser(c context.Context, user *models.User) error

	// SaveRefreshToken securely stores a new refresh token hash in the database.
	SaveRefreshToken(ctx context.Context, userID string, refreshToken string) error

	// Finds a refresh token by its hash and, if valid, marks it as revoked.
	FindAndInvalidate(ctx context.Context, userID string, refreshToken string) error

	// FindRefreshToken finds a refresh token by its hash without invalidating it.
	FindRefreshToken(ctx context.Context, refreshToken string) (*models.RefreshToken, error)

	// CountByEmail returns the number of users with the given email.
	CountByEmail(c context.Context, email string) (int64, error)

	// CountByPhone returns the number of users with the given phone number.
	CountByPhone(c context.Context, phone string) (int64, error)

	// FindByEmail retrieves a user by their email address.
	FindByEmail(c context.Context, email string) (*models.User, error)

	// FindByPhone retrieves a user by their phone number.
	FindByPhone(c context.Context, phone string) (*models.User, error)

	// FindByID retrieves a user by their ID.
	FindByID(c context.Context, id string) (*models.User, error)

	// UpdateUser updates an existing user in the database.
	UpdateUser(c context.Context, user *models.User) error

    // UpdateTokens updates the access and refresh tokens for a user.
	UpdateTokens(c context.Context, userID string, accessToken string, refreshToken string) error

	// IsEmailVerified checks if the user's email is verified.
	// IsEmailVerified(c context.Context, userID string) (bool, error)
}
