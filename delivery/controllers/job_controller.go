package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tsigemariamzewdu/JobMate-backend/delivery/dto"
	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"
	"github.com/tsigemariamzewdu/JobMate-backend/infrastructure/ai"
	"github.com/tsigemariamzewdu/JobMate-backend/repositories"
	"github.com/tsigemariamzewdu/JobMate-backend/usecases"
)

type JobController struct {
	JobUsecase  *usecases.JobUsecase
	JobChatRepo *repositories.JobChatRepository
	GroqClient  *ai.GroqClient
}

func NewJobController(jobUsecase *usecases.JobUsecase, jobChatRepo *repositories.JobChatRepository, groqClient *ai.GroqClient) *JobController {
	return &JobController{
		JobUsecase:  jobUsecase,
		JobChatRepo: jobChatRepo,
		GroqClient:  groqClient,
	}
}

func (jc *JobController) SuggestJobs(c *gin.Context) {
	var req dto.JobSuggestionRequest
	chatID := c.Query("chat_id")
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	// input validation
	if req.UserID == "" || req.Field == "" || req.LookingFor == "" || (req.LookingFor != "local" && req.LookingFor != "remote" && req.LookingFor != "freelance") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid field(s) in request"})
		return
	}
	if req.Language != "en" && req.Language != "am" {
		req.Language = "en"
	}

	// convert DTO to domain model
	domainReq := models.JobSuggestionRequest{
		LookingFor: req.LookingFor,
		Field:      req.Field,
		Skills:     req.Skills,
		Experience: req.Experience,
		Language:   req.Language,
	}

	// retrieve previous chat if chat_id is provided
	var chatMsgs []models.JobChatMessage
	if chatID != "" {
		prevChat, err := jc.JobChatRepo.GetJobChatByID(c.Request.Context(), chatID)
		if err == nil && prevChat != nil {
			chatMsgs = append(chatMsgs, prevChat.Messages...)
		}
	}
	// append new chat history from request
	for _, m := range req.ChatHistory {
		chatMsgs = append(chatMsgs, models.JobChatMessage{
			Role:      m.Role,
			Message:   m.Message,
			Timestamp: time.Now(),
		})
	}

	jobs, aiResp, msg, newChatID, err := jc.JobUsecase.SuggestJobs(c.Request.Context(), req.UserID, domainReq, chatMsgs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	// convert domain jobs to DTOs
	var jobDTOs []dto.JobDTO
	for _, job := range jobs {
		jobDTOs = append(jobDTOs, dto.JobDTO{
			Title:        job.Title,
			Company:      job.Company,
			Location:     job.Location,
			Requirements: job.Requirements,
			Type:         job.Type,
			Source:       job.Source,
			Link:         job.Link,
			Language:     job.Language,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"jobs":       jobDTOs,
		"ai_message": aiResp,
		"message":    msg,
		"chat_id":    newChatID,
	})
}
