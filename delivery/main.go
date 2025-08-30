package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tsigemariamzewdu/JobMate-backend/delivery/controllers"
	"github.com/tsigemariamzewdu/JobMate-backend/delivery/routes"
	groqpkg "github.com/tsigemariamzewdu/JobMate-backend/infrastructure/ai"
	"github.com/tsigemariamzewdu/JobMate-backend/infrastructure/ai_service"
	authinfra "github.com/tsigemariamzewdu/JobMate-backend/infrastructure/auth"
	config "github.com/tsigemariamzewdu/JobMate-backend/infrastructure/config"
	emailinfra "github.com/tsigemariamzewdu/JobMate-backend/infrastructure/email"
	"github.com/tsigemariamzewdu/JobMate-backend/infrastructure/job_service"

	mongoclient "github.com/tsigemariamzewdu/JobMate-backend/infrastructure/db/mongo"
	// utils "github.com/tsigemariamzewdu/JobMate-backend/infrastructure/util"
	file_parser "github.com/tsigemariamzewdu/JobMate-backend/infrastructure/file_parser"
	"github.com/tsigemariamzewdu/JobMate-backend/repositories"
	"github.com/tsigemariamzewdu/JobMate-backend/usecases"
)

func main() {

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to MongoDB
	client := mongoclient.NewMongoClient()
	db := client.Database(cfg.DBName)
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting MongoDB client: %v", err)
		}
	}()

	// Initialize repositories
	otpRepo := repositories.NewOTPRepository(db)
	authRepo := repositories.NewAuthRepository(db)
	userRepo := repositories.NewUserRepository(db)
	cvRepo := repositories.NewCVRepository(db)
	feedbackRepo := repositories.NewFeedbackRepository(db)
	skillGapRepo := repositories.NewSkillGapRepository(db)
	// use the name conversationRepo because feature branch used it
	conversationRepo := repositories.NewConversationRepository(db)

	providersConfigs, err := config.BuildProviderConfigs()
	if err != nil {
		log.Fatal("error: ", err)
	}

	// Initialize services
	phoneValidator := &authinfra.PhoneValidatorImpl{}
	// email service (feature branch addition)
	emailService := emailinfra.NewSMTPService(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUsername, cfg.SMTPPassword, cfg.EmailFrom)

	otpSender, err := authinfra.NewOTPSenderFromEnv(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize OTP sender: %v", err)
	}
	otpSenderTyped := otpSender
	jwtService := authinfra.NewJWTService(cfg.JWTSecretKey, fmt.Sprint(cfg.JWTExpirationMinutes))
	passwordService := authinfra.NewPasswordService()
	authMiddleware := authinfra.NewAuthMiddleware(jwtService)
	oauthService, err := authinfra.NewOAuth2Service(providersConfigs)
	aiService := ai_service.NewGeminiAISuggestionService("gemini-1.5-flash", cfg.AIApiKey) // to be loaded from config later

	textExtractor := file_parser.NewFileTextExtractor()

	if err != nil {
		log.Fatalf("Failed to initialize OAuth2 service: %v", err)
	}

	// Initialize AI client (avoid alias/variable collision)
	groqClient := groqpkg.NewGroqClient(cfg)

	// Initialize use cases
	// Feature branch expected emailService as an extra arg for NewOTPUsecase
	otpUsecase := usecases.NewOTPUsecase(otpRepo, phoneValidator, otpSenderTyped, emailService)
	// Feature branch expected otpRepo in the auth usecase constructor
	authUsecase := usecases.NewAuthUsecase(authRepo, passwordService, jwtService, cfg.BaseURL, otpRepo, time.Second*10)
	userUsecase := usecases.NewUserUsecase(userRepo, time.Second*10)

	cvUsecase := usecases.NewCVUsecase(cvRepo, feedbackRepo, skillGapRepo, aiService, textExtractor, time.Second*15)
	chatUsecase := usecases.NewChatUsecase(conversationRepo, groqClient, cfg)

	// Job Matching Feature
	jobRepo := job_service.NewJobService(cfg.JobDataApiKey)
	jobChatRepo := repositories.NewJobChatRepository(db)
	// usecase expects job service and jobChatRepo + groq client
	jobUsecase := usecases.NewJobUsecase(jobRepo, jobChatRepo, groqClient)
	jobController := controllers.NewJobController(jobUsecase, jobChatRepo, groqClient)

	// Initialize controllers
	otpController := controllers.NewOtpController(otpUsecase)
	authController := controllers.NewAuthController(authUsecase)
	userController := controllers.NewUserController(userUsecase)
	oauthController := controllers.NewOAuth2Controller(oauthService, authUsecase)
	cvController := controllers.NewCVController(cvUsecase)
	chatController := controllers.NewChatController(chatUsecase)

	// Setup router (add more controllers as you add features)
	router := routes.SetupRouter(authMiddleware, userController, authController, otpController, oauthController, cvController, chatController, jobController)

	// Security: Add CORS and secure headers middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("Referrer-Policy", "no-referrer")
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'")
		c.Next()
	})
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Get port from config or environment variable
	port := cfg.AppPort
	if port == "" {
		port = os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
	}

	// Start server
	log.Printf("Server starting on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
