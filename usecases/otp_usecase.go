package usecases

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	repo "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/repositories"
	svc "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/services"
	uc "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/usecases"
	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"

	"golang.org/x/crypto/bcrypt"
)

const (
    otpLength         = 6
    otpExpiryMinutes  = 5
    otpRateLimitCount = 3
    otpRateLimitWindow = 10 * time.Minute
)

var (
    ErrRateLimited = errors.New("too many OTP requests, please try again later")
    ErrInvalidPhone = errors.New("invalid phone number")
)

type AuthUsecase struct {
    OTPRepo        repo.IOTPRepository
    PhoneValidator uc.IPhoneValidator
    OTPSender      svc.IOTPSender
}

func NewAuthUsecase(repo repo.IOTPRepository, validator uc.IPhoneValidator, sender svc.IOTPSender) *AuthUsecase {
    return &AuthUsecase{
        OTPRepo: repo,
        PhoneValidator: validator,
        OTPSender: sender,
    }
}

func (u *AuthUsecase) RequestOTP(ctx context.Context, req *models.OTPRequest) error {
    // Normalize and validate phone
    normalizedPhone, err := u.PhoneValidator.Normalize(req.Phone)
    if err != nil {
        return ErrInvalidPhone
    }
    if err := u.PhoneValidator.Validate(normalizedPhone); err != nil {
        return ErrInvalidPhone
    }

    // Rate limiting by phone
    since := time.Now().Add(-otpRateLimitWindow)
    count, err := u.OTPRepo.GetRecentRequestsByPhone(ctx, normalizedPhone, since)
    if err == nil && count >= otpRateLimitCount {
        return ErrRateLimited
    }
    // Rate limiting by IP
    if req.RequestorIP != "" {
        ipCount, err := u.OTPRepo.GetRecentRequestsByIP(ctx, req.RequestorIP, since)
        if err == nil && ipCount >= otpRateLimitCount {
            return ErrRateLimited
        }
    }

    // Generate OTP
    otp, err := generateOTP(otpLength)
    if err != nil {
        return errors.New("failed to generate OTP")
    }
    // Hash OTP
    otpHash, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
    if err != nil {
        return errors.New("failed to hash OTP")
    }
    // Persist verification code
    code := &models.UserVerificationCode{
        Phone:     normalizedPhone,
        CodeHash:  string(otpHash),
        Type:      "registration",
        ExpiresAt: time.Now().Add(otpExpiryMinutes * time.Minute),
        Used:      false,
        CreatedAt: time.Now(),
    }
    if err := u.OTPRepo.CreateVerificationCode(ctx, code); err != nil {
        return errors.New("failed to save verification code")
    }

    // Send SMS (stub: print to log)
    if err := u.OTPSender.SendOTP(normalizedPhone, otp); err != nil {
        // Do not leak info
        return errors.New("failed to send OTP")
    }

    // Always return nil (generic response handled in controller)
    return nil
}

// generateOTP generates a secure random n-digit OTP
func generateOTP(length int) (string, error) {
    var num uint32
    err := binary.Read(rand.Reader, binary.LittleEndian, &num)
    if err != nil {
        return "", err
    }
    min := int32(1)
    for i := 1; i < length; i++ {
        min *= 10
    }
    max := int32(1)
    for i := 0; i < length; i++ {
        max *= 10
    }
    otp := int32(num % uint32(max-min)) + min
    return fmt.Sprintf("%0*d", length, otp), nil
}