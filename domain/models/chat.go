 package models

import (
  "time"
)

// UserConversation represents a single message in the chat history

type UserConversation struct {
  ConversationID string
  UserID         string
  Message    string
  IsFromUser     bool
  MessageType    string // e.g., 'text', 'quick_reply'
  Intent         string // Detected intent: cv_review, job_search, interview_practice, career_advice
  Context        map[string]interface{} // Conversation context for continuity (e.g., interview_question_index)
  CreatedAt      time.Time
}

// AIMessage represents a message structure for AI models (e.g., Groq, OpenAI)
type AIMessage struct {
  Role    string
  Content string
}
