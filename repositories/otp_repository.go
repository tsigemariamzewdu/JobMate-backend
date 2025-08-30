package repositories

import (
	"context"
	"time"

	repo "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/repositories"
	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (r *OTPRepositoryImpl) GetLatestCodeByEmail(ctx context.Context, email string) (*models.UserVerificationCode, error) {
	filter := bson.M{"email": email, "used": false, "expires_at": bson.M{"$gt": time.Now()}}
	opts := options.FindOne().SetSort(bson.M{"created_at": -1})
	
	var result bson.M
	err := r.otpCollection.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	
	return r.mapToUserVerificationCode(result), nil
}
func (r *OTPRepositoryImpl) GetLatestCodeByPhone(ctx context.Context, phone string) (*models.UserVerificationCode, error) {
	filter := bson.M{"phone": phone, "used": false, "expires_at": bson.M{"$gt": time.Now()}}
	opts := options.FindOne().SetSort(bson.M{"created_at": -1})
	
	var result bson.M
	err := r.otpCollection.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	
	return r.mapToUserVerificationCode(result), nil
}
func (r *OTPRepositoryImpl) MarkCodeAsUsed(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"used": true}}
	
	_, err = r.otpCollection.UpdateOne(ctx, filter, update)
	return err
}
func (r *OTPRepositoryImpl) DeleteExpiredCodes(ctx context.Context) error {
	filter := bson.M{"expires_at": bson.M{"$lt": time.Now()}}
	_, err := r.otpCollection.DeleteMany(ctx, filter)
	return err
}

// Helper method to map MongoDB document to UserVerificationCode model
func (r *OTPRepositoryImpl) mapToUserVerificationCode(doc bson.M) *models.UserVerificationCode {
	var userID *string
	if doc["user_id"] != nil {
		uid := doc["user_id"].(string)
		userID = &uid
	}
	
	var phone *string
	if doc["phone"] != nil {
		p := doc["phone"].(string)
		phone = &p
	}
	
	var email *string
	if doc["email"] != nil {
		e := doc["email"].(string)
		email = &e
	}

	// Handle primitive.DateTime conversion for expires_at
	var expiresAt time.Time
	if dt, ok := doc["expires_at"].(primitive.DateTime); ok {
		expiresAt = dt.Time()
	} else if t, ok := doc["expires_at"].(time.Time); ok {
		expiresAt = t
	}

	// Handle primitive.DateTime conversion for created_at
	var createdAt time.Time
	if dt, ok := doc["created_at"].(primitive.DateTime); ok {
		createdAt = dt.Time()
	} else if t, ok := doc["created_at"].(time.Time); ok {
		createdAt = t
	}
	
	return &models.UserVerificationCode{
		ID:         doc["_id"].(primitive.ObjectID).Hex(),
		UserID:     userID,
		Phone:      phone,
		Email:      email,
		CodeHash:   doc["code"].(string),
		Type:       doc["type"].(string),
		ExpiresAt:  expiresAt,
		Used:       doc["used"].(bool),
		CreatedAt:  createdAt,
	}
}
