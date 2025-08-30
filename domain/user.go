package domain

import (
	"context"
	"time"
)

// PreferredLanguage enum
type PreferredLanguage string

const (
	LanguageAmharic PreferredLanguage = "am"
	LanguageEnglish PreferredLanguage = "en"
)

// EducationLevel enum
type EducationLevel string

const (
	EducationHighSchool EducationLevel = "high_school"
	EducationDiploma    EducationLevel = "diploma"
	EducationBachelor   EducationLevel = "bachelor"
	EducationMaster     EducationLevel = "master"
	EducationPhD        EducationLevel = "phd"
	EducationOther      EducationLevel = "other"
)

type User struct {
	UserID            string
	// OauthID           int
	Phone             *string
	
	IsVerified        bool
	Email             *string
	FirstName         *string
	LastName          *string
	PreferredLanguage *PreferredLanguage
	EducationLevel    *EducationLevel
	FieldOfStudy      *string
	Password          *string
	PasswordHash      *string
	YearsExperience   *int
	CareerInterests   *string
	CareerGoals       *string
	ProfilePicture    *string
	RefreshToken   	  *string
	AccessToken       *string
	OTP               *string
	// LastActiveAt      time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time

	Provider          string
}
type IEmailService interface {
	SendEmail(to, subject, body string) error
}

type IUserRepository interface {
	// GetProfile(ctx context.Context) (*User, error)
	UpdateProfile(ctx context.Context, user *User) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
}

type IUserUsecase interface {
	UpdateProfile(ctx context.Context, user *User) (*User, error)
	GetProfile(ctx context.Context, userID string) (*User, error)
}
