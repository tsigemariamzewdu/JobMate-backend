package repositories

import (
	"context"

	repo "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/repositories"
	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type feedbackRepository struct {
	collection *mongo.Collection
}

func NewFeedbackRepository(db *mongo.Database) repo.FeedbackRepository {
	return &feedbackRepository{
		collection: db.Collection("cv_feedback_sessions"),
	}
}

func (r *feedbackRepository) Create(ctx context.Context, f *models.CVFeedback) (string, error) {
	doc := map[string]interface{}{
		"user_id":                 f.UserID,
		"cv_id":                   f.CVID,
		"strengths":               f.Strengths,
		"weaknesses":              f.Weaknesses,
		"improvement_suggestions": f.ImprovementSuggestions,
		"generated_at":            f.GeneratedAt,
	}

	res, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return "", err
	}

	id := res.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}
