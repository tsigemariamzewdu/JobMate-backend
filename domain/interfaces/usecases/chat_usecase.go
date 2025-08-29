package interfaces

import (
	"context"

	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"
)

type IChatUsecase interface {
	SendMessage(ctx context.Context, userID string, message string) (*models.UserConversation, error)
	GetConversationHistory(ctx context.Context, userID string, limit int64) ([]models.UserConversation, error)
}
