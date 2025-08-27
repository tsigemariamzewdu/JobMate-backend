package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/tsigemariamzewdu/JobMate-backend/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepository struct {
	userCollection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) domain.IUserRepository {
	collection := db.Collection("user")
	return &userRepository{userCollection: collection}
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	filter := bson.M{"user_id": id}

	if err := r.userCollection.FindOne(ctx, filter).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil

}

func (r *userRepository) UpdateProfile(ctx context.Context, user *domain.User) (*domain.User, error) {
	user.UpdatedAt = time.Now()

	filter := bson.M{"user_id": user.UserID}
	update := bson.M{"$set": user}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedUser domain.User
	if err := r.userCollection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedUser); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &updatedUser, nil
}

// FindByID returns a user by their ID
func (r *userRepository) FindByID(ctx context.Context, userID string) (*domain.User, error) {
	var user domain.User
	filter := bson.M{"user_id": userID}
	if err := r.userCollection.FindOne(ctx, filter).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}