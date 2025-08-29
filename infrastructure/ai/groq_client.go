package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	dto "github.com/tsigemariamzewdu/JobMate-backend/delivery/dto"
	svc "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/services"
	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"
	config "github.com/tsigemariamzewdu/JobMate-backend/infrastructure/config"
)

type GroqClient struct {
	APIKey     string
	Model      string
	BaseURL    string
	Temperature float32 
	HTTPClient *http.Client
}


var _ svc.IAIClient = (*GroqClient)(nil)

// NewGroqClient creates a new GroqClient instance
func NewGroqClient(cfg *config.Config) *GroqClient {
	return &GroqClient{
		APIKey:     cfg.AIApiKey,
		Model:      cfg.AIModelName,
		BaseURL:    cfg.AIApiBaseUrl, 
		Temperature: cfg.AITemperature, 
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// GetChatCompletion sends a request to the Groq API and returns the AI's response
func (gc *GroqClient) GetChatCompletion(ctx context.Context, domainMessages []models.AIMessage) (string, error) {
	
	var dtoMessages []dto.AIMessageDTO
	for _, msg := range domainMessages {
		dtoMessages = append(dtoMessages, dto.AIMessageDTO{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	requestBody := dto.GroqAPIRequest{
		Messages:    dtoMessages, 
		Model:       gc.Model,
		Temperature: gc.Temperature,
		MaxTokens:   1000,
		Stream:      false,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal Groq API request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/chat/completions", gc.BaseURL), bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create Groq API request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", gc.APIKey))

	resp, err := gc.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to Groq API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse struct {
			Error struct {
				Message string `json:"message"`
				Type    string `json:"type"`
			} `json:"error"`
		}
		if decodeErr := json.NewDecoder(resp.Body).Decode(&errorResponse); decodeErr == nil {
			return "", fmt.Errorf("groq API returned error status %d: %s (Type: %s)", resp.StatusCode, errorResponse.Error.Message, errorResponse.Error.Type)
		}
		return "", fmt.Errorf("groq API returned error status: %s", resp.Status)
	}

	var groqResponse dto.GroqAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&groqResponse); err != nil {
		return "", fmt.Errorf("failed to decode Groq API response: %w", err)
	}

	if len(groqResponse.Choices) > 0 {
		return groqResponse.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("groq API returned no choices")
}