package models

import (
	"time"

	"github.com/tsigemariamzewdu/JobMate-backend/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RefreshTokenModel represents the refresh token document in MongoDB.
type RefreshTokenModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    string             `bson:"user_id"`
	TokenHash string             `bson:"token_hash"`
	IsRevoked bool               `bson:"is_revoked"`
	ExpiresAt time.Time          `bson:"expires_at"`
	CreatedAt time.Time          `bson:"created_at"`
}

type User struct {
	UserID         primitive.ObjectID `bson:"_id,omitempty"`
	FirstName      *string            `bson:"first_name"`
	LastName       *string            `bson:"last_name"`
	ProfilePicture *string            `bson:"profile_picture"`
	IsVerified     bool               `bson:"is_verified"`
	Email          string             `bson:"email"`
	Password       *string            `bson:"password"`
	RefreshToken   *string            `bson:"refresh_token"`
	AccessToken    *string            `bson:"access_token"`
	CreatedAt      time.Time          `bson:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at"`

	Provider string         `bson:"provider"`
}

func UserFromDomain(u domain.User) (*User, error) {
	var objID primitive.ObjectID
	var err error

	if u.UserID != "" {
		objID, err = primitive.ObjectIDFromHex(u.UserID)
		if err != nil {
			return nil, domain.ErrInvalidUserID
		}
	}

	return &User{
		UserID: objID,

		FirstName:      u.FirstName,
		LastName:       u.LastName,
		ProfilePicture: u.ProfilePicture,
		IsVerified:     u.IsVerified,
		Email:          *u.Email,
		Password:       u.Password,
		RefreshToken:   u.RefreshToken,
		AccessToken:    u.AccessToken,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,

		Provider: u.Provider,

	}, nil
}

func (u *User) ToDomain() domain.User {
	return domain.User{
		UserID:         u.UserID.Hex(),
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		ProfilePicture: u.ProfilePicture,
		IsVerified:     u.IsVerified,
		Email:          &u.Email,
		Password:       u.Password,
		RefreshToken:   u.RefreshToken,
		AccessToken:    u.AccessToken,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,

		Provider: u.Provider,
	}
}
