package interfaces

import (
	"context"
	"mime/multipart"

	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"
)

type ICVUsecase interface {
	Upload(ctx context.Context, userID string, rawText string, file *multipart.FileHeader) (*models.CV, error)

	Analyze(ctx context.Context, cvID string) (*models.AISuggestions, error)
}
