package main

import (
	"log"
	"task_manager_testify/Delivery/routers"
	"task_manager_testify/Infrastructure"
	"task_manager_testify/Repositories"
	"task_manager_testify/Usecases"
)

func main() {
	// init DB
	if err := Infrastructure.InitMongo(); err != nil {
		log.Fatalf("mongo init: %v", err)
	}
	defer Infrastructure.CloseMongo()

	// create repository implementations
	userRepo := Repositories.NewMongoUserRepository(Infrastructure.GetDB())
	taskRepo := Repositories.NewMongoTaskRepository(Infrastructure.GetDB())

	userUC := Usecases.NewUserUsecase(userRepo, Infrastructure.NewPasswordService(), Infrastructure.NewJWTService())
	taskUC := Usecases.NewTaskUsecase(taskRepo)

	// create router with controllers wired to usecases
	r := routers.SetupRouter(userUC, taskUC)

	port := Infrastructure.GetEnv("PORT", "8080")

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server run: %v", err)
	}
}
