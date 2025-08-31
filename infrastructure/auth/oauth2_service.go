package auth

// imports
import (
	"context"
	"fmt"

	"github.com/tsigemariamzewdu/JobMate-backend/domain/models"
	"github.com/tsigemariamzewdu/JobMate-backend/infrastructure/auth/providers"
	svc "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/services"
)

type oauth2Service struct {
	providers map[string]svc.IOAuth2Provider
}

// creates a new OAuth2 service with the given providers
func NewOAuth2Service(providerConfigs map[string]svc.OAuth2ProviderConfig) (svc.IOAuth2Service, error) {
	service := &oauth2Service{
		providers: make(map[string]svc.IOAuth2Provider),
	}

	for name, config := range providerConfigs {
		
		var provider svc.IOAuth2Provider
		
		switch name {
		case "google":
			provider = providers.NewGoogleProvider(config)
		default:
			return nil, fmt.Errorf("unsupported OAuth2 provider: %s", name)
		}
		
		service.providers[name] = provider
	}

	return service, nil
}

func (o2serv *oauth2Service) SupportedProviders() []string {
	
	providers := make([]string, 0, len(o2serv.providers))
	for name := range o2serv.providers {
		providers = append(providers, name)
	}
	return providers
}

func (o2serv *oauth2Service) GetAuthorizationURL(provider string, state string) (string, error) {
	
	p, ok := o2serv.providers[provider]
	if !ok {
		return "", fmt.Errorf("provider %s is not supported", provider)
	}
	return p.GetAuthorizationURL(state), nil
}

func (o2serv *oauth2Service) Authenticate(ctx context.Context, provider string, code string) (*models.User, error) {
	
	p, ok := o2serv.providers[provider]
	if !ok {
		return nil, fmt.Errorf("provider %s is not supported", provider)
	}
	return p.Authenticate(ctx, code)
}