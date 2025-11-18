package routers

import (
	"task_manager_clean/Delivery/controllers"
	"task_manager_clean/Infrastructure"
	"task_manager_clean/Usecases"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userUC Usecases.UserUsecase, taskUC Usecases.TaskUsecase) *gin.Engine {
	r := gin.Default()

	jwtSvc := Infrastructure.NewJWTService()

	ctr := controllers.NewController(userUC, taskUC, jwtSvc)

	// public
	r.POST("/register", ctr.Register)
	r.POST("/login", ctr.Login)
	r.GET("/", ctr.Home)

	// protected
	auth := r.Group("/")
	auth.Use(Infrastructure.AuthMiddleware(jwtSvc))
	{
		auth.GET("/tasks", ctr.GetAllTasks)
		auth.GET("/tasks/:id", ctr.GetTask)
	}

	// admin-only
	admin := r.Group("/")
	admin.Use(Infrastructure.AuthMiddleware(jwtSvc), Infrastructure.AdminOnly())
	{
		admin.POST("/tasks", ctr.CreateTask)
		admin.PUT("/tasks/:id", ctr.UpdateTask)
		admin.DELETE("/tasks/:id", ctr.DeleteTask)
		admin.POST("/promote", ctr.Promote)
	}

	return r
}
