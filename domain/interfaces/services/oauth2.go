package interfaces


import (
	"context"

	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"

	"golang.org/x/oauth2"
)

type OAuth2ProviderConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
	Endpoint     oauth2.Endpoint
}

// IOAuth2Provider defines the interface for individual OAuth2 providers
type IOAuth2Provider interface {
	Name() string // provider name
	Authenticate(ctx context.Context, code string) (*models.User, error)
	GetAuthorizationURL(state string) string
}

// IOAuth2Service defines the interface for managing multiple OAuth2 providers
type IOAuth2Service interface {
	SupportedProviders() []string
	GetAuthorizationURL(provider string, state string) (string, error)
	Authenticate(ctx context.Context, provider string, code string) (*models.User, error)
}
