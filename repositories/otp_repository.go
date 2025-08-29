package repositories

import (
	"context"
	"time"

	repo "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/repositories"
	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type OTPRepositoryImpl struct {
	otpCollection  *mongo.Collection
   
	
}

func NewOTPRepository(db *mongo.Database) repo.IOTPRepository {
	return &OTPRepositoryImpl{
		otpCollection:  db.Collection("user_verification_codes"),
       	
	}
}

func (r *OTPRepositoryImpl) CreateVerificationCode(ctx context.Context, code *models.UserVerificationCode) error {
	doc := bson.M{
		"user_id":    code.UserID,
		"phone":      code.Phone,
		"email":    code.Email,
		"code":       code.CodeHash,
		"type":       code.Type,
		"expires_at": code.ExpiresAt,
		"used":       code.Used,
		"created_at": code.CreatedAt,
	}
	_, err := r.otpCollection.InsertOne(ctx, doc)
	return err
}

func (r *OTPRepositoryImpl) GetRecentRequestsByPhone(ctx context.Context, phone string, since time.Time) (int, error) {
	filter := bson.M{"phone": phone, "created_at": bson.M{"$gte": since}}
	count, err := r.otpCollection.CountDocuments(ctx, filter)
	return int(count), err
}
func (r *OTPRepositoryImpl) GetRecentRequestsByEmail(ctx context.Context, email string, since time.Time) (int, error) {
	filter := bson.M{"email": email, "created_at": bson.M{"$gte": since}}
	count, err := r.otpCollection.CountDocuments(ctx, filter)
	return int(count), err
}

func (r *OTPRepositoryImpl) GetRecentRequestsByIP(ctx context.Context, ip string, since time.Time) (int, error) {
	return 0, nil
}


