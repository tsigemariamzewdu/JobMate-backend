package interfaces

import (
	"context"
	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"
	"time"
)

// OTPRepository defines DB operations for OTP codes
// (user_verification_codes)
type IOTPRepository interface {
	// Store a newly generated verification code
	CreateVerificationCode(ctx context.Context, code *models.UserVerificationCode) error

	//  helpers: limit requests per phone/email/IP
	GetRecentRequestsByPhone(ctx context.Context, phone string, since time.Time) (int, error)
	GetRecentRequestsByEmail(ctx context.Context, email string, since time.Time) (int, error)
	GetRecentRequestsByIP(ctx context.Context, ip string, since time.Time) (int, error)

	// Verification helpers
	GetLatestCodeByEmail(ctx context.Context, email string) (*models.UserVerificationCode, error)
	GetLatestCodeByPhone(ctx context.Context, phone string) (*models.UserVerificationCode, error)

	// Mark the code as used once verified
	MarkCodeAsUsed(ctx context.Context, id string) error

	// Optional cleanup
	DeleteExpiredCodes(ctx context.Context) error
}

