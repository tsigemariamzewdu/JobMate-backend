package usecases

import (
	"context"
	"errors"
	"fmt"

	"time"
	"unicode"

	"github.com/tsigemariamzewdu/JobMate-backend/domain"
	repo "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	AuthRepo        domain.IAuthRepository
	OTPRepo        repo.IOTPRepository      
	PasswordService domain.IPasswordService
	JWTService      domain.IJWTService
	EmailService    domain.IEmailService
	BaseURL         string
	ContextTimeout  time.Duration
}

func NewAuthUsecase(repo domain.IAuthRepository, ps domain.IPasswordService, jw domain.IJWTService, bs string,OTPRepo repo.IOTPRepository , timeout time.Duration) domain.IAuthUsecase {
	return &AuthUsecase{
		AuthRepo:        repo,
		PasswordService: ps,
		JWTService:      jw,
		BaseURL:         bs,
		OTPRepo: OTPRepo,
		ContextTimeout:  timeout,
	}
}

// register usecase

// Register handles user registration, supporting both traditional and OAuth-based flows
func (uc *AuthUsecase) Register(ctx context.Context, input *domain.User, oauthUser *domain.User) (*domain.User, error) {

	var email *string
	if oauthUser != nil {
		email = oauthUser.Email
	} else {
		email = input.Email

		// check password strength (min 8 chars, at least one number and one letter)
		if !validatePasswordStrength(*input.Password) {
			return nil, fmt.Errorf("%w", domain.ErrWeakPassword)
		}
	}

	
	// check if email already exists
	count, err := uc.AuthRepo.CountByEmail(ctx, *email)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrDatabaseOperationFailed, err)
	}
	if count > 0 {
		return nil, fmt.Errorf("%w", domain.ErrEmailAlreadyExists)
	}

	
    // OTP verification (only for normal registration, not OAuth)
if oauthUser == nil {
    if input.OTP == nil {
        return nil, fmt.Errorf("%w", errors.New("input otp is empty "))
    }

    // fetch latest OTP for email - FIXED: Use GetLatestCodeByEmail instead of GetRecentRequestsByEmail
    code, err := uc.OTPRepo.GetLatestCodeByEmail(ctx, *email)
    if err != nil {
        return nil, fmt.Errorf("%w: %v", domain.ErrDatabaseOperationFailed, err)
    }
    if code == nil {
        return nil, fmt.Errorf("%w", errors.New("error getting the latest code by email"))
    }

    if code.Used || time.Now().After(code.ExpiresAt) {
        return nil, fmt.Errorf("%w", domain.ErrOTPExpired)
    }

    if err := bcrypt.CompareHashAndPassword([]byte(code.CodeHash), []byte(*input.OTP)); err != nil {
        return nil, fmt.Errorf("%w", domain.ErrInvalidOTP)
    }

    // mark OTP as used - FIXED: Method name should be MarkCodeAsUsed
    if err := uc.OTPRepo.MarkCodeAsUsed(ctx, code.ID); err != nil {
        return nil, fmt.Errorf("%w: %v", domain.ErrOTPUseFailed, err)
    }
}

	// // check if phone already exists
	// var phone *string
	// if input != nil {
	// 	phone = input.Phone
	// } else if oauthUser != nil {
	// 	phone = oauthUser.Phone
	// }
	// if phone != nil {
	// 	count, err = uc.AuthRepo.CountByPhone(ctx, *phone)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("%w: %v", domain.ErrDatabaseOperationFailed, err)
	// 	}
	// 	if count > 0 {
	// 		return nil, fmt.Errorf("%w", domain.ErrPhoneAlreadyExists)
	// 	}
	// }

	var hashedPassword *string
	if oauthUser == nil {
		hashed, err := uc.PasswordService.HashPassword(*input.Password)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", domain.ErrPasswordHashingFailed, err)
		}
		hashedPassword = &hashed
	}

	// construct user model
	newUser := domain.User{
		FirstName: chooseNonEmpty(get(input, func(u *domain.User) *string { return u.FirstName }), get(oauthUser, func(u *domain.User) *string { return u.FirstName })),
		LastName:  chooseNonEmpty(get(input, func(u *domain.User) *string { return u.LastName }), get(oauthUser, func(u *domain.User) *string { return u.LastName })),

		Email:          email,
		IsVerified: true,
		Password:       hashedPassword,
		ProfilePicture: oauthUserPicture(oauthUser),
		Provider:       oauthUserProvider(oauthUser),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// save user to the database
	err = uc.AuthRepo.CreateUser(ctx, &newUser)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrUserCreationFailed, err)
	}
	


	return &newUser, nil
}

