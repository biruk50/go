package Repositories

import (
	"context"
	"errors"
	"time"
	"FMS/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReportRepository interface {
	Create(t *Domain.Report) error
	GetAll() ([]Domain.Report, error)
	GetByID(id string) (*Domain.Report, error)
	Update(id string, t *Domain.Report) error
	Delete(id string) error
}

type mongoReportRepo struct {
	coll *mongo.Collection
}

func NewMongoReportRepository(db *mongo.Database) ReportRepository {
	return &mongoReportRepo{coll: db.Collection("reports")}
}

func (r *mongoReportRepo) Create(t *Domain.Report) error {
	t.ID = primitive.NewObjectID()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.coll.InsertOne(ctx, t)
	return err
}

func (r *mongoReportRepo) GetAll() ([]Domain.Report, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []Domain.Report
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *mongoReportRepo) GetByID(id string) (*Domain.Report, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var task Domain.Report
	if err := r.coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&task); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("task not found")
		}
		return nil, err
	}

	return &task, nil
}

func (r *mongoReportRepo) Update(id string, t *Domain.Report) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"title":       t.Title,
			"description": t.Description,
			"status":      t.Status,
			"due_date":    t.DueDate,
		},
	}

	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("task not found")
	}

	return nil
}

func (r *mongoReportRepo) Delete(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := r.coll.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("task not found")
	}

	return nil
}
