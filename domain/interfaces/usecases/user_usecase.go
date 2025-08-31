package interfaces

import(
	"context"
	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"
)
type IUserUsecase interface {
	UpdateProfile(ctx context.Context, user *models.User) (*models.User, error)
	GetProfile(ctx context.Context, userID string) (*models.User, error)
}