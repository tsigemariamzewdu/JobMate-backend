package usecases

import (
	"context"
	"fmt"
	"log"

	"mime/multipart"
	"time"

	"github.com/tsigemariamzewdu/JobMate-backend/domain"
	repo "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/repositories"

	service "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/services"

	usecase "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/usecases"

	model "github.com/tsigemariamzewdu/JobMate-backend/domain/models"
)

type CVUsecase struct {
	cvRepo        repo.CVRepository
	feedbackRepo  repo.FeedbackRepository
	skillGapRepo  repo.SkillGapRepository
	aiService     service.AISuggestionService
	textExtractor service.TextExtractor
	timeout       time.Duration
}

func NewCVUsecase(
	cvRepo repo.CVRepository,
	feedbackRepo repo.FeedbackRepository,
	skillGapRepo repo.SkillGapRepository,
	aiService service.AISuggestionService,
	textExtractor service.TextExtractor,
	timeout time.Duration,
) usecase.ICVUsecase {
	return &CVUsecase{
		cvRepo:        cvRepo,
		feedbackRepo:  feedbackRepo,
		skillGapRepo:  skillGapRepo,
		aiService:     aiService,
		textExtractor: textExtractor,
		timeout:       timeout,
	}
}

func (uc *CVUsecase) Upload(ctx context.Context, userID string, rawText string, file *multipart.FileHeader) (*model.CV, error) {

	c, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	if rawText == "" && file != nil {
		text, err := uc.textExtractor.Extract(file)
		if err != nil {
			return nil, fmt.Errorf("failed to extract text from file: %w", err)
		}

		rawText = text
	}

	cv := &model.CV{
		UserID:       userID,
		FileName:     "",
		OriginalText: rawText,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if file != nil {
		cv.FileName = file.Filename
	}

	id, err := uc.cvRepo.Create(c, cv)
	if err != nil {
		return nil, fmt.Errorf("failed to create CV in repository: %w", err)
	}
	cv.ID = id
	return cv, nil
}

func (uc *CVUsecase) Analyze(ctx context.Context, cvID string) (*model.AISuggestions, error) {
	c, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	cv, err := uc.cvRepo.GetByID(c, cvID)
	if err != nil {
		return nil, err
	}

	// Generate AI suggestions
	suggestions, err := uc.aiService.Analyze(c, cv.OriginalText)
	if err != nil {
		return nil, err
	}

	// Update CV
	cv.ExtractedSkills = suggestions.CVs.ExtractedSkills
	cv.ExtractedExperience = suggestions.CVs.ExtractedExperience
	cv.ExtractedEducation = suggestions.CVs.ExtractedEducation
	cv.Summary = suggestions.CVs.Summary
	cv.UpdatedAt = time.Now()

	if err := uc.cvRepo.Update(c, cv); err != nil {
		return nil, domain.ErrCVUpdateFailed
	}

	// Save feedback
	feedback := &model.CVFeedback{
		UserID:                 cv.UserID,
		CVID:                   cv.ID,
		Strengths:              suggestions.CVFeedback.Strengths,
		Weaknesses:             suggestions.CVFeedback.Weaknesses,
		ImprovementSuggestions: suggestions.CVFeedback.ImprovementSuggestions,
		GeneratedAt:            time.Now(),
	}

	if _, err := uc.feedbackRepo.Create(c, feedback); err != nil {
		log.Printf("failed to save CV feedback: %v", err)
	}

	// Save skill gaps
	var gaps []*model.SkillGap
	for _, g := range suggestions.SkillGaps {
		gaps = append(gaps, &model.SkillGap{
			UserID:                 cv.UserID,
			SkillName:              g.SkillName,
			CurrentLevel:           g.CurrentLevel,
			RecommendedLevel:       g.RecommendedLevel,
			Importance:             model.Importance(g.Importance),
			ImprovementSuggestions: g.ImprovementSuggestions,
			CreatedAt:              time.Now(),
			UpdatedAt:              time.Now(),
		})
	}

	if len(gaps) > 0 {
		if err := uc.skillGapRepo.CreateMany(c, gaps); err != nil {
			log.Printf("failed to save skill gaps: %v", err)
		}
	}

	return suggestions, nil
}
