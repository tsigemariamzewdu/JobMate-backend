package dto

import (
	"time"

	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"
)

type UserDTO struct {
	UserID            string                 `json:"user_id,omitempty"`
	Phone             *string                 `json:"phone,omitempty"`
	Email             *string                 `json:"email,omitempty"`
	FirstName         *string                 `json:"first_name,omitempty"`
	LastName          *string                 `json:"last_name,omitempty"`
	IsVerified        bool                   `json:"is_verified,omitempty"`
	PreferredLanguage *models.PreferredLanguage `json:"preferred_language,omitempty"`
	EducationLevel    *models.EducationLevel    `json:"education_level,omitempty"`
	FieldOfStudy      *string                 `json:"field_of_study,omitempty"`
	YearsExperience   *int                    `json:"years_experience,omitempty"`
	CareerInterests   *string                 `json:"career_interests,omitempty"`
	CareerGoals       *string                 `json:"career_goals,omitempty"`
	ProfilePicture    *string                 `json:"profile_picture,omitempty"`
	Provider          string                 `json:"provider,omitempty"`
	OTP               *string                 `json:"otp,omitempty"`
	CreatedAt         time.Time              `json:"created_at,omitempty"`
	UpdatedAt         time.Time              `json:"updated_at,omitempty"`
}

func ToUserDTO(user *models.User) *UserDTO {
	if user == nil {
		return nil
	}

	return &UserDTO{
		UserID:            user.UserID,
		Phone:             user.Phone,
		Email:             user.Email,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		IsVerified:        user.IsVerified,
		PreferredLanguage: user.PreferredLanguage,
		EducationLevel:    user.EducationLevel,
		FieldOfStudy:      user.FieldOfStudy,
		YearsExperience:   user.YearsExperience,
		CareerInterests:   user.CareerInterests,
		CareerGoals:       user.CareerGoals,
		ProfilePicture:    user.ProfilePicture,
		Provider:          user.Provider,
		OTP:               user.OTP,
		CreatedAt:         user.CreatedAt,
		UpdatedAt:         user.UpdatedAt,
	}
}

type UpdateUserProfile struct {
	FirstName         *string                  `json:"first_name,omitempty"`
	LastName          *string                  `json:"last_name,omitempty"`
	PreferredLanguage *models.PreferredLanguage `json:"preferred_language,omitempty"`
	EducationLevel    *models.EducationLevel    `json:"education_level,omitempty"`
	FieldOfStudy      *string                  `json:"field_of_study,omitempty"`
	YearsExperience   *int                     `json:"years_experience,omitempty"`
	CareerInterests   *string                  `json:"career_interests,omitempty"`
	CareerGoals       *string                  `json:"career_goals,omitempty"`
	ProfilePicture    *string                  `json:"profile_picture,omitempty"`
}

type UserProfileResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	User    *UserDTO `json:"user"`
}
type UserProfileView struct {
	UserID            string  `json:"user_id"`
	FirstName         *string `json:"first_name,omitempty"`
	LastName          *string `json:"last_name,omitempty"`
	PreferredLanguage string  `json:"preferred_language"`
	EducationLevel    string  `json:"education_level"`
	FieldOfStudy      *string `json:"field_of_study,omitempty"`
	YearsExperience   *int    `json:"years_experience,omitempty"`
	CareerInterests   *string `json:"career_interests,omitempty"`
	CareerGoals       *string `json:"career_goals,omitempty"`
	ProfilePicture    *string `json:"profile_picture,omitempty"`
}

func ToUserProfileView(u *models.User) *UserProfileView {
	if u == nil {
		return nil
	}
	return &UserProfileView{
		UserID:            u.UserID,
		FirstName:         u.FirstName,
		LastName:          u.LastName,
		PreferredLanguage: string(*u.PreferredLanguage),
		EducationLevel:    string(*u.EducationLevel),
		FieldOfStudy:      u.FieldOfStudy,
		YearsExperience:   u.YearsExperience,
		CareerInterests:   u.CareerInterests,
		CareerGoals:       u.CareerGoals,
		ProfilePicture:    u.ProfilePicture,
	}
}