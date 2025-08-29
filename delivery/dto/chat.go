package dto

import (
	"time"

	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"
)

type ChatRequest struct {
  UserID  string `json:"user_id" binding:"required"`
  Message string `json:"message" binding:"required"`
  IsFromUser bool `json:"is_from_user"`
}

type ChatResponse struct {
  ID         string `json:"id"`        
  UserID     string `json:"user_id"`   
  Message    string `json:"message"`   
  IsFromUser bool   `json:"is_from_user"`
  MessageType string `json:"message_type,omitempty"`
  CreatedAt  time.Time `json:"created_at"` 
  Intent     string `json:"intent,omitempty"`
  Context    map[string]interface{} `json:"context,omitempty"`
}

// AIMessageDTO represents a message structure for AI models for DTO layer
type AIMessageDTO struct {
  Role    string `json:"role"`
  Content string `json:"content"`
}

// GroqAPIRequest represents the request body for the Groq API
type GroqAPIRequest struct {
  Messages    []AIMessageDTO `json:"messages"` // Use the DTO version
  Model       string             `json:"model"`
  Temperature float32            `json:"temperature"`
  MaxTokens   int                `json:"max_tokens"`
  Stream      bool               `json:"stream"`
}

// GroqAPIResponse represents the response body from the Groq API
type GroqAPIResponse struct {
  Choices []struct {
    Message struct {
      Content string `json:"content"`
      Role    string `json:"role"`
    } `json:"message"`
    FinishReason string `json:"finish_reason"`
    Index        int    `json:"index"`
  } `json:"choices"`
  Created int    `json:"created"`
  ID      string `json:"id"`
  Model   string `json:"model"`
  Object  string `json:"object"`
  Usage   struct {
    CompletionTokens int `json:"completion_tokens"`
    PromptTokens     int `json:"prompt_tokens"`
    TotalTokens      int `json:"total_tokens"`
  } `json:"usage"`
}

func ToChatResponse(conv *models.UserConversation) *ChatResponse {
  return &ChatResponse{
    ID:        conv.ConversationID,
    UserID:    conv.UserID,
    Message:   conv.Message,
    IsFromUser: conv.IsFromUser,
    MessageType: conv.MessageType,
    CreatedAt: conv.CreatedAt,
    Intent:    conv.Intent,
    Context:   conv.Context,
  }
}
