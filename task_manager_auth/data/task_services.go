package data

import (
	"context"
	"errors"
	"log"
	"os"
	"task_manager/models"
	"time"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
  collection *mongo.Collection
  client *mongo.Client
  ctx=context.Background()
)



func InitMongo() error {
	if err := godotenv.Load(); err != nil {
		return errors.New("no .env file found")
	}

	uri := os.Getenv("MONGODB_URL")
	
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil { return err }

	collection = client.Database("task_manager").Collection("tasks")
	if collection == nil { return errors.New("failed to connect to collection") }
	return nil

}

func CloseMongo(){
	if client == nil { return  }
	client.Disconnect(ctx) 
}

// GetAllTasks returns all tasks.
func GetAllTasks() ([]models.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	cursor,err := collection.Find(ctx, bson.D{})
	if err != nil { return nil, err }

	var tasks []models.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil

}

// GetTaskByID returns a single task by ID.
func GetTaskByID(id string) (*models.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}


	var task models.Task

	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&task)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("task not found")
	}
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func AddTask(task models.Task) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, task)
	return err
}

func UpdateTask(id string, updated models.Task) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"title":       updated.Title,
			"description": updated.Description,
			"due_date":    updated.DueDate,
			"status":      updated.Status,
		},
	}

	res, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("task not found")
	}
	return nil
}

// DeleteTask removes a task by its ID.
func DeleteTask(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("format")
		return errors.New("invalid ID format")
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	res, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	
	if err != nil { 
		log.Printf("Error deleting task: %v", err)
		return err}
	log.Printf("Delete result: %+v", res)
	

	return nil
}