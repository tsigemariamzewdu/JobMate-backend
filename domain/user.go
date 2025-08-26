package models

import "time"

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
	PasswordHash      *string
	YearsExperience   *int
	CareerInterests   *string
	CareerGoals       *string
	ProfilePicture    *string
	RefreshToken   	  *string
	AccessToken       *string
	// LastActiveAt      time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
