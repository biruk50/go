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

type BudgetRepository interface {
	Create(t *Domain.Budget) error
	GetAll() ([]Domain.Budget, error)
	GetByID(id string) (*Domain.Budget, error)
	Update(id string, t *Domain.Budget) error
	Delete(id string) error
}

type mongoBudgetRepo struct {
	coll *mongo.Collection
}

func NewMongoBudgetRepository(db *mongo.Database) BudgetRepository {
	return &mongoBudgetRepo{coll: db.Collection("budgets")}
}

func (r *mongoBudgetRepo) Create(t *Domain.Budget) error {
	t.ID = primitive.NewObjectID()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.coll.InsertOne(ctx, t)
	return err
}

func (r *mongoBudgetRepo) GetAll() ([]Domain.Budget, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var budgets []Domain.Budget
	if err := cursor.All(ctx, &budgets); err != nil {
		return nil, err
	}

	return budgets, nil
}

func (r *mongoBudgetRepo) GetByID(id string) (*Domain.Budget, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var b Domain.Budget
	if err := r.coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&b); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("budget not found")
		}
		return nil, err
	}

	return &b, nil
}

func (r *mongoBudgetRepo) Update(id string, t *Domain.Budget) error {
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
			"amount":      t.Amount,
			"remaining":   t.Remaining,
		},
	}

	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("budget not found")
	}

	return nil
}

func (r *mongoBudgetRepo) Delete(id string) error {
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
		return errors.New("budget not found")
	}

	return nil
}
