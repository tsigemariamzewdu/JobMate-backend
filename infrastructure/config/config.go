package infrastructure

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv  string
	AppPort string
	BaseURL string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBDriver   string
	DBUri      string

	JWTSecretKey              string
	JWTExpirationMinutes      int
	RefreshTokenSecret        string
	AccessTokenSecret         string
	RefreshTokenExpirationMin int

	SMTPHost      string
	SMTPPort      int
	SMTPUsername  string
	SMTPPassword  string
	SMTPKey       string
	EmailFrom     string
	EmailFromName string

	// General AI Configuration
	AIApiKey       string 
	AIModelName    string 
	AIApiBaseUrl   string 
	AIProvider     string 
	AITemperature float32 
	
	// Separate config for OpenAI if needed later for CV specific
	OpenAIApiKey string 
	OpenAIModelName string 


	DefaultPageSize int
	MaxPageSize     int

	AllowedOrigins []string
	LogLevel       string
	Timezone       string

	GoogleClientID         string
	GoogleClientSecret     string
	GoogleRedirectURL      string
	GithubClientID         string
	GithubClientSecret     string
	GithubRedirectURL      string
	FacebookClientID       string
	FacebookClientSecret   string
	FacebookRedirectURL    string

	AfricaTalkingUsername string
	AfricaTalkingApiKey   string
	AfricaTalkingSenderId string

	// Twilio fields
	TwilioAccountSID string
	TwilioAuthToken  string
	TwilioFromNumber string

	// JobData
	JobDataApiKey    string
}

// LoadConfig loads config.env from project root (if present) and also supports environment variables.
// If config.env is missing, it falls back to environment variables.
func LoadConfig() (*Config, error) {
	// Tell viper to look for "config.env" in the root folder
	viper.AddConfigPath(".") 
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}


	// Now populate your config struct
	cfg := &Config{
		AppEnv:  viper.GetString("APP_ENV"),
		AppPort: viper.GetString("APP_PORT"),
		BaseURL: viper.GetString("BASE_URL"),

		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     viper.GetString("DB_PORT"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
		DBUri:      viper.GetString("DB_URI"),

		JWTSecretKey:              viper.GetString("JWT_SECRET_KEY"),
		JWTExpirationMinutes:      viper.GetInt("JWT_EXPIRATION_MINUTES"),
		RefreshTokenSecret:        viper.GetString("REFRESH_TOKEN_SECRET"),
		RefreshTokenExpirationMin: viper.GetInt("REFRESH_TOKEN_EXPIRATION_MINUTES"),

		SMTPHost:      viper.GetString("SMTP_HOST"),
		SMTPPort:      viper.GetInt("SMTP_PORT"),
		SMTPUsername:  viper.GetString("SMTP_USERNAME"),
		SMTPPassword:  viper.GetString("SMTP_PASSWORD"),
		SMTPKey:       viper.GetString("SMTP_KEY"),
		EmailFrom:     viper.GetString("EMAIL_FROM"),
		EmailFromName: viper.GetString("EMAIL_FROM_NAME"),

		// General AI Configuration
		AIApiKey: viper.GetString("AI_API_KEY"),        
		AIModelName:  viper.GetString("AI_MODEL_NAME"),  
		AIApiBaseUrl: viper.GetString("AI_API_BASE_URL"), 
		AIProvider:   viper.GetString("AI_PROVIDER"),    
		AITemperature:         float32(viper.GetFloat64("AI_TEMPERATURE")), 
		
		// OpenAI Specific (for CV analysis, if separate)
		OpenAIApiKey: viper.GetString("OPENAI_API_KEY"),
		OpenAIModelName: viper.GetString("OPENAI_MODEL_NAME"),

		DefaultPageSize: viper.GetInt("DEFAULT_PAGE_SIZE"),
		MaxPageSize:     viper.GetInt("MAX_PAGE_SIZE"),

		AllowedOrigins: strings.Split(viper.GetString("ALLOWED_ORIGINS"), ","),
		LogLevel:       viper.GetString("LOG_LEVEL"),
		Timezone:       viper.GetString("TIMEZONE"),

		GoogleClientID:       viper.GetString("GOOGLE_CLIENT_ID"),
		GoogleClientSecret:   viper.GetString("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:    viper.GetString("GOOGLE_REDIRECT_URL"),
		GithubClientID:       viper.GetString("GITHUB_CLIENT_ID"),
		GithubClientSecret:   viper.GetString("GITHUB_CLIENT_SECRET"),
		GithubRedirectURL:    viper.GetString("GITHUB_REDIRECT_URL"),
		FacebookClientID:     viper.GetString("FACEBOOK_CLIENT_ID"),
		FacebookClientSecret: viper.GetString("FACEBOOK_CLIENT_SECRET"),
		FacebookRedirectURL:  viper.GetString("FACEBOOK_REDIRECT_URL"),

		// Twilio (new)
		TwilioAccountSID: viper.GetString("TWILIO_ACCOUNT_SID"),
		TwilioAuthToken:  viper.GetString("TWILIO_AUTH_TOKEN"),
		TwilioFromNumber: viper.GetString("TWILIO_FROM_NUMBER"),

		// Africa's Talking 
		AfricaTalkingUsername: viper.GetString("AFRICASTALKING_USERNAME"),
		AfricaTalkingApiKey:   viper.GetString("AFRICASTALKING_API_KEY"),
		AfricaTalkingSenderId: viper.GetString("AFRICASTALKING_SENDER_ID"),

		// JobData
		JobDataApiKey: viper.GetString("JOBDATA_API_KEY"),
	}

	return cfg, nil
}
