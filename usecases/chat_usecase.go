package usecases

import (
	"context"
	"fmt"
	"strings"
	"time"

	repositories "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/repositories"
	chatUsecaseI "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/usecases"
	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"
	groqClient "github.com/tsigemariamzewdu/JobMate-backend/infrastructure/ai"
	config "github.com/tsigemariamzewdu/JobMate-backend/infrastructure/config"
)

// Interview questions (for simplicity, stored directly in usecase for Week 1)
var interviewQuestions = []string{
	"Tell me about yourself.",
	"What are your strengths and weaknesses?",
	"Why do you want to work for this company?",
	"Where do you see yourself in five years?",
	"Do you have any questions for me?",
}

type chatUsecase struct {
	ConversationRepository repositories.IUserConversationRepository
	GroqClient             *groqClient.GroqClient
	AppConfig              *config.Config
}

func NewChatUsecase(convRepo repositories.IUserConversationRepository, groq *groqClient.GroqClient, cfg *config.Config) chatUsecaseI.IChatUsecase {
	return &chatUsecase{
		ConversationRepository: convRepo,
		GroqClient:             groq,
		AppConfig:              cfg,
	}
}

func (u *chatUsecase) SendMessage(ctx context.Context, userID string, message string) (*models.UserConversation, error) {
	// Save user message
	userConversation := &models.UserConversation{
		UserID:      userID,
		IsFromUser:  true,
		Message:     message,
		MessageType: "text", // Default to 'text' as per schema
		CreatedAt:   time.Now(), // Added CreatedAt
	}

	err := u.ConversationRepository.SaveConversationMessage(ctx, userConversation)
	if err != nil {
		return nil, fmt.Errorf("failed to save user conversation: %w", err)
	}

	// Fetch conversation continuity (last ~5 messages)
	history, err := u.ConversationRepository.GetConversationHistory(ctx, userID, 5)
	if err != nil {
		// Log error, but proceed with basic AI interaction if history fails
		fmt.Printf("Error fetching chat history for user %s: %v\n", userID, err)
		history = []models.UserConversation{} // Use empty history
	}

	// Define System Prompt & Build AI Request Messages
	aiMessages := u.buildAIMessages(history, message)

	// Call the AI client (Groq for chat)
	aiRawResponse, err := u.GroqClient.GetChatCompletion(ctx, aiMessages)
	if err != nil {
		// Handle fallback if AI call fails
		fmt.Printf("Error calling AI client for user %s: %v\n", userID, err)
		return u.createFallbackResponse(userID, message, err)
	}

	// Parse Groq Response (Clean Text, Detect Errors) & Intent Tagging
	aiCleanedResponse, intent, newContext := u.parseAIResponseAndTagIntent(history, message, aiRawResponse)

	// Update context for interview mode
	
	if currentMode, ok := u.extractModeFromHistory(history); ok && currentMode == "interview" {
		if idx, ok := u.extractQuestionIndexFromHistory(history); ok {
			if idx < len(interviewQuestions)-1 {
				newContext["question_index"] = idx + 1
				aiCleanedResponse = interviewQuestions[idx+1] // Send next question
				intent = "interview_practice"
			} else {
				newContext["mode"] = "general" // End interview mode
				aiCleanedResponse = "That was the last question! How did you feel about the practice?"
				intent = "interview_practice"
			}
		}
	} else if strings.Contains(strings.ToLower(message), "start interview") {
		// User wants to start interview mode
		intent = "interview_practice"
		aiCleanedResponse = interviewQuestions[0]
		newContext["mode"] = "interview"
		newContext["question_index"] = 0
	} else if aiCleanedResponse == "" {
		aiCleanedResponse = "Can you rephrase that? I'm not sure how to help."
		intent = "general"
	}


	// Construct the AI's response message model for saving
	aiConversation := &models.UserConversation{
		UserID:      userID,
		Message:     aiCleanedResponse,
		IsFromUser:  false,
		MessageType: "text", // Can be dynamic later
		Intent:      intent,
		Context:     newContext,
		CreatedAt:   time.Now(),
	}

	err = u.ConversationRepository.SaveConversationMessage(ctx, aiConversation)
	if err != nil {
		return nil, fmt.Errorf("failed to save AI conversation: %w", err)
	}

	return aiConversation, nil
}

