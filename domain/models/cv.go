package models

import "time"

type Language string

const (
	LanguageAm Language = "am"
	LanguageEn Language = "en"
)

type Importance string

const (
	ImportanceCritical   Importance = "critical"
	ImportanceImportant  Importance = "important"
	ImportanceNiceToHave Importance = "nice_to_have"
)

type CV struct {
	ID                  string
	UserID              string
	FileName            string
	OriginalText        string
	ExtractedSkills     []string
	ExtractedExperience []string
	ExtractedEducation  []string
	Summary             string
	Language            Language
	IsActive            bool
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type CVFeedback struct {
	ID                     string
	SessionID              string
	UserID                 string
	CVID                   string
	Strengths              string
	Weaknesses             string
	ImprovementSuggestions string
	GeneratedAt            time.Time
}

type SkillGap struct {
	ID                     string
	UserID                 string
	SkillName              string
	CurrentLevel           int // 1-5 scale
	RecommendedLevel       int // 1-5 scale
	Importance             Importance
	ImprovementSuggestions string
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

type AISuggestions struct {
	CVs struct {
		ExtractedSkills     []string
		ExtractedExperience []string
		ExtractedEducation  []string
		Summary             string
	}
	CVFeedback struct {
		Strengths              string
		Weaknesses             string
		ImprovementSuggestions string
	}
	SkillGaps []struct {
		SkillName              string
		CurrentLevel           int
		RecommendedLevel       int
		Importance             string
		ImprovementSuggestions string
	}
}
