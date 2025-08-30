package domain

import (
	"context"
	"time"

	"golang.org/x/oauth2"
)

type RefreshToken struct {
	ID        	*string    
	UserID    	*string  
	TokenHash 	*string    
	IsRevoked 	bool      
	ExpiresAt 	time.Time 
	CreatedAt 	time.Time 
}

// IAuthUsecase defines the business logic for authentication operations.
type IAuthUsecase interface {
	// Register registers a new user using provided user data.
	// It should validate the input, hash the password, store the user in the database,
	// and send a verification email.
	Register(ctx context.Context, user, oauthUser *User) (*User, error)
	
	// Login authenticates a user using an identifier (username or email) and password.
	// It should verify credentials, check if the email is verified, and return accesstoken, refreshtoken, user data.
	Login(ctx context.Context, input *User) (*LoginResult, error)

	// Logout logs out a user by invalidating their session or deleting the stored refresh token.
	// This ensures the user can no longer refresh their access token.
	Logout(ctx context.Context, userID string) error
	
	// RefreshToken validates the provided refresh token, invalidates it,
	// and issues a new access token and a new refresh token.
	// It returns the new access token, new refresh token, access token duration, and an error.
	RefreshToken(ctx context.Context, refreshToken string) (*string, *string, time.Duration, error)



	// OAuthLogin handles login/registration via an external OAuth2 provider.
	OAuthLogin(ctx context.Context, oauthUser *User) (*LoginResult, error)
}

// IJWTService defines the contract for our JWT service.
type IJWTService interface {
	// GenerateAccessToken creates a short-lived access token with user claims.
	GenerateAccessToken(userID string, preferredLanguage string) (string, time.Duration, error)

	// GenerateRefreshToken creates a long-lived refresh token.
	GenerateRefreshToken(userID string) (string, error)

	// ValidateAccessToken validates the access token's signature and expiration,
	// and returns the user ID and preferred language from its claims.
	ValidateAccessToken(tokenString string) (userID string, preferredLanguage string, err error)

	// ValidateRefreshToken validates the refresh token's signature.
	ValidateRefreshToken(tokenString string) (userID string, err error)

	GenerateVerificationToken(userID string) (string, error)
}

type IOTPUsecase interface {
	// RequestOTP generates and sends a one-time password (OTP) to the user's phone.
	RequestOTP(ctx context.Context, phone string) error

	// VerifyOTP checks the validity of the provided OTP.
	VerifyOTP(ctx context.Context, phone string, otp string) (bool, error)
}

// IAuthRepository defines the contract for authentication-related database operations.
type IAuthRepository interface {
	// CreateUser saves a new user to the database.
	CreateUser(c context.Context, user *User) error

	// SaveRefreshToken securely stores a new refresh token hash in the database.
	SaveRefreshToken(ctx context.Context, userID string, refreshToken string) error

	// Finds a refresh token by its hash and, if valid, marks it as revoked.
	FindAndInvalidate(ctx context.Context, userID string, refreshToken string) error

	// FindRefreshToken finds a refresh token by its hash without invalidating it.
	FindRefreshToken(ctx context.Context, refreshToken string) (*RefreshToken, error)

	// CountByEmail returns the number of users with the given email.
	CountByEmail(c context.Context, email string) (int64, error)

	// CountByPhone returns the number of users with the given phone number.
	CountByPhone(c context.Context, phone string) (int64, error)

	// FindByEmail retrieves a user by their email address.
	FindByEmail(c context.Context, email string) (*User, error)

	// FindByPhone retrieves a user by their phone number.
	FindByPhone(c context.Context, phone string) (*User, error)

	// FindByID retrieves a user by their ID.
	FindByID(c context.Context, id string) (*User, error)

	// UpdateUser updates an existing user in the database.
	UpdateUser(c context.Context, user *User) error

    // UpdateTokens updates the access and refresh tokens for a user.
	UpdateTokens(c context.Context, userID string, accessToken string, refreshToken string) error

	// IsEmailVerified checks if the user's email is verified.
	// IsEmailVerified(c context.Context, userID string) (bool, error)
}

//PasswordService Interface
type IPasswordService interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, inputPassword string) bool
}

type OAuth2ProviderConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
	Endpoint     oauth2.Endpoint
}

// OAuth2 providers interface
type IOAuth2Provider interface {
	Name() string // provider name
	Authenticate(ctx context.Context, code string) (*User, error)
	GetAuthorizationURL(state string) string
}

type IOAuth2Service interface {
	SupportedProviders() []string
	GetAuthorizationURL(provider string, state string) (string, error)
	Authenticate(ctx context.Context, provider string, code string) (*User, error)
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    time.Duration
	User         *User
}

