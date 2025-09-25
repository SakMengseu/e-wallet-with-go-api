package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoConfig struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func ConnectDB() (*MongoConfig, error) {
	
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	
	uri := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("MONGO_DB")

	fmt.Println("MONGODB_URI:", uri) // Debugging line

	if uri == "" {
		return nil, fmt.Errorf("MONGODB_URI not set")
	}

	if dbName == "" {
		return nil, fmt.Errorf("MONGO_DB not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri).SetTimeout(10 * time.Second))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB!")

	return &MongoConfig{
		Client: client,
		DB:     client.Database(dbName),
	}, nil
}
