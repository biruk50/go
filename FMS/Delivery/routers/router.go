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

	ctr := controllers.NewController(userUC, budgetUC, cashRequestUC, expenseUC, reportUC, jwtSvc)

	// public
	r.POST("/register", ctr.Register)
	r.POST("/login", ctr.Login)
	r.GET("/", ctr.Home)

	user := r.Group("/users")
	user.Use(Infrastructure.AuthMiddleware(jwtSvc))
	{	
		user.GET("/", ctr.GetAllUsers)
		user.GET("/me", ctr.GetMyProfile)
	}

	user.Use(Infrastructure.AuthMiddleware(jwtSvc), Infrastructure.FinanceOnly())
	{	
		user.PUT("/:id/role", ctr.UpdateUser)
	}
	
	
	// protected
	budget := r.Group("/budgets")
	budget.Use(Infrastructure.AuthMiddleware(jwtSvc))
	{	

		budget.GET("/", ctr.GetAllBudgets)
		budget.GET("/:id", ctr.GetBudgetByID)
		budget.GET("/:id/summary", ctr.GetBudgetSummary)

		budget.PUT("/:id", ctr.UpdateBudget)

		budget.POST("/", ctr.CreateBudget)	
	}

	budget.Use(Infrastructure.AuthMiddleware(jwtSvc), Infrastructure.FinanceOnly())
	{	

		budget.PUT("/:id", ctr.UpdateBudget)

		budget.POST("/:id/approve", ctr.ApproveBudget)
		budget.POST("/:id/reject", ctr.RejectBudget)
		
	}



	cashRequest := r.Group("/cash-requests")
	cashRequest.Use(Infrastructure.AuthMiddleware(jwtSvc))
	{	

		cashRequest.GET("/", ctr.GetAllCashRequests)
		cashRequest.GET("/:id", ctr.GetCashRequest)

		cashRequest.POST("/", ctr.CreateCashRequest)
	}

	cashRequest.Use(Infrastructure.AuthMiddleware(jwtSvc), Infrastructure.FinanceOnly())
	{	
		cashRequest.POST("/:id/approve", ctr.ApproveCashRequest)
		cashRequest.POST("/:id/reject", ctr.RejectCashRequest)
		cashRequest.POST("/:id/disburse", ctr.DisburseCashRequest)
	}

	expense := r.Group("/expenses")
	expense.Use(Infrastructure.AuthMiddleware(jwtSvc))
	{	

		expense.GET("/", ctr.GetAllExpenses)
		expense.GET("/:id", ctr.GetExpense)
		expense.GET("/:id/summary", ctr.GetExpenseSummary)

		expense.POST("/", ctr.CreateExpense)
		expense.POST("/:id/receipts", ctr.CreateExpenseReceipt)

		expense.PUT("/:id/verify", ctr.VerifyExpense)
	}

	report := r.Group("/reports")
	report.Use(Infrastructure.AuthMiddleware(jwtSvc))
	{	
		report.GET("/overview", ctr.GetOverviewReport)
		report.GET("/cash-requests", ctr.GetCashRequestReport)
		report.GET("/budgets", ctr.GetBudgetReport)
		report.GET("/expenses", ctr.GetExpenseReport)
	}

	return r
}
