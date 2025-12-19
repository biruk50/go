package routers

import (
	"FMS/Delivery/controllers"
	"FMS/Infrastructure"
	"FMS/Usecases"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userUC Usecases.UserUsecase, budgetUC Usecases.BudgetUsecase, cashRequestUC Usecases.CashRequestUsecase, expenseUC Usecases.ExpenseUsecase, reportUC Usecases.ReportUsecase) *gin.Engine {
	r := gin.Default()

	jwtSvc := Infrastructure.NewJWTService()

	userCtr := controllers.NewUserController(userUC, jwtSvc)
	budgetCtr := controllers.NewBudgetController(budgetUC)
	cashRequestCtr := controllers.NewCashRequestController(cashRequestUC)
	expenseCtr := controllers.NewExpenseController(expenseUC)
	reportCtr := controllers.NewReportController(reportUC)


	// public
	r.POST("/register", userCtr.Register)
	r.POST("/login", userCtr.Login)
	r.GET("/", userCtr.Home)

	user := r.Group("/users")
	user.Use(Infrastructure.AuthMiddleware(jwtSvc))
	{
		user.GET("/", userCtr.GetAllUsers)
		user.GET("/me", userCtr.GetMyProfile)
	}

	user.Use(Infrastructure.AuthMiddleware(jwtSvc), Infrastructure.FinanceOnly())
	{
		user.PUT("/:id/role", userCtr.UpdateUser)
	}

	// protected
	budget := r.Group("/budgets")
	budget.Use(Infrastructure.AuthMiddleware(jwtSvc))
	{

		budget.GET("/", budgetCtr.GetAllBudgets)
		budget.GET("/:id", budgetCtr.GetBudgetByID)
		budget.GET("/:id/summary", budgetCtr.GetBudgetSummary)

		budget.PUT("/:id", budgetCtr.UpdateBudget)
		budget.POST("/", budgetCtr.CreateBudget)
	}

	budget.Use(Infrastructure.AuthMiddleware(jwtSvc), Infrastructure.FinanceOnly())
	{
		budget.POST("/:id/approve", budgetCtr.ApproveBudget)
		budget.POST("/:id/reject", budgetCtr.RejectBudget)
	}

	cashRequest := r.Group("/cash-requests")
	cashRequest.Use(Infrastructure.AuthMiddleware(jwtSvc))
	{

		cashRequest.GET("/", cashRequestCtr.GetAllCashRequests)
		cashRequest.GET("/:id", cashRequestCtr.GetCashRequest)

		cashRequest.POST("/", cashRequestCtr.CreateCashRequest)
	}

	cashRequest.Use(Infrastructure.AuthMiddleware(jwtSvc), Infrastructure.FinanceOnly())
	{
		cashRequest.POST("/:id/approve", cashRequestCtr.ApproveCashRequest)
		cashRequest.POST("/:id/reject", cashRequestCtr.RejectCashRequest)
		cashRequest.POST("/:id/disburse", cashRequestCtr.DisburseCashRequest)
	}

	expense := r.Group("/expenses")
	expense.Use(Infrastructure.AuthMiddleware(jwtSvc))
	{

		expense.GET("/", expenseCtr.GetAllExpenses)
		expense.GET("/:id", expenseCtr.GetExpense)
		expense.GET("/:id/summary", expenseCtr.GetExpenseSummary)

		expense.POST("/", expenseCtr.CreateExpense)
		expense.POST("/:id/receipts", expenseCtr.CreateExpenseReceipt)
	}

	expense.Use(Infrastructure.AuthMiddleware(jwtSvc), Infrastructure.FinanceOnly())
	{
		expense.PUT("/:id/verify", expenseCtr.VerifyExpense)
	}

	report := r.Group("/reports")
	report.Use(Infrastructure.AuthMiddleware(jwtSvc), Infrastructure.FinanceOnly())
	{
		report.GET("/overview", reportCtr.GetOverviewReport)
		report.GET("/cash-requests", reportCtr.GetCashRequestReport)
		report.GET("/budgets", reportCtr.GetBudgetReport)
		report.GET("/expenses", reportCtr.GetExpenseReport)
	}

	return r
}
