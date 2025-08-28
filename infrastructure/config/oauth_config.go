package infrastructure

import (
	"github.com/tsigemariamzewdu/JobMate-backend/domain"
	"golang.org/x/oauth2/google"
)

func BuildProviderConfigs() (map[string]domain.OAuth2ProviderConfig, error) {
	configs, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	return map[string]domain.OAuth2ProviderConfig{
		"google": {
			ClientID:     configs.GoogleClientID,
			ClientSecret: configs.GoogleClientSecret,
			RedirectURL:  configs.BaseURL + configs.GoogleRedirectURL,
			Scopes:       []string{"profile", "email"},
			Endpoint:     google.Endpoint,
		},
	}, nil
}
