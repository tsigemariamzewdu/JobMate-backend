package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	dto "github.com/tsigemariamzewdu/JobMate-backend/delivery/dto"
	chatUsecase "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/usecases"
)

type ChatController struct {
  ChatUsecase chatUsecase.IChatUsecase
}

func NewChatController(chatUsecase chatUsecase.IChatUsecase) *ChatController {
  return &ChatController{
    ChatUsecase: chatUsecase,
  }
}

func (c *ChatController) SendMessage(gCtx *gin.Context) {
  var request dto.ChatRequest

  if err := gCtx.ShouldBindJSON(&request); err != nil {
    gCtx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
    return
  }

  ctx, cancel := context.WithTimeout(gCtx.Request.Context(), 10*time.Second) // Added timeout context
  defer cancel()

  conversation, err := c.ChatUsecase.SendMessage(ctx, request.UserID, request.Message)
  if err != nil {
    gCtx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
    return
  }

  gCtx.JSON(http.StatusOK, dto.ToChatResponse(conversation))
}

func (c *ChatController) GetConversationHistory(gCtx *gin.Context) {
  userIDStr := gCtx.Query("user_id")
  if userIDStr == "" {
    gCtx.JSON(http.StatusBadRequest, gin.H{"message": "user_id is required"})
    return
  }

  limitStr := gCtx.Query("limit")
  limit := int64(5) // Default limit
  if limitStr != "" {
    parsedLimit, err := strconv.ParseInt(limitStr, 10, 64)
    if err != nil {
      gCtx.JSON(http.StatusBadRequest, gin.H{"message": "invalid limit"})
      return
    }
    limit = parsedLimit
  }

  ctx, cancel := context.WithTimeout(gCtx.Request.Context(), 10*time.Second) // Added timeout context
  defer cancel()

  history, err := c.ChatUsecase.GetConversationHistory(ctx, userIDStr, limit)
  if err != nil {
    gCtx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
    return
  }

  var responses []*dto.ChatResponse
  for _, conv := range history {
    responses = append(responses, dto.ToChatResponse(&conv))
  }

  gCtx.JSON(http.StatusOK, responses)
}
