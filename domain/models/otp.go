package models

import "time"

// OTPRequest represents the input for requesting an OTP
// (business logic only, not for transport)
type OTPRequest struct {
    Phone       *string
    Email       *string
    RequestorIP string // for rate limiting
}

// UserVerificationCode represents a verification code (business concept)
type UserVerificationCode struct {
    ID         string
    UserID     *string
    Phone      *string
    Email      *string
    CodeHash   string // store hashed code
    Type       string // e.g., 'registration', 'password_reset'
    ExpiresAt  time.Time
    Used       bool
    CreatedAt  time.Time
}
