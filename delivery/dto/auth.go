package dto

import (
	"time"

	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"
)

type RefreshTokenDTO struct {
	ID        string    `json:"id,omitempty"`
	UserID    string    `json:"user_id,omitempty"`
	TokenHash string    `json:"token_hash,omitempty"` 
	IsRevoked bool      `json:"is_revoked,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func FromDomainRefreshToken(token *models.RefreshToken) *RefreshTokenDTO {
	if token == nil {
		return nil
	}
	return &RefreshTokenDTO{
		ID:        token.ID,
		UserID:    token.UserID,
		TokenHash: token.TokenHash,
		IsRevoked: token.IsRevoked,
		ExpiresAt: token.ExpiresAt,
		CreatedAt: token.CreatedAt,
	}
}

type LoginResponseDTO struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresIn    time.Duration `json:"expires_in"`
	User         *UserDTO      `json:"user"` 
}
