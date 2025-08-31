package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/tsigemariamzewdu/JobMate-backend/domain"
	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	repo "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/repositories"

)

type userRepository struct {
	userCollection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) repo.IUserRepository {
	collection := db.Collection("user")
	return &userRepository{userCollection: collection}
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	filter := bson.M{"user_id": id}

	if err := r.userCollection.FindOne(ctx, filter).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil

}

func (r *userRepository) UpdateProfile(ctx context.Context, user *models.User) (*models.User, error) {
	user.UpdatedAt = time.Now()

	filter := bson.M{"user_id": user.UserID}
	update := bson.M{"$set": user}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedUser models.User
	if err := r.userCollection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedUser); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &updatedUser, nil
}
