package interfaces

import (
	"context"
	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"
)

// IAIClient defines the contract for interacting with an AI chat completion service
type IAIClient interface {
	GetChatCompletion(ctx context.Context, messages []models.AIMessage) (string, error)
}
