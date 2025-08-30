package domain

import "errors"

var (
	// Token-related errors
	ErrInvalidInput            = errors.New("invalid input")
	ErrTokenVerificationFailed = errors.New("token verification failed")
	ErrTokenGenerationFailed   = errors.New("token generation failed")
	ErrTokenUsed               = errors.New("refresh token has already been used")

	// ─── Repository-Level Errors ──────────────────────────────────────────
	ErrQueryFailed         = errors.New("failed to execute MongoDB query")
	ErrDocumentDecoding    = errors.New("failed to decode MongoDB document")
	ErrCursorFailed        = errors.New("cursor encountered an error during iteration")
	ErrInsertingDocuments  = errors.New("failed to insert document(s)")
	ErrRetrievingDocuments = errors.New("failed to retrieve documents")
	ErrDecodingDocument    = errors.New("failed to decode document")
	ErrUpdatingDocument    = errors.New("failed to update document")
	ErrDeletingDocument    = errors.New("failed to delete document")
	ErrCursorIteration     = errors.New("database cursor iteration error")

	// User related errors
	ErrInvalidUserID                    = errors.New("user_id is required")
	ErrUserNotFound                     = errors.New("user not found")
	ErrDatabaseOperationFailed          = errors.New("database operation failed")
	ErrInvalidToken                     = errors.New("invalid token")
	ErrTokenRevocationFailed            = errors.New("token revocation failed")
	ErrInvalidEmailFormat               = errors.New("invalid email format")
	ErrEmailAlreadyExists               = errors.New("email already exists")
	ErrPhoneAlreadyExists               = errors.New("phone already exists")
	ErrInvalidCredentials               = errors.New("invalid credentials")
	ErrEmailNotVerified                 = errors.New("email not verified")
	ErrPasswordHashingFailed            = errors.New("password hashing failed")
	ErrEmailSendingFailed               = errors.New("email sending failed")
	ErrUserCreationFailed               = errors.New("user creation failed")
	ErrOAuthUserCannotLoginWithPassword = errors.New("OAuth user cannot login with password")
	ErrUserUpdateFailed                 = errors.New("user update failed")
	ErrEmailVerficationFailed           = errors.New("email verification failed")
	ErrPasswordMismatch                 = errors.New("passwords do not match")
	ErrUserVerified                     = errors.New("user already verified")
	ErrGetTokenExpiryFailed             = errors.New("failed to get token expiration time")
	ErrWeakPassword                     = errors.New("password is too weak")
	ErrInvalidOAuthUserData             = errors.New("invalid OAuth user data")
	ErrOAuthProviderMismatch            = errors.New("OAuth provider mismatch for this account")

	// Cv-related errors



	ErrCVUpdateFailed = errors.New("cv update failed")

	ErrCVNotFound  = errors.New("cv not found")
	ErrInvalidCVID = errors.New("invalid cv id")

	//otp realted errors
	ErrMissingOTP=errors.New("otp not found")
	ErrOTPExpired=errors.New("otp is expired")
	ErrInvalidOTP=errors.New("otp is invalid")
	ErrOTPUseFailed=errors.New("otp has failed ")

)
