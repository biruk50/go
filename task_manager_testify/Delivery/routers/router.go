package routers

import (
	"task_manager_testify/Delivery/controllers"
	"task_manager_testify/Infrastructure"
	"task_manager_testify/Usecases"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userUC Usecases.UserUsecase, taskUC Usecases.TaskUsecase, jwtSvcOpt ...Infrastructure.JWTService) *gin.Engine {
	r := gin.Default()

	var jwtSvc Infrastructure.JWTService
	if len(jwtSvcOpt) > 0 {
		jwtSvc = jwtSvcOpt[0]
	} else {
		jwtSvc = Infrastructure.NewJWTService()
	}

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
