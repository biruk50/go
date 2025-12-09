package Repositories

import (
	"context"
	"errors"
	"time"
	"task_manager_clean/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Create(u *Domain.User) error
	FindByUsername(username string) (*Domain.User, error)
	Count() (int64, error)
	Promote(username string) error
}

type mongoUserRepo struct {
	coll *mongo.Collection
}

func NewMongoUserRepository(db *mongo.Database) UserRepository {
	return &mongoUserRepo{coll: db.Collection("users")}
}

func (r *mongoUserRepo) Create(u *Domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.coll.InsertOne(ctx, u)
	return err
}

func (r *mongoUserRepo) FindByUsername(username string) (*Domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var u Domain.User
	if err := r.coll.FindOne(ctx, bson.M{"username": username}).Decode(&u); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &u, nil
}

func (r *mongoUserRepo) Count() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.coll.CountDocuments(ctx, bson.D{})
}

func (r *mongoUserRepo) Promote(username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := r.coll.UpdateOne(ctx, bson.M{"username": username}, bson.M{"$set": bson.M{"role": "admin"}})
	if err != nil { return err }
	if res.MatchedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}