// login usecase

// Login handles user login usecase
func (uc *AuthUsecase) Login(ctx context.Context, input *domain.User) (*domain.LoginResult, error) {

	// find user by email or username
	var user *domain.User
	var err error

	if validateEmail(*input.Email) {
		user, err = uc.AuthRepo.FindByEmail(ctx, *input.Email)
	} else {
		user, err = uc.AuthRepo.FindByPhone(ctx, *input.Phone)
	}

	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidCredentials, err)
	}

	// reject login if registered via OAuth
	if user.Provider != "" {
		return nil, fmt.Errorf("%w", domain.ErrOAuthUserCannotLoginWithPassword)
	}

	// check if email is verified
	isVerified, err := uc.AuthRepo.IsEmailVerified(ctx, user.UserID)
	if err != nil {
		return nil, fmt.Errorf("%w", domain.ErrEmailVerficationFailed)
	}
	if !isVerified {
		return nil, fmt.Errorf("%w", domain.ErrEmailNotVerified)
	}

	// compare passwords
	if user.Password == nil || !uc.PasswordService.ComparePassword(*user.Password, *input.Password) {
		return nil, fmt.Errorf("%w", domain.ErrInvalidCredentials)
	}

	// generate access token (handle nil PreferredLanguage)
	lang := "en"
	if user.PreferredLanguage != nil {
		lang = string(*user.PreferredLanguage)
	}
	accessToken, expiresIn, err := uc.JWTService.GenerateAccessToken(user.UserID, lang)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrTokenGenerationFailed, err)
	}

	// generate refresh token
	refreshToken, err := uc.JWTService.GenerateRefreshToken(user.UserID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrTokenGenerationFailed, err)
	}

	user.AccessToken = &accessToken
	user.RefreshToken = &refreshToken

	user.UpdatedAt = time.Now()

	// update the user (save the tokens into database)
	err = uc.AuthRepo.UpdateTokens(ctx, user.UserID, accessToken, refreshToken)
	if err != nil {
		return nil, domain.ErrDatabaseOperationFailed
	}
	result := domain.LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		User:         user,
	}

	return &result, nil
}

// OAuthLogin logs in or registers a user via an OAuth2 provider
func (uc *AuthUsecase) OAuthLogin(ctx context.Context, oauthUser *domain.User) (*domain.LoginResult, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ContextTimeout)
	defer cancel()

	if oauthUser == nil || oauthUser.Email == nil {
		return nil, domain.ErrInvalidOAuthUserData
	}

	// check if the user exists
	user, err := uc.AuthRepo.FindByEmail(ctx, *oauthUser.Email)
	if err != nil {
		// if user doesn't exist, register them
		if err == domain.ErrUserNotFound {
			user, err = uc.Register(ctx, nil, oauthUser)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, domain.ErrDatabaseOperationFailed
		}
	}

	// ensure this user was created via the same provider
	if user.Provider != oauthUser.Provider {
		return nil, fmt.Errorf("%w: expected %s but got %s", domain.ErrOAuthProviderMismatch, user.Provider, oauthUser.Provider)
	}

	// generate access token (handle nil PreferredLanguage)
	lang := "en"
	if user.PreferredLanguage != nil {
		lang = string(*user.PreferredLanguage)
	}
	accessToken, expiresIn, err := uc.JWTService.GenerateAccessToken(user.UserID, lang)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrTokenGenerationFailed, err)
	}

	// generate refresh token
	refreshToken, err := uc.JWTService.GenerateRefreshToken(user.UserID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrTokenGenerationFailed, err)
	}

	// update tokens in db
	err = uc.AuthRepo.UpdateTokens(ctx, user.UserID, accessToken, refreshToken)
	if err != nil {
		return nil, domain.ErrDatabaseOperationFailed
	}

	user.AccessToken = &accessToken
	user.RefreshToken = &refreshToken
	user.UpdatedAt = time.Now()

	return &domain.LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		User:         user,
	}, nil
}

