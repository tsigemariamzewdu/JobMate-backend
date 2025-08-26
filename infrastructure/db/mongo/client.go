package mongo

import (
	"context"
	"log"
	"time"

	infrastructure "github.com/tsigemariamzewdu/JobMate-backend/infrastructure/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// load configuration
	cfg, err := infrastructure.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// connect to MongoDB using the URI from config
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DBUri))
	if err != nil {
		log.Fatalf("Mongo connection error: %v", err)
	}

	return client
}