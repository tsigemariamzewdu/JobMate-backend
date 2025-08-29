package interfaces

import (
	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"
)

type JobSuggestionService interface {
	SuggestJobs(req models.JobSuggestionRequest) ([]models.Job, string, error)
}