// helper functions
func chooseNonEmpty(primary *string, fallback *string) *string {
	if primary != nil && *primary != "" {
		return primary
	}
	if fallback != nil && *fallback != "" {
		return fallback
	}
	return nil
}

func get(u *domain.User, f func(*domain.User) *string) *string {
	if u == nil {
		return nil
	}
	return f(u)
}

func oauthUserPicture(oauthUser *domain.User) *string {
	if oauthUser == nil || *oauthUser.ProfilePicture == "" {
		return nil
	}
	return oauthUser.ProfilePicture
}

func oauthUserProvider(oauthUser *domain.User) string {
	if oauthUser == nil {
		return ""
	}
	return oauthUser.Provider
}

// logout usecase
func (uc *AuthUsecase) Logout(ctx context.Context, userID string) error {

	//check if empty
	if userID == "" {
		return fmt.Errorf("%w", domain.ErrInvalidUserID)
	}

	//find the users refresh token from db

	user, err := uc.AuthRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("%w", domain.ErrDatabaseOperationFailed)
	}
	//make the refresh token null
	user.RefreshToken = nil

	user.UpdatedAt = time.Now()

	//update the user
	err = uc.AuthRepo.UpdateUser(ctx, user)
	if err != nil {
		return domain.ErrDatabaseOperationFailed
	}

	return nil
}

// refresh token
func (uc *AuthUsecase) RefreshToken(ctx context.Context, userID string) (*string, *string, time.Duration, error) {
	emptyToken := ""

	if userID == "" {
		return &emptyToken, &emptyToken, 0, fmt.Errorf("%w", domain.ErrInvalidInput)
	}

	user, err := uc.AuthRepo.FindByID(ctx, userID)
	if err != nil {
		return &emptyToken, &emptyToken, 0, domain.ErrDatabaseOperationFailed
	}

	if user.RefreshToken == nil {
		return &emptyToken, &emptyToken, 0, domain.ErrInvalidInput
	}

	userIDFromToken, err := uc.JWTService.ValidateRefreshToken(*user.RefreshToken)
	if err != nil {
		return &emptyToken, &emptyToken, 0, domain.ErrTokenVerificationFailed
	}

	lang := "en"
	if user.PreferredLanguage != nil {
		lang = string(*user.PreferredLanguage)
	}
	newAccessToken, expiryTime, err := uc.JWTService.GenerateAccessToken(userIDFromToken, lang)
	if err != nil {
		return &emptyToken, &emptyToken, 0, domain.ErrTokenGenerationFailed
	}

	newRefreshToken, err := uc.JWTService.GenerateRefreshToken(userIDFromToken)
	if err != nil {
		return &emptyToken, &emptyToken, 0, domain.ErrTokenGenerationFailed
	}

	user.AccessToken = &newAccessToken
	user.RefreshToken = &newRefreshToken
	user.UpdatedAt = time.Now()

	err = uc.AuthRepo.UpdateUser(ctx, user)
	if err != nil {
		return &emptyToken, &emptyToken, 0, domain.ErrDatabaseOperationFailed
	}

	return &newAccessToken, &newRefreshToken, time.Duration(expiryTime), nil
}


// function to validate password strength

func validatePasswordStrength(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasLetter := false
	hasNumber := false

	for _, c := range password {
		switch {
		case unicode.IsLetter(c):
			hasLetter = true
		case unicode.IsNumber(c):
			hasNumber = true
		}
	}

	return hasLetter && hasNumber
}


