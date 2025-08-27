package dto

import (
	"github.com/tsigemariamzewdu/JobMate-backend/domain"
)

type UpdateUserProfile struct {
	FirstName         *string                   `json:"first_name,omitempty"`
	LastName          *string                   `json:"last_name,omitempty"`
	PreferredLanguage *domain.PreferredLanguage `json:"preferred_language,omitempty"`
	EducationLevel    *domain.EducationLevel    `json:"education_level,omitempty"`
	FieldOfStudy      *string                   `json:"field_of_study,omitempty"`
	YearsExperience   *int                      `json:"years_experience,omitempty"`
	CareerInterests   *string                   `json:"career_interests,omitempty"`
	CareerGoals       *string                   `json:"career_goals,omitempty"`
	ProfilePicture    *string                   `json:"profile_picture,omitempty"`
}

type UserProfileResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	User    *UserProfileView `json:"user"`
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

func ToUserProfileView(u *domain.User) *UserProfileView {
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
