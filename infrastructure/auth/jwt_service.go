package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/tsigemariamzewdu/JobMate-backend/domain"
)

// JWTService implements the domain.IJWTService interface.
type JWTService struct {
	accessSecret 	[]byte
	refreshSecret 	[]byte
}

func NewJWTService(accessSecret string, refreshSecret string) domain.IJWTService {
	return &JWTService{
		accessSecret: 	[]byte(accessSecret),
		refreshSecret: 	[]byte(refreshSecret),
	}
}

// GenerateAccessToken creates a short-lived JWT for user authentication.
func (j *JWTService) GenerateAccessToken(userID string, preferredLanguage string) (string, time.Duration, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"lang": preferredLanguage,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.accessSecret)
	if err != nil {
		return "", 0, err
	}
	return tokenString, 15 * time.Minute, nil
}

// GenerateRefreshToken creates a long-lived token for refreshing the access token.
func (j *JWTService) GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.refreshSecret)
}

// parseToken is a private helper to parse and validate a token with a specific secret.
func (j *JWTService) parseToken(tokenString string, secret []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// ValidateAccessToken verifies the access token's signature and expiration.
func (j *JWTService) ValidateAccessToken(tokenString string) (string, string, error) {
	claims, err := j.parseToken(tokenString, j.accessSecret)
	if err != nil {
		return "", "", err
	}
	sub, ok := claims["sub"].(string)
	if !ok {
		return "", "", errors.New("invalid subject in token")
	}
	lang, ok := claims["lang"].(string)
	if !ok {
		return "", "", errors.New("invalid preferred language")
	}
	return sub, lang, nil
}

// ValidateRefreshToken verifies the refresh token's signature.
func (j *JWTService) ValidateRefreshToken(tokenString string) (string, error) {
	claims, err := j.parseToken(tokenString, j.refreshSecret)
	if err != nil {
		return "", err
	}
	sub, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("invalid subject in token")
	}
	return sub, nil
}