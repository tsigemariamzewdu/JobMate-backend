package repositories

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	repositories "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/repositories"
	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"
)

type conversationRepository struct {
  Collection *mongo.Collection
}

// userConversationDoc is an internal struct for MongoDB interaction,
type userConversationDoc struct {
  ID             primitive.ObjectID `bson:"_id,omitempty"`
  UserID         string             `bson:"user_id"`
  Message        string             `bson:"message"`
  IsFromUser     bool               `bson:"is_from_user"`
  MessageType    string             `bson:"message_type,omitempty"`
  Intent         string             `bson:"intent,omitempty"`
  Context        map[string]interface{} `bson:"context,omitempty"`
  CreatedAt      time.Time          `bson:"created_at"`
}

func NewConversationRepository(db *mongo.Database) repositories.IUserConversationRepository {
  return &conversationRepository{
    Collection: db.Collection("user_conversations"),
  }
}

func (r *conversationRepository) SaveConversationMessage(ctx context.Context, conversation *models.UserConversation) error {
  doc := userConversationDoc{
    ID:             primitive.NewObjectID(), 
    UserID:         conversation.UserID,
    Message:        conversation.Message,
    IsFromUser:     conversation.IsFromUser,
    MessageType:    conversation.MessageType,
    Intent:         conversation.Intent,
    Context:        conversation.Context,
    CreatedAt:      time.Now(),
  }

  result, err := r.Collection.InsertOne(ctx, doc)
  if err != nil {
    return err
  }

  // Update the ConversationID in the domain model with the generated _id
  if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
    conversation.ConversationID = oid.Hex()
  } else {
    return fmt.Errorf("failed to assert InsertedID to primitive.ObjectID")
  }
  
  return nil
}

func (r *conversationRepository) GetConversationHistory(ctx context.Context, userID string, limit int64) ([]models.UserConversation, error) {
  var docs []userConversationDoc
  var conversations []models.UserConversation

  opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(limit)
  filter := bson.M{"user_id": userID}

  cursor, err := r.Collection.Find(ctx, filter, opts)
  if err != nil {
    return nil, err
  }
  defer cursor.Close(ctx)

  if err = cursor.All(ctx, &docs); err != nil {
    return nil, err
  }

  for _, doc := range docs {
    conversations = append(conversations, models.UserConversation{
      ConversationID: doc.ID.Hex(), 
      UserID:         doc.UserID,
      Message:        doc.Message,
      IsFromUser:     doc.IsFromUser,
      MessageType:    doc.MessageType,
      Intent:         doc.Intent,
      Context:        doc.Context,
      CreatedAt:      doc.CreatedAt,
    })
  }

  // Reverse the order to show oldest first
  for i, j := 0, len(conversations)-1; i < j; i, j = i+1, j-1 {
    conversations[i], conversations[j] = conversations[j], conversations[i]
  }

  return conversations, nil
}
