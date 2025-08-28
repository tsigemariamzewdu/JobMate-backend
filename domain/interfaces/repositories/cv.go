package interfaces

import (
	"context"

	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"
)

type CVRepository interface {
	
	Create(ctx context.Context, cv *models.CV) (string, error)

	GetByID(ctx context.Context, id string) (*models.CV, error)

	Update(ctx context.Context, cv *models.CV) error
}