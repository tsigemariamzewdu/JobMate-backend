package interfaces

import (
	"context"

	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"
)

type FeedbackRepository interface {
	Create(ctx context.Context, f *models.CVFeedback) (string, error)
}
