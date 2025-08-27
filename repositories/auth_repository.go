package repositories

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/tsigemariamzewdu/JobMate-backend/domain"
	"github.com/tsigemariamzewdu/JobMate-backend/repositories/models"
)

// AuthRepository handles MongoDB operations for authentication.
type AuthRepository struct {
	tokensCollection *mongo.Collection
}

// NewAuthRepository creates a new instance of AuthRepository.
func NewAuthRepository(dbClient *mongo.Client, dbName string) domain.IAuthRepository {
	collection := dbClient.Database(dbName).Collection("refresh_tokens")
	return &AuthRepository{
		tokensCollection: collection,
	}
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

	_, err := r.tokensCollection.InsertOne(ctx, model)
	if err != nil {
		return fmt.Errorf("failed to save refresh token: %w", err)
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
	domainEntity := &domain.RefreshToken{
		ID:        model.ID.Hex(), // Convert MongoDB's ObjectID to string
		UserID:    model.UserID,
		TokenHash: model.TokenHash,
		IsRevoked: model.IsRevoked,
		ExpiresAt: model.ExpiresAt,
		CreatedAt: model.CreatedAt,
	}

	return domainEntity, nil
}

// hashToken is a private helper function to securely hash the token.
func hashToken(token string) string {
	hasher := sha256.New()
	hasher.Write([]byte(token))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}