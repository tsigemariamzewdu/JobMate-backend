package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/tsigemariamzewdu/JobMate-backend/domain"
)

type AuthUseCase struct {
	UserRepo       domain.IUserRepository  
	AuthRepo       domain.IAuthRepository 
	JWTService     domain.IJWTService
	ContextTimeout time.Duration
}

// RefreshToken validates the provided refresh token and issues a new access and refresh token.
func (uc *AuthUseCase) RefreshToken(ctx context.Context, refreshToken string) (*string, *string, time.Duration, error) {
	if refreshToken == "" {
		return nil, nil, 0, fmt.Errorf("refresh token is missing: %w", domain.ErrInvalidInput)
	}

	// Validate the incoming token's signature to get the userID
	userID, err := uc.JWTService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("invalid refresh token: %w", domain.ErrTokenVerificationFailed)
	}

	// Look up and invalidate the token in the database
	err = uc.AuthRepo.FindAndInvalidate(ctx, userID, refreshToken)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("token not found or already used: %w", err)
	}

	// Fetch user data (preferred language) from the main user table
	user, err := uc.UserRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("user not found: %w", domain.ErrUserNotFound)
	}

	// Generate a new pair of tokens

	var lang string
	if user.PreferredLanguage != nil {
		lang = string(*user.PreferredLanguage)
	} else {
		lang = "en"
	}
	newAccessToken, expiryTime, err := uc.JWTService.GenerateAccessToken(user.UserID, lang)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to generate new access token: %w", domain.ErrTokenGenerationFailed)
	}

	newRefreshToken, err := uc.JWTService.GenerateRefreshToken(user.UserID)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to generate new refresh token: %w", domain.ErrTokenGenerationFailed)
	}

	// Store the new token's hash in the database
	err = uc.AuthRepo.SaveRefreshToken(ctx, userID, newRefreshToken)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to save new token: %w", err)
	}

	return &newAccessToken, &newRefreshToken, expiryTime, nil
}