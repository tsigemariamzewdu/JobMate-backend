package infrastructure

import (
	
	interfaces "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/services"
	"golang.org/x/oauth2/google"
)

func BuildProviderConfigs() (map[string]interfaces.OAuth2ProviderConfig, error) {
	configs, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	return map[string]interfaces.OAuth2ProviderConfig{
		"google": {
			ClientID:     configs.GoogleClientID,
			ClientSecret: configs.GoogleClientSecret,
			RedirectURL:  configs.BaseURL + configs.GoogleRedirectURL,
			Scopes:       []string{"profile", "email"},
			Endpoint:     google.Endpoint,
		},
	}, nil
}
