package main

import (
	"FMS/Delivery/routers"
	"FMS/Infrastructure"
	"FMS/Repositories"
	"FMS/Usecases"
	"log"
)

func main() {
	// init DB
	if err := Infrastructure.InitMongo(); err != nil {
		log.Fatalf("mongo init: %v", err)
	}
	defer Infrastructure.CloseMongo()

	// create repository implementations
	userRepo := Repositories.NewMongoUserRepository(Infrastructure.GetDB())
	budgetRepo := Repositories.NewMongoBudgetRepository(Infrastructure.GetDB())
	cashRepo := Repositories.NewMongoCashRequestRepository(Infrastructure.GetDB())
	expenseRepo := Repositories.NewMongoExpenseRepository(Infrastructure.GetDB())

	userUC := Usecases.NewUserUsecase(userRepo, Infrastructure.NewPasswordService(), Infrastructure.NewJWTService())
	budgetUC := Usecases.NewBudgetUsecase(budgetRepo)
	cashUC := Usecases.NewCashRequestUsecase(cashRepo)
	expenseUC := Usecases.NewExpenseUsecase(expenseRepo)
	reportUC := Usecases.NewReportUsecase(budgetRepo, cashRepo, expenseRepo)

	// create router with controllers wired to usecases
	r := routers.SetupRouter(userUC, budgetUC, cashUC, expenseUC, reportUC)

	port := Infrastructure.GetEnv("PORT", "8080")

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server run: %v", err)
	}
}
