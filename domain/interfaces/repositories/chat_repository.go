package interfaces

import (
	"context"

	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"
)

// IUserConversationRepository defines DB operations for chat messages
type IUserConversationRepository interface {
  SaveConversationMessage(ctx context.Context, msg *models.UserConversation) error
  GetConversationHistory(ctx context.Context, userID string, limit int64) ([]models.UserConversation, error)
}
