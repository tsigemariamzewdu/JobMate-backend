package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	repo "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/repositories"

	"github.com/tsigemariamzewdu/JobMate-backend/domain"
	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type cvModel struct {
	ID                  primitive.ObjectID `bson:"_id"`
	UserID              string             `bson:"user_id"`
	FileName            string             `bson:"file_name"`
	OriginalText        string             `bson:"original_text"`
	ExtractedSkills     []string           `bson:"extracted_skills"`
	ExtractedExperience []string           `bson:"extracted_experience"`
	ExtractedEducation  []string           `bson:"extracted_education"`
	Summary             string             `bson:"summary"`
	IsActive            bool               `bson:"is_active"`
	CreatedAt           time.Time          `bson:"created_at"`
	UpdatedAt           time.Time          `bson:"updated_at"`
}

func toDomainCV(m cvModel) *models.CV {
	return &models.CV{
		ID:                  m.ID.Hex(),
		UserID:              m.UserID,
		FileName:            m.FileName,
		OriginalText:        m.OriginalText,
		ExtractedSkills:     m.ExtractedSkills,
		ExtractedExperience: m.ExtractedExperience,
		ExtractedEducation:  m.ExtractedEducation,
		Summary:             m.Summary,
		IsActive:            m.IsActive,
		CreatedAt:           m.CreatedAt,
		UpdatedAt:           m.UpdatedAt,
	}
}

func toCVModel(d models.CV) (*cvModel, error) {
	id := primitive.NewObjectID()
	if d.ID != "" {
		var err error
		id, err = primitive.ObjectIDFromHex(d.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid CV ID: %w", err)
		}
	}

	return &cvModel{
		ID:                  id,
		UserID:              d.UserID,
		FileName:            d.FileName,
		OriginalText:        d.OriginalText,
		ExtractedSkills:     d.ExtractedSkills,
		ExtractedExperience: d.ExtractedExperience,
		ExtractedEducation:  d.ExtractedEducation,
		Summary:             d.Summary,
		IsActive:            d.IsActive,
		CreatedAt:           d.CreatedAt,
		UpdatedAt:           d.UpdatedAt,
	}, nil
}

type cvRepository struct {
	collection *mongo.Collection
}

func NewCVRepository(db *mongo.Database) repo.CVRepository {
	return &cvRepository{collection: db.Collection("cvs")}
}

func (r *cvRepository) Create(ctx context.Context, cv *models.CV) (string, error) {
	model, err := toCVModel(*cv)
	if err != nil {
		return "", err
	}

	_, err = r.collection.InsertOne(ctx, model)
	if err != nil {
		return "", fmt.Errorf("failed to insert CV: %w", err)
	}

	return model.ID.Hex(), nil
}
func (r *cvRepository) GetByID(ctx context.Context, id string) (*models.CV, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidCVID
	}

	var model cvModel
	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&model)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrCVNotFound
		}
		return nil, err
	}

	return toDomainCV(model), nil
}

func (r *cvRepository) Update(ctx context.Context, cv *models.CV) error {
	oid, err := primitive.ObjectIDFromHex(cv.ID)
	if err != nil {
		return domain.ErrInvalidCVID
	}

	filter := bson.M{"_id": oid}
	update := bson.M{
		"$set": bson.M{
			"extracted_skills":     cv.ExtractedSkills,
			"extracted_experience": cv.ExtractedExperience,
			"extracted_education":  cv.ExtractedEducation,
			"summary":              cv.Summary,
			"updated_at":           time.Now(),
		},
	}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return domain.ErrCVNotFound
	}
	return nil
}
