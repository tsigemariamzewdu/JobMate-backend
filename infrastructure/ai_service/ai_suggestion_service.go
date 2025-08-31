package ai_service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	svc "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/services"

	model "github.com/tsigemariamzewdu/JobMate-backend/domain/models"
	"google.golang.org/genai"
)

type aiResponse struct {
	CVs struct {
		ExtractedSkills     []string `json:"extracted_skills"`
		ExtractedExperience []string `json:"extracted_experience"`
		ExtractedEducation  []string `json:"extracted_education"`
		Summary             string   `json:"summary"`
	} `json:"cvs"`
	CVFeedback struct {
		Strengths              string `json:"strengths"`
		Weaknesses             string `json:"weaknesses"`
		ImprovementSuggestions string `json:"improvement_suggestions"`
	} `json:"cv_feedback"`
	SkillGaps []struct {
		SkillName              string `json:"skill_name"`
		CurrentLevel           int    `json:"current_level"`
		RecommendedLevel       int    `json:"recommended_level"`
		Importance             string `json:"importance"`
		ImprovementSuggestions string `json:"improvement_suggestions"`
	} `json:"skill_gaps"`
}

type GeminiAISuggestionService struct {
	model  string
	apiKey string
}

func NewGeminiAISuggestionService(model, apiKey string) svc.AISuggestionService {
	if model == "" {
		model = "gemini-1.5-flash"
	}
	return &GeminiAISuggestionService{
		model:  model,
		apiKey: apiKey,
	}
}

func (s *GeminiAISuggestionService) Analyze(ctx context.Context, cvText string) (*model.AISuggestions, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: s.apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	prompt := fmt.Sprintf(`You are a career coach AI. Analyze the following CV text and return **only JSON**, strictly matching this structure. Use empty arrays or empty strings if there is no data:

{
  "cvs": {
    "extracted_skills": ["skill1", "skill2"],
    "extracted_experience": ["experience1"],
    "extracted_education": ["education1"],
    "summary": "Concise professional summary"
  },
  "cv_feedback": {
    "strengths": "Highlight strong points",
    "weaknesses": "Highlight weak points",
    "improvement_suggestions": "Actionable suggestions"
  },
  "skill_gaps": [
    {
      "skill_name": "Name",
      "current_level": 1,
      "recommended_level": 5,
      "importance": "critical",
      "improvement_suggestions": "How to improve"
    }
  ]
}

CV Text:
%s
`, cvText)

	result, err := client.Models.GenerateContent(ctx, s.model, genai.Text(prompt), nil)
	if err != nil {
		return nil, fmt.Errorf("AI generation failed: %w", err)
	}
	resp := strings.TrimSpace(result.Text())
	resp = strings.TrimPrefix(resp, "```json")
	resp = strings.TrimPrefix(resp, "```")
	resp = strings.TrimSuffix(resp, "```")
	resp = strings.TrimSpace(resp)

	var aiResp aiResponse
	if err := json.Unmarshal([]byte(resp), &aiResp); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w\nAI output: %s", err, resp)
	}

	// Map AI response to domain.AISuggestions
	suggestions := &model.AISuggestions{
		CVs: struct {
			ExtractedSkills     []string
			ExtractedExperience []string
			ExtractedEducation  []string
			Summary             string
		}{
			ExtractedSkills:     aiResp.CVs.ExtractedSkills,
			ExtractedExperience: aiResp.CVs.ExtractedExperience,
			ExtractedEducation:  aiResp.CVs.ExtractedEducation,
			Summary:             aiResp.CVs.Summary,
		},
		CVFeedback: struct {
			Strengths              string
			Weaknesses             string
			ImprovementSuggestions string
		}{
			Strengths:              aiResp.CVFeedback.Strengths,
			Weaknesses:             aiResp.CVFeedback.Weaknesses,
			ImprovementSuggestions: aiResp.CVFeedback.ImprovementSuggestions,
		},
	}

	type skillGapType = struct {
		SkillName              string `json:"skill_name"`
		CurrentLevel           int    `json:"current_level"`
		RecommendedLevel       int    `json:"recommended_level"`
		Importance             string `json:"importance"`
		ImprovementSuggestions string `json:"improvement_suggestions"`
	}

	if aiResp.SkillGaps == nil {
		aiResp.SkillGaps = make([]skillGapType, 0)
	} else {
		for _, g := range aiResp.SkillGaps {
			suggestions.SkillGaps = append(suggestions.SkillGaps, struct {
				SkillName              string
				CurrentLevel           int
				RecommendedLevel       int
				Importance             string
				ImprovementSuggestions string
			}{
				SkillName:              g.SkillName,
				CurrentLevel:           g.CurrentLevel,
				RecommendedLevel:       g.RecommendedLevel,
				Importance:             g.Importance,
				ImprovementSuggestions: g.ImprovementSuggestions,
			})
		}
	}

	return suggestions, nil
}
