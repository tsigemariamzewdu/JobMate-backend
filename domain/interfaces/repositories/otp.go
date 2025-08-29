package interfaces

import (
	"context"
	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"
	"time"
)

// OTPRepository defines DB operations for OTP codes
// (user_verification_codes)
type IOTPRepository interface {
	CreateVerificationCode(ctx context.Context, code *models.UserVerificationCode) error
	GetRecentRequestsByPhone(ctx context.Context, phone string, since time.Time) (int, error)
	GetRecentRequestsByEmail(ctx context.Context, phone string, since time.Time) (int, error)
	GetRecentRequestsByIP(ctx context.Context, ip string, since time.Time) (int, error)
}
