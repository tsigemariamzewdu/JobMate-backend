package repositories

import (
	"context"
	"time"

	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type JobChatRepository struct {
	collection *mongo.Collection
}

func NewJobChatRepository(db *mongo.Database) *JobChatRepository {
	return &JobChatRepository{
		collection: db.Collection("job_chats"),
	}
}

func (r *JobChatRepository) CreateJobChat(ctx context.Context, userID string, query map[string]any, jobResults []models.Job, messages []models.JobChatMessage) (string, error) {
	chat := models.JobChat{
		UserID:         userID,
		Messages:       messages,
		JobSearchQuery: query,
		JobResults:     jobResults,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	result, err := r.collection.InsertOne(ctx, chat)
	if err != nil {
		return "", err
	}
	id := result.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

func (r *JobChatRepository) AppendMessage(ctx context.Context, chatID string, message models.JobChatMessage) error {
	objID, err := primitive.ObjectIDFromHex(chatID)
	if err != nil {
		return err
	}
	update := bson.M{
		"$push": bson.M{"messages": message},
		"$set":  bson.M{"updated_at": time.Now()},
	}
	_, err = r.collection.UpdateByID(ctx, objID, update)
	return err
}

func (r *JobChatRepository) GetJobChatByID(ctx context.Context, chatID string) (*models.JobChat, error) {
	objID, err := primitive.ObjectIDFromHex(chatID)
	if err != nil {
		return nil, err
	}
	var chat models.JobChat
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&chat)
	if err != nil {
		return nil, err
	}
	return &chat, nil
}

func (r *JobChatRepository) GetJobChatsByUserID(ctx context.Context, userID string) ([]*models.JobChat, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	var chats []*models.JobChat
	if err := cursor.All(ctx, &chats); err != nil {
		return nil, err
	}
	return chats, nil
}
