package job_service

import (
	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"
)

type JobSuggestionService struct {
	JobRepo JobService
}

func NewJobSuggestionService(jobRepo JobService) *JobSuggestionService {
	return &JobSuggestionService{JobRepo: jobRepo}
}

func (s *JobSuggestionService) SuggestJobs(req models.JobSuggestionRequest) ([]models.Job, string, error) {
	return s.JobRepo.GetCuratedJobs(req.Field, req.LookingFor, req.Experience, req.Skills, req.Language)
}
