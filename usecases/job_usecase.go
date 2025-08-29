package usecases

import (
	"context"
	"time"

	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"
	"github.com/tsigemariamzewdu/JobMate-backend/infrastructure/ai"
	"github.com/tsigemariamzewdu/JobMate-backend/infrastructure/job_service"
	"github.com/tsigemariamzewdu/JobMate-backend/repositories"
)

type JobUsecase struct {
	JobService  *job_service.JobService
	JobChatRepo *repositories.JobChatRepository
	GroqClient  *ai.GroqClient
}

func NewJobUsecase(jobService *job_service.JobService, jobChatRepo *repositories.JobChatRepository, groqClient *ai.GroqClient) *JobUsecase {
	return &JobUsecase{
		JobService:  jobService,
		JobChatRepo: jobChatRepo,
		GroqClient:  groqClient,
	}
}

// SuggestJobs handles the full job chat flow: fetch jobs, store chat, call AI, return all
func (uc *JobUsecase) SuggestJobs(ctx context.Context, userID string, req models.JobSuggestionRequest, chatMsgs []models.JobChatMessage) (jobs []models.Job, aiMessage string, msg string, chatID string, err error) {
	// Fetch jobs
	jobs, msg, err = uc.JobService.GetCuratedJobs(req.Field, req.LookingFor, req.Experience, req.Skills, req.Language)
	if err != nil {
		return nil, "", "No jobs found for your criteria", "", err
	}

	// Save or update job chat
	query := map[string]any{
		"looking_for": req.LookingFor,
		"field":       req.Field,
		"skills":      req.Skills,
		"experience":  req.Experience,
		"language":    req.Language,
	}
	// Always create a new chat for now (could be updated to upsert by chatID)
	chatID, _ = uc.JobChatRepo.CreateJobChat(ctx, userID, query, jobs, chatMsgs)

	// Prepare context for Groq AI
	var aiMessages []models.AIMessage
	for _, m := range chatMsgs {
		aiMessages = append(aiMessages, models.AIMessage{
			Role:    m.Role,
			Content: m.Message,
		})
	}
	// Add job results as a system message
	if len(jobs) > 0 {
		jobSummary := "Job search results:\n"
		for _, job := range jobs {
			jobSummary += "- " + job.Title + " at " + job.Company + " (" + job.Location + ")\n"
		}
		aiMessages = append(aiMessages, models.AIMessage{
			Role:    "system",
			Content: jobSummary,
		})
	}

	// Call Groq AI
	aiResp, _ := uc.GroqClient.GetChatCompletion(context.Background(), aiMessages)

	// Save AI response to chat
	_ = uc.JobChatRepo.AppendMessage(ctx, chatID, models.JobChatMessage{
		Role:      "assistant",
		Message:   aiResp,
		Timestamp: time.Now(),
	})

	return jobs, aiResp, msg, chatID, nil
}
