package repositories_integration

import (
	"context"
	"testing"
	"time"

	"task_manager_testify/Domain"
	"task_manager_testify/Infrastructure"
	"task_manager_testify/Repositories"

	"github.com/stretchr/testify/assert"
)

func TestTaskRepositoryIntegration(t *testing.T) {
	Infrastructure.LoadEnv()

	uri := Infrastructure.GetEnv("MONGODB_URL", "")
	if uri == "" {
		t.Skip("MONGODB_URL not set – skipping integration test")
	}

	// init DB (this uses your Infrastructure Connect/init)
	err := Infrastructure.InitMongo()
	if err != nil {
		t.Fatalf("init mongo: %v", err)
	}
	defer Infrastructure.CloseMongo()
	db := Infrastructure.GetDB()
	repo := Repositories.NewMongoTaskRepository(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Cleanup collection for test isolation
	_ = db.Collection("tasks").Drop(ctx)

	t1 := &Domain.Task{Title: "integration task"}
	err = repo.Create(t1)
	assert.NoError(t, err)
	// find by object id
	found, err := repo.GetByID(t1.ID.Hex())
	assert.NoError(t, err)
	assert.Equal(t, t1.Title, found.Title)

	// cleanup
	_ = db.Collection("tasks").Drop(ctx)
}