func (u *chatUsecase) GetConversationHistory(ctx context.Context, userID string, limit int64) ([]models.UserConversation, error) {
	return u.ConversationRepository.GetConversationHistory(ctx, userID, limit)
}

// buildAIMessages constructs the array of messages to send to the AI
func (u *chatUsecase) buildAIMessages(history []models.UserConversation, currentMessage string) []models.AIMessage {
	// Define the core system prompt for JobMate
	systemPrompt := "You are JobMate, a helpful, friendly, and supportive career buddy for young job seekers in Ethiopia. Keep answers short, actionable, and culturally relevant. Speak in the same language as the user. Your primary goal is to assist with CV feedback, job matching, and interview practice."

	messages := []models.AIMessage{
		{Role: "system", Content: systemPrompt},
	}

	// Add historical messages for conversation continuity
	for _, conv := range history {
		role := "user"
		if !conv.IsFromUser {
			role = "assistant"
		}
		messages = append(messages, models.AIMessage{Role: role, Content: conv.Message})
	}

	// Add the current user message
	messages = append(messages, models.AIMessage{Role: "user", Content: currentMessage})

	return messages
}

// parseAIResponseAndTagIntent is where Person B's core logic for intent and parsing lives
func (u *chatUsecase) parseAIResponseAndTagIntent(history []models.UserConversation, userMessage, rawAIResponse string) (cleanedResponse string, intent string, newContext map[string]interface{}) {
	cleanedResponse = strings.TrimSpace(rawAIResponse)
	intent = "general" // Default intent

	newContext = make(map[string]interface{})

	// let's keep intent tagging simple (can be AI-driven later)
	lowerUserMessage := strings.ToLower(userMessage)

	if strings.Contains(lowerUserMessage, "interview") || strings.Contains(lowerUserMessage, "practice") {
		intent = "interview_practice"
	} else if strings.Contains(lowerUserMessage, "cv") || strings.Contains(lowerUserMessage, "resume") {
		intent = "cv_review"
	} else if strings.Contains(lowerUserMessage, "job") || strings.Contains(lowerUserMessage, "opportunity") || strings.Contains(lowerUserMessage, "gigs") {
		intent = "job_search"
	} else if strings.Contains(lowerUserMessage, "salary") || strings.Contains(lowerUserMessage, "skills") || strings.Contains(lowerUserMessage, "advice") {
		intent = "career_advice"
	}


	return cleanedResponse, intent, newContext
}

// createFallbackResponse generates a generic error response for the user
func (u *chatUsecase) createFallbackResponse(userID, userMessage string, aiErr error) (*models.UserConversation, error) {
	fallbackMessage := "I apologize, but I'm currently experiencing some technical difficulties. Please try again in a moment."
	if aiErr != nil {
		// Log the actual AI error for debugging purposes (Person A's task to see logs)
		fmt.Printf("Fallback triggered for user %s due to AI error: %v\n", userID, aiErr)
	}

	return &models.UserConversation{
		UserID:      userID,
		Message: fallbackMessage,
		IsFromUser:  false,
		MessageType: "text",
		Intent:      "error_fallback",
		Context:     nil, 
		CreatedAt:   time.Now(),
	}, nil
}

func (u *chatUsecase) extractModeFromHistory(history []models.UserConversation) (string, bool) {
	if len(history) > 0 && history[len(history)-1].Context != nil {
		if mode, ok := history[len(history)-1].Context["mode"].(string); ok {
			return mode, true
		}
	}
	return "", false
}

func (u *chatUsecase) extractQuestionIndexFromHistory(history []models.UserConversation) (int, bool) {
	if len(history) > 0 && history[len(history)-1].Context != nil {
		if idx, ok := history[len(history)-1].Context["question_index"].(float64); ok {
			return int(idx), true 
		}
	}
	return 0, false
}
