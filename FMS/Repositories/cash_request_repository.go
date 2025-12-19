package Repositories

import (
	"FMS/Domain"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CashRequestRepository interface {
	Create(t *Domain.CashRequest) error
	GetAll() ([]Domain.CashRequest, error)
	GetByID(id string) (*Domain.CashRequest, error)
	Update(id string, t *Domain.CashRequest) error
	Delete(id string) error
}

type mongoCashRequestRepo struct {
	coll *mongo.Collection
}

func NewMongoCashRequestRepository(db *mongo.Database) CashRequestRepository {
	return &mongoCashRequestRepo{coll: db.Collection("cash_requests")}
}

func (r *mongoCashRequestRepo) Create(t *Domain.CashRequest) error {
	t.ID = primitive.NewObjectID()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.coll.InsertOne(ctx, t)
	return err
}

func (r *mongoCashRequestRepo) GetAll() ([]Domain.CashRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []Domain.CashRequest
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *mongoCashRequestRepo) GetByID(id string) (*Domain.CashRequest, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cr Domain.CashRequest
	if err := r.coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&cr); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("cash request not found")
		}
		return nil, err
	}

	return &cr, nil
}

func (r *mongoCashRequestRepo) Update(id string, t *Domain.CashRequest) error {
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
			"amount":      t.Amount,
			"budget_id":   t.BudgetID,
			"requester":   t.Requester,
			"created_at":  t.CreatedAt,
			"status":      t.Status,
		},
	}

	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("cash request not found")
	}

	return nil
}

func (r *mongoCashRequestRepo) Delete(id string) error {
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
		return errors.New("cash request not found")
	}

	return nil
}
