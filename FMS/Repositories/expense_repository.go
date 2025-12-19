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

type ExpenseRepository interface {
	Create(t *Domain.Expense) error
	GetAll() ([]Domain.Expense, error)
	GetByID(id string) (*Domain.Expense, error)
	Update(id string, t *Domain.Expense) error
	Delete(id string) error
}

type mongoExpenseRepo struct {
	coll *mongo.Collection
}

func NewMongoExpenseRepository(db *mongo.Database) ExpenseRepository {
	return &mongoExpenseRepo{coll: db.Collection("expenses")}
}

func (r *mongoExpenseRepo) Create(t *Domain.Expense) error {
	t.ID = primitive.NewObjectID()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.coll.InsertOne(ctx, t)
	return err
}

func (r *mongoExpenseRepo) GetAll() ([]Domain.Expense, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []Domain.Expense
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *mongoExpenseRepo) GetByID(id string) (*Domain.Expense, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var e Domain.Expense
	if err := r.coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&e); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("expense not found")
		}
		return nil, err
	}

	return &e, nil
}

func (r *mongoExpenseRepo) Update(id string, t *Domain.Expense) error {
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
			"receipt_url": t.ReceiptURL,
			"budget_id":   t.BudgetID,
			"created_at":  t.CreatedAt,
			"status":      t.Status,
			"due_date":    t.DueDate,
		},
	}

	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("expense not found")
	}

	return nil
}

func (r *mongoExpenseRepo) Delete(id string) error {
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
