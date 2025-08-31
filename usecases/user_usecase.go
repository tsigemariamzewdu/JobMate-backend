package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/tsigemariamzewdu/JobMate-backend/domain"
	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"
	uc "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/usecases"
	repo "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/repositories"

)

type UserUsecase struct {
	userRepository repo.IUserRepository
	contextTimeout time.Duration
}


func NewUserUsecase(userRepo repo.IUserRepository, timeout time.Duration) uc.IUserUsecase {
	return &UserUsecase{
		userRepository: userRepo,
		contextTimeout: timeout,
	}
}

func (uc *UserUsecase) GetProfile(ctx context.Context, userID string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	user, err := uc.userRepository.GetByID(ctx, userID)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	return user, nil
}

func (uc *UserUsecase) UpdateProfile(ctx context.Context, user *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	existing, err := uc.userRepository.GetByID(ctx, user.UserID)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	if user.FirstName != nil {
		existing.FirstName = user.FirstName
	}
	if user.LastName != nil {
		existing.LastName = user.LastName
	}
	if user.PreferredLanguage != nil {
		existing.PreferredLanguage = user.PreferredLanguage
	}
	if user.EducationLevel != nil {
		existing.EducationLevel = user.EducationLevel
	}
	if user.FieldOfStudy != nil {
		existing.FieldOfStudy = user.FieldOfStudy
	}
	if user.YearsExperience != nil {
		existing.YearsExperience = user.YearsExperience
	}
	if user.CareerInterests != nil {
		existing.CareerInterests = user.CareerInterests
	}
	if user.CareerGoals != nil {
		existing.CareerGoals = user.CareerGoals
	}
	if user.ProfilePicture != nil {
		existing.ProfilePicture = user.ProfilePicture
	}


	updated, err := uc.userRepository.UpdateProfile(ctx, existing)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return updated, nil
}

