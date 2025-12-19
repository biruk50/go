package Infrastructure

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client
var db *mongo.Database

func InitMongo() error {
	_=LoadEnv()
	uri := GetEnv("MONGODB_URL", "")
	if uri == "" {
		return errors.New("MONGODB_URL not set")
	}
	dbName := GetEnv("MONGO_DB", "taskdb")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	db = client.Database(dbName)
	return nil
}

func GetDB() *mongo.Database {
	return db
}

func CloseMongo() {
	if client == nil { return }
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = client.Disconnect(ctx)
}
