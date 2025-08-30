package usecases

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/tsigemariamzewdu/JobMate-backend/domain"
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
    ErrInvalidEmail=errors.New("invaid email address")
    ErrEmailValidationFailed=errors.New("email validation failed")
)

type OTPUsecase struct {
    OTPRepo        repo.IOTPRepository
    PhoneValidator uc.IPhoneValidator
    OTPSender      svc.IOTPSender
   
    EmailService   domain.IEmailService
}

func NewOTPUsecase(repo repo.IOTPRepository, phonevalidator uc.IPhoneValidator, sender svc.IOTPSender,emailService domain.IEmailService) *OTPUsecase {
    return &OTPUsecase{
        OTPRepo:       repo,
        PhoneValidator: phonevalidator,
        OTPSender:     sender,
        
        EmailService: emailService,
    }
}

func (u *OTPUsecase) RequestOTP(ctx context.Context, req *models.OTPRequest) error {
    
    //validate email
   
    // email format validation
	if !validateEmail(*req.Email) {
		return  ErrEmailValidationFailed
	}

    
    // Rate limiting by email
    // since := time.Now().Add(-otpRateLimitWindow)
    // count, err := u.OTPRepo.GetRecentRequestsByEmail(ctx, *req.Email, since)
    // if err == nil && count >= otpRateLimitCount {
    //     return ErrRateLimited
    // }
    
    // // Normalize and validate phone
    // normalizedPhone, err := u.PhoneValidator.Normalize(*req.Phone)
    // if err != nil {
    //     return ErrInvalidPhone
    // }
    // if err := u.PhoneValidator.Validate(normalizedPhone); err != nil {
    //     return ErrInvalidPhone
    // }


    // // Rate limiting by phone
    // since = time.Now().Add(-otpRateLimitWindow)
    // count, err = u.OTPRepo.GetRecentRequestsByPhone(ctx, normalizedPhone, since)
    // if err == nil && count >= otpRateLimitCount {
    //     return ErrRateLimited
    // }
    // // Rate limiting by IP
    // if req.RequestorIP != "" {
    //     ipCount, err := u.OTPRepo.GetRecentRequestsByIP(ctx, req.RequestorIP, since)
    //     if err == nil && ipCount >= otpRateLimitCount {
    //         return ErrRateLimited
    //     }
    // }

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
        Email:     req.Email,
        CodeHash:  string(otpHash),
        Type:      "registration",
        ExpiresAt: time.Now().Add(otpExpiryMinutes * time.Minute),
        Used:      false,
        CreatedAt: time.Now(),
    }
    if err := u.OTPRepo.CreateVerificationCode(ctx, code); err != nil {
        return errors.New("failed to save verification code")
    }
    //send an email with the otp that is generated
    emailBody := generateVerificationEmailBody(otp)
    if err = u.EmailService.SendEmail(*req.Email, "Verify Your Email Address", emailBody); err != nil {
        fmt.Println("email sending failed:", err)
    }

    // // Send SMS (stub: print to log)
    // if err := u.OTPSender.SendOTP(normalizedPhone, otp); err != nil {
    //     // Do not leak info
    //     return errors.New("failed to send OTP")
    // }

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

//function to generate verification email body

func generateVerificationEmailBody(otp string) string {
	return fmt.Sprintf(`
    <html>
  <body style="font-family: Arial, sans-serif; line-height: 1.6; background-color: #f9f9f9; padding: 20px;">
    <div style="max-width: 600px; margin: auto; background: white; border-radius: 8px; padding: 20px; box-shadow: 0 2px 6px rgba(0,0,0,0.1);">
      <h2 style="color: #333;">Welcome!</h2>
      <p>Thanks for signing up. Please use the following One-Time Password (OTP) to verify your email address on the registration page:</p>
      
      <p style="text-align: center; margin: 30px 0;">
        <span style="font-size: 24px; font-weight: bold; letter-spacing: 4px; background: #f3f3f3; padding: 10px 20px; border-radius: 6px; display: inline-block;">
          %s
        </span>
      </p>
      
      <p>This OTP is valid for <strong>10 minutes</strong> and can only be used once.</p>
      <p>If you didn’t request this, you can safely ignore this email.</p>
      <p style="margin-top: 40px;">— The Team</p>
    </div>
  </body>
</html>

  `, otp)
}

//function to validate email

func validateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)

}