package models

type JobSuggestionRequest struct {
	LookingFor string
	Field      string
	Skills     []string
	Experience string
	Language   string
}
