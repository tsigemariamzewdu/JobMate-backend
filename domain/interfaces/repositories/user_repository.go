package interfaces

import(
	"context"
	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"
)
type IUserRepository interface {
	// GetProfile(ctx context.Context) (*User, error)
	UpdateProfile(ctx context.Context, user *models.User) (*models.User, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
}