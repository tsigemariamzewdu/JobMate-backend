package repositories

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	repo "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/repositories"
	models "github.com/tsigemariamzewdu/JobMate-backend/domain/models"
)

var ErrInvalidUserID = errors.New("invalid user ID")
var ErrUserNotFound = errors.New("user not found")
var ErrRefreshTokenNotFound = errors.New("refresh token not found")
var ErrUserCreationFailed = errors.New("user creation failed")
var ErrDecodingDocument = errors.New("failed to decode document")

// MongoDB-specific persistence models
type User struct {
	UserID                primitive.ObjectID        `bson:"_id,omitempty"`
	FirstName         string                    `bson:"first_name,omitempty"`
	LastName          string                    `bson:"last_name,omitempty"`
	ProfilePicture    string                    `bson:"profile_picture,omitempty"`
	IsVerified        bool                      `bson:"is_verified"`
	Email             string                    `bson:"email"`
	Phone             string                    `bson:"phone,omitempty"`
	Password          string                    `bson:"password,omitempty"` // stored password/hash
	PreferredLanguage models.PreferredLanguage `bson:"preferred_language,omitempty"`
	EducationLevel    models.EducationLevel    `bson:"education_level,omitempty"`
	FieldOfStudy      string                    `bson:"field_of_study,omitempty"`
	YearsExperience   int                       `bson:"years_experience,omitempty"`
	CareerInterests   string                    `bson:"career_interests,omitempty"`
	CareerGoals       string                    `bson:"career_goals,omitempty"`
	RefreshToken      string                    `bson:"refresh_token,omitempty"`
	AccessToken       string                    `bson:"access_token,omitempty"`
	OTP               string                    `bson:"otp,omitempty"`
	CreatedAt         time.Time                 `bson:"created_at"`
	UpdatedAt         time.Time                 `bson:"updated_at"`
	Provider          string                    `bson:"provider,omitempty"`
}

type RefreshTokenModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    string             `bson:"user_id"`
	TokenHash string             `bson:"token_hash"`
	IsRevoked bool               `bson:"is_revoked"`
	ExpiresAt time.Time          `bson:"expires_at"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time         `bson:"updated_at"`
}

// helper converters (domain uses pointer fields) 

func fromPtrString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func toPtrString(s string) *string {
	if s == "" {
		return nil
	}
	v := s
	return &v
}

func fromPtrInt(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

func toPtrInt(i int) *int {
	if i == 0 {
		return nil
	}
	v := i
	return &v
}

func fromPtrPreferredLanguage(p *models.PreferredLanguage) models.PreferredLanguage {
	if p == nil {
		return ""
	}
	return *p
}

func toPtrPreferredLanguage(pl models.PreferredLanguage) *models.PreferredLanguage {
	if pl == "" {
		return nil
	}
	v := pl
	return &v
}

func fromPtrEducationLevel(p *models.EducationLevel) models.EducationLevel {
	if p == nil {
		return ""
	}
	return *p
}

func toPtrEducationLevel(el models.EducationLevel) *models.EducationLevel {
	if el == "" {
		return nil
	}
	v := el
	return &v
}

// --------- mapping functions ----------

// domain -> persistence
func userFromDomain(u *models.User) (*User, error) {
	var objID primitive.ObjectID
	var err error

	if u == nil {
		return nil, errors.New("nil user")
	}

	if u.UserID != "" {
		objID, err = primitive.ObjectIDFromHex(u.UserID)
		if err != nil {
			return nil, ErrInvalidUserID
		}
	}

	// prefer password hash if present in domain, otherwise use Password
	password := fromPtrString(u.PasswordHash)
	if password == "" {
		password = fromPtrString(u.Password)
	}

	return &User{
		UserID:                objID,
		FirstName:         fromPtrString(u.FirstName),
		LastName:          fromPtrString(u.LastName),
		ProfilePicture:    fromPtrString(u.ProfilePicture),
		IsVerified:        u.IsVerified,
		Email:             fromPtrString(u.Email),
		Phone:             fromPtrString(u.Phone),
		Password:          password,
		PreferredLanguage: fromPtrPreferredLanguage(u.PreferredLanguage),
		EducationLevel:    fromPtrEducationLevel(u.EducationLevel),
		FieldOfStudy:      fromPtrString(u.FieldOfStudy),
		YearsExperience:   fromPtrInt(u.YearsExperience),
		CareerInterests:   fromPtrString(u.CareerInterests),
		CareerGoals:       fromPtrString(u.CareerGoals),
		RefreshToken:      fromPtrString(u.RefreshToken),
		AccessToken:       fromPtrString(u.AccessToken),
		OTP:               fromPtrString(u.OTP),
		CreatedAt:         u.CreatedAt,
		UpdatedAt:         u.UpdatedAt,
		Provider:          u.Provider,
	}, nil
}

// persistence -> domain
func (u *User) toDomain() *models.User {
	id := ""
	if !u.UserID.IsZero() {
		id = u.UserID.Hex()
	}

	return &models.User{
		UserID:            id,
		FirstName:         toPtrString(u.FirstName),
		LastName:          toPtrString(u.LastName),
		ProfilePicture:    toPtrString(u.ProfilePicture),
		IsVerified:        u.IsVerified,
		Email:             toPtrString(u.Email),
		Phone:             toPtrString(u.Phone),
		Password:          nil,                    // never set plaintext password from DB
		PasswordHash:      toPtrString(u.Password), // stored value (likely hash)
		PreferredLanguage: toPtrPreferredLanguage(u.PreferredLanguage),
		EducationLevel:    toPtrEducationLevel(u.EducationLevel),
		FieldOfStudy:      toPtrString(u.FieldOfStudy),
		YearsExperience:   toPtrInt(u.YearsExperience),
		CareerInterests:   toPtrString(u.CareerInterests),
		CareerGoals:       toPtrString(u.CareerGoals),
		RefreshToken:      toPtrString(u.RefreshToken),
		AccessToken:       toPtrString(u.AccessToken),
		OTP:               toPtrString(u.OTP),
		CreatedAt:         u.CreatedAt,
		UpdatedAt:         u.UpdatedAt,
		Provider:          u.Provider,
	}
}

func refreshTokenFromDomain(rt *models.RefreshToken) (*RefreshTokenModel, error) {
	var objID primitive.ObjectID
	var err error
	if rt == nil {
		return nil, errors.New("nil refresh token")
	}
	if rt.ID != "" {
		objID, err = primitive.ObjectIDFromHex(rt.ID)
		if err != nil {
			return nil, ErrInvalidUserID
		}
	}

	return &RefreshTokenModel{
		ID:        objID,
		UserID:    rt.UserID,
		TokenHash: rt.TokenHash,
		IsRevoked: rt.IsRevoked,
		ExpiresAt: rt.ExpiresAt,
		CreatedAt: rt.CreatedAt,
	}, nil
}

func (rtm *RefreshTokenModel) toDomain() *models.RefreshToken {
	id := ""
	if !rtm.ID.IsZero() {
		id = rtm.ID.Hex()
	}

	return &models.RefreshToken{
		ID:        id,
		UserID:    rtm.UserID,
		TokenHash: rtm.TokenHash,
		IsRevoked: rtm.IsRevoked,
		ExpiresAt: rtm.ExpiresAt,
		CreatedAt: rtm.CreatedAt,
	}
}


type AuthRepository struct {
	userCollection   *mongo.Collection
	tokensCollection *mongo.Collection
}

func NewAuthRepository(db *mongo.Database) repo.IAuthRepository {
	return &AuthRepository{
		userCollection:   db.Collection("users"),
		tokensCollection: db.Collection("refresh_tokens"),
	}
}

func (r *AuthRepository) CreateUser(ctx context.Context, user *models.User) error {
	userModel, err := userFromDomain(user)
	if err != nil {
		return err
	}

	// If caller didn't set CreatedAt/UpdatedAt we set them here.
	if userModel.CreatedAt.IsZero() {
		userModel.CreatedAt = time.Now()
	}
	userModel.UpdatedAt = time.Now()

	result, err := r.userCollection.InsertOne(ctx, userModel)
	if err != nil {
		return ErrUserCreationFailed
	}

	objID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return ErrUserCreationFailed
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
	return r.userCollection.CountDocuments(ctx, filter)
}

// SaveRefreshToken hashes the token and stores it in the database.
func (r *AuthRepository) SaveRefreshToken(ctx context.Context, userID string, refreshToken string) error {
	hashedToken := hashToken(refreshToken)

	model := RefreshTokenModel{
		UserID:    userID,
		TokenHash: hashedToken,
		IsRevoked: false,
		ExpiresAt: time.Now().Add(24 * time.Hour * 60), // 60-day expiry
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := r.tokensCollection.InsertOne(ctx, model)
	if err != nil {
		return fmt.Errorf("failed to save refresh token: %w", err)
	}
	return nil
}

// FindByEmail query the user collection based on specified email
func (ur *AuthRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	filter := bson.D{{Key: "email", Value: email}}
	result := ur.userCollection.FindOne(ctx, filter)

	if err := consolidateUserError(result.Err()); err != nil {
		return nil, err
	}

	var userModel User
	if err := result.Decode(&userModel); err != nil {
		return nil, ErrDecodingDocument
	}

	user := userModel.toDomain()

	return user, nil
}

// FindByID retrieves a user by their ID.
func (ur *AuthRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidUserID
	}

	filter := bson.D{{Key: "_id", Value: objID}}
	result := ur.userCollection.FindOne(ctx, filter)

	if err := consolidateUserError(result.Err()); err != nil {
		return nil, err
	}

	var userModel User
	if err := result.Decode(&userModel); err != nil {
		return nil, ErrDecodingDocument
	}
	user := userModel.toDomain()

	return user, nil
}

func (ur *AuthRepository) FindByPhone(ctx context.Context, phone string) (*models.User, error) {
	filter := bson.D{{Key: "phone", Value: phone}}
	result := ur.userCollection.FindOne(ctx, filter)

	if err := consolidateUserError(result.Err()); err != nil {
		return nil, err
	}

	var userModel User
	if err := result.Decode(&userModel); err != nil {
		return nil, ErrDecodingDocument
	}
	user := userModel.toDomain()

	return user, nil
}

// Updates a user completely, it returns ErrInvalidUserID or ErrUserNotFound
func (ur *AuthRepository) UpdateUser(ctx context.Context, user *models.User) error {
	if user == nil || user.UserID == "" {
		return ErrInvalidUserID
	}

	objID, err := primitive.ObjectIDFromHex(user.UserID)
	if err != nil {
		return ErrInvalidUserID
	}

	userData, err := userFromDomain(user)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	updateFields := bson.M{
		"first_name":         userData.FirstName,
		"last_name":          userData.LastName,
		"profile_picture":    userData.ProfilePicture,
		"is_verified":        userData.IsVerified,
		"email":              userData.Email,
		"phone":              userData.Phone,
		"password":           userData.Password,
		"preferred_language": userData.PreferredLanguage,
		"education_level":    userData.EducationLevel,
		"field_of_study":     userData.FieldOfStudy,
		"years_experience":   userData.YearsExperience,
		"career_interests":   userData.CareerInterests,
		"career_goals":       userData.CareerGoals,
		"provider":           userData.Provider,
		"otp":                userData.OTP,
		"updated_at":         time.Now(),
	}
	result, err := ur.userCollection.UpdateOne(ctx, filter, bson.M{"$set": updateFields})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrUserNotFound
	}
	return nil
}

// FindAndInvalidate finds the token by hash and marks it as revoked.
func (r *AuthRepository) FindAndInvalidate(ctx context.Context, userID string, refreshToken string) error {
	hashedToken := hashToken(refreshToken)

	filter := bson.M{
		"user_id":    userID,
		"token_hash": hashedToken,
		"is_revoked": false,
		"expires_at": bson.M{"$gt": time.Now()},
	}

	update := bson.M{
		"$set": bson.M{"is_revoked": true, "updated_at": time.Now()},
	}

	result, err := r.tokensCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("database operation failed: %w", err)
	}
	if result.ModifiedCount == 0 {
		return ErrRefreshTokenNotFound
	}

	return nil
}

// FindRefreshToken finds a refresh token by its hash without invalidating it.
func (r *AuthRepository) FindRefreshToken(ctx context.Context, refreshToken string) (*models.RefreshToken, error) {
	hashedToken := hashToken(refreshToken)

	var model RefreshTokenModel
	err := r.tokensCollection.FindOne(ctx, bson.M{"token_hash": hashedToken}).Decode(&model)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrRefreshTokenNotFound
		}
		return nil, fmt.Errorf("database operation failed: %w", err)
	}

	return model.toDomain(), nil
}

// UpdateTokens updates only the access and refresh tokens for a user
func (ur *AuthRepository) UpdateTokens(ctx context.Context, userID string, accessToken, refreshToken string) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return ErrInvalidUserID
	}

	filter := bson.M{"_id": objID}
	update := bson.M{
		"$set": bson.M{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"updated_at":    time.Now(),
		},
	}

	result, err := ur.userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return ErrUserNotFound
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
		return ErrUserNotFound
	}
	return err
}
