package auth

import (
	"errors"
	"regexp"
	"strings"

	uc "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/usecases"
)

type PhoneValidatorImpl struct{}

var _ uc.IPhoneValidator = (*PhoneValidatorImpl)(nil)

// Normalize converts phone to E.164 (+251...) and strips spaces, dashes, etc.
func (v *PhoneValidatorImpl) Normalize(phone string) (string, error) {
    phone = strings.TrimSpace(phone)
    phone = strings.ReplaceAll(phone, " ", "")
    phone = strings.ReplaceAll(phone, "-", "")
    phone = strings.ReplaceAll(phone, "(", "")
    phone = strings.ReplaceAll(phone, ")", "")
    // Handle 09... to +2519...
    if strings.HasPrefix(phone, "09") && len(phone) == 10 {
        phone = "+251" + phone[1:]
    }
    // Handle 2519... to +2519...
    if strings.HasPrefix(phone, "2519") && len(phone) == 12 {
        phone = "+" + phone
    }
    // Already in +2519... format
    if strings.HasPrefix(phone, "+2519") && len(phone) == 13 {
        return phone, nil
    }
    // Fallback: invalid
    return phone, errors.New("invalid phone format")
}

// Validate checks if phone is a valid Ethiopian mobile number
func (v *PhoneValidatorImpl) Validate(phone string) error {
    re := regexp.MustCompile(`^\+2519\d{8}$`)
    if !re.MatchString(phone) {
        return errors.New("invalid Ethiopian phone number")
    }
    return nil
}
