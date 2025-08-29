package repositories

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/tsigemariamzewdu/JobMate-backend/domain"
	"github.com/tsigemariamzewdu/JobMate-backend/repositories/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AuthRepository handles MongoDB operations for authentication.
type AuthRepository struct {
	tokensCollection *mongo.Collection
	userCollection   *mongo.Collection
}

// NewAuthRepository creates a new instance of AuthRepository.
func NewAuthRepository(db *mongo.Database) domain.IAuthRepository {
	
	return &AuthRepository{
		userCollection: db.Collection("users"),
		tokensCollection: db.Collection("refresh_tokens"),
	}
}

// CreateUser inserts the specified data into user collection
func (r *AuthRepository) CreateUser(ctx context.Context, user *domain.User) error {
	userModel, err := models.UserFromDomain(*user)
	if err != nil {
		return err
	}
	result, err := r.userCollection.InsertOne(ctx, userModel)
	if err != nil {
		return domain.ErrUserCreationFailed
	}

	objID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return err
	}
	user.UserID = objID.Hex()
	return nil
}

// CountByEmail counts how many entry exist in the collection with the specified email
func (r *AuthRepository) CountByEmail(ctx context.Context, email string) (int64, error) {
	filter := bson.D{{Key: "email", Value: email}}
	return r.userCollection.CountDocuments(ctx, filter)
}

// CountByPhone counts how many entry exist in the collection with the specified phone
func (r *AuthRepository) CountByPhone(ctx context.Context, phone string) (int64, error) {
	filter := bson.D{{Key: "phone", Value: phone}}
	return r.tokensCollection.CountDocuments(ctx, filter)
}

// SaveRefreshToken hashes the token and stores it in the database.
func (r *AuthRepository) SaveRefreshToken(ctx context.Context, userID string, refreshToken string) error {
	hashedToken := hashToken(refreshToken)

	model := models.RefreshTokenModel{
		UserID:    userID,
		TokenHash: hashedToken,
		IsRevoked: false,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour), // 60-day expiry
		CreatedAt: time.Now(),
	}

	_, err := r.userCollection.InsertOne(ctx, model)
	if err != nil {
		return fmt.Errorf("failed to save refresh token: %w", err)
	}
	return nil
}

// FindByEmail query the user collection based on specified email
func (ur *AuthRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	filter := bson.D{{Key: "email", Value: email}}
	result := ur.userCollection.FindOne(ctx, filter)

	if err := consolidateUserError(result.Err()); err != nil {
		return nil, err
	}

	var userModel models.User
	if err := result.Decode(&userModel); err != nil {
		return nil, domain.ErrDecodingDocument
	}

	user := userModel.ToDomain()

	return &user, nil
}

// FindByID retrieves a user by their ID.
func (ur *AuthRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidUserID
	}

	// prepare the filter
	filter := bson.D{{Key: "_id", Value: objID}}
	result := ur.userCollection.FindOne(ctx, filter)

	if err := consolidateUserError(result.Err()); err != nil {
		return nil, err
	}

	var userModel models.User
	if err := result.Decode(&userModel); err != nil {
		return nil, domain.ErrDecodingDocument
	}
	user := userModel.ToDomain()

	return &user, nil
}

func (ur *AuthRepository) FindByPhone(ctx context.Context, phone string) (*domain.User, error) {
	filter := bson.D{{Key: "phone", Value: phone}}
	result := ur.userCollection.FindOne(ctx, filter)

	if err := consolidateUserError(result.Err()); err != nil {
		return nil, err
	}

	var userModel models.User
	if err := result.Decode(&userModel); err != nil {
		return nil, domain.ErrDecodingDocument
	}
	user := userModel.ToDomain()

	return &user, nil
}

// IsEmailVerified query and check if the specified id is verified or not
func (ur *AuthRepository) IsEmailVerified(ctx context.Context, id string) (bool, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, domain.ErrInvalidUserID
	}
	filter := bson.D{{Key: "_id", Value: objID}}
	result := ur.userCollection.FindOne(ctx, filter)

	if err := consolidateUserError(result.Err()); err != nil {
		return false, err
	}

	var userModel models.User
	if err := result.Decode(&userModel); err != nil {
		return false, domain.ErrDecodingDocument
	}
	return userModel.IsVerified, nil
}

// Updates a user completely, it returns ErrInvalidUserID or ErrUserNotFound
func (ur *AuthRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	userData, err := models.UserFromDomain(*user)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: userData.UserID}}
	result, err := ur.userCollection.ReplaceOne(ctx, filter, userData)
	if err != nil {
		return err
	}
	// check if a user is found
	if result.MatchedCount == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

// FindAndInvalidate finds the token by hash and marks it as revoked.
func (r *AuthRepository) FindAndInvalidate(ctx context.Context, userID string, refreshToken string) error {
	hashedToken := hashToken(refreshToken)

	// Build the filter to find the correct token.
	filter := bson.M{
		"user_id":    userID,
		"token_hash": hashedToken,
		"is_revoked": false,
		"expires_at": bson.M{"$gt": time.Now()},
	}

	// Update the token's status.
	update := bson.M{
		"$set": bson.M{"is_revoked": true},
	}

	// Find the token and update its status in a single operation.
	result, err := r.tokensCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("database operation failed: %w", err)
	}
	if result.ModifiedCount == 0 {
		return errors.New("token not found, already revoked, or expired")
	}

	return nil
}

// FindRefreshToken finds a refresh token by its hash without invalidating it.
func (r *AuthRepository) FindRefreshToken(ctx context.Context, refreshToken string) (*domain.RefreshToken, error) {
	hashedToken := hashToken(refreshToken)

	var model models.RefreshTokenModel
	err := r.tokensCollection.FindOne(ctx, bson.M{"token_hash": hashedToken}).Decode(&model)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("token not found")
		}
		return nil, fmt.Errorf("database operation failed: %w", err)
	}

	// Convert the repository model to the domain entity.
	idStr := model.ID.Hex()
	domainEntity := &domain.RefreshToken{
		ID:        &idStr,
		UserID:    &model.UserID,
		TokenHash: &model.TokenHash,
		IsRevoked: model.IsRevoked,
		ExpiresAt: model.ExpiresAt,
		CreatedAt: model.CreatedAt,
	}

	return domainEntity, nil
}

// UpdateTokens updates only the access and refresh tokens for a user
func (ur *AuthRepository) UpdateTokens(ctx context.Context, userID string, accessToken, refreshToken string) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return domain.ErrInvalidUserID
	}
	filter := bson.M{"_id": objID}
	update := bson.M{
		"$set": bson.M{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"updated_at":    time.Now(),
		},
	}

	result, err := ur.tokensCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// hashToken is a private helper function to securely hash the token.
func hashToken(token string) string {
	hasher := sha256.New()
	hasher.Write([]byte(token))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

// consolidateUserError extracts the type of error and returns it
func consolidateUserError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, mongo.ErrNoDocuments) {
		return domain.ErrUserNotFound
	}
	return err
}
