package interfaces

import "github.com/tsigemariamzewdu/JobMate-backend/domain/models"

type IJobRepository interface {
	GetCuratedJobs(field, lookingFor, experience string, skills []string, language string) ([]models.Job, string, error)
}
