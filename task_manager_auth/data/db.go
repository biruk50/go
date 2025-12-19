package data

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Global DB variables accessible to other data files
var (
	Client       *mongo.Client
	DB           *mongo.Database
	UsersColl    *mongo.Collection
	TasksColl    *mongo.Collection
	JWTSecretKey string
)

// InitMongo initializes MongoDB client and sets up collections
func InitMongo() error {
	_ = godotenv.Load()

	mongoURI := os.Getenv("MONGODB_URL")
	if mongoURI == "" {
		return errors.New("MONGODB_URL not set")
	}

	JWTSecretKey = os.Getenv("JWT_SECRET")

	dbName := getEnv("MONGO_DB", "taskdb")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	Client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
	}

	DB = Client.Database(dbName)
	UsersColl = DB.Collection(getEnv("MONGO_USERS_COLLECTION", "users"))
	TasksColl = DB.Collection(getEnv("MONGO_TASKS_COLLECTION", "tasks"))

	return nil
}

// CloseMongo closes the MongoDB client connection
func CloseMongo() {
	if Client == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = Client.Disconnect(ctx)
}

// Helper to get env var with fallback
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// GetJWTSecret returns the secret key for JWT signing
func GetJWTSecret() string {
	return JWTSecretKey
}
