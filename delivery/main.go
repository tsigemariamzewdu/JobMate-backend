package main

import (
	"context"
	"log"
	"os"


	"github.com/tsigemariamzewdu/JobMate-backend/delivery/controllers"
	"github.com/tsigemariamzewdu/JobMate-backend/delivery/routers"
	authinfra "github.com/tsigemariamzewdu/JobMate-backend/infrastructure/auth"
	config "github.com/tsigemariamzewdu/JobMate-backend/infrastructure/config"
	mongoclient "github.com/tsigemariamzewdu/JobMate-backend/infrastructure/db/mongo"
	svc "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/services"
	"github.com/tsigemariamzewdu/JobMate-backend/repositories"
	repo "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/repositories"
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
	var otpRepo repo.IOTPRepository = repositories.NewOTPRepository(db)

	// Initialize services
	phoneValidator := &authinfra.PhoneValidatorImpl{}
	otpSender, err := authinfra.NewOTPSenderFromEnv(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize OTP sender: %v", err)
	}
	var otpSenderTyped svc.IOTPSender = otpSender

	// Initialize use case
	authUsecase := usecases.NewAuthUsecase(otpRepo, phoneValidator, otpSenderTyped)

	// Initialize controller
	authController := controllers.NewAuthController(authUsecase)

	// Setup router (add more controllers as you add features)
	router := routers.SetupRouter(authController)

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
